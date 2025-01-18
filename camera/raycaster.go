// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package camera

import (
	"github.com/derekmu/g3n/core"
	"github.com/derekmu/g3n/gls"
	"github.com/derekmu/g3n/graphic"
	"github.com/derekmu/g3n/material"
	"github.com/derekmu/g3n/math32"
	"sort"
)

// Raycaster represents an empty object that can cast rays and check for ray intersections.
type Raycaster struct {
	math32.Ray
	// The distance from the ray origin to the intersected points
	// must be greater than the value of this field to be considered.
	// The default value is 0.0
	Near float32
	// The distance from the ray origin to the intersected points
	// must be less than the value of this field to be considered.
	// The default value is +Infinity.
	Far float32
	// Minimum distance in world coordinates between the ray and
	// a line segment when checking intersects with lines.
	// The default value is 0.1
	LinePrecision float32
	// Minimum distance in world coordinates between the ray and
	// a point when checking intersects with points.
	// The default value is 0.1
	PointPrecision float32
	// This field must be set with the camera view matrix used
	// when checking for sprite intersections.
	// It is set automatically when using camera.SetRaycaster
	ViewMatrix math32.Matrix4
}

// Intersection describes the intersection between a ray and an object
type Intersection struct {
	// Distance between the origin of the ray and the intersection
	Distance float32
	// Point of intersection in world coordinates
	Point math32.Vector3
	// Intersected node
	Object core.INode
	// If the geometry has indices, this field is the
	// index in the Indices buffer of the vertex intersected
	// or the first vertex of the intersected face.
	// If the geometry doesn't have indices, this field is the
	// index in the positions buffer of the vertex intersected
	// or the first vertex of the intersected face.
	Index uint32
}

// NewRaycaster creates and returns a pointer to a new raycaster object
// with the specified origin and direction.
func NewRaycaster(origin, direction math32.Vector3) *Raycaster {
	rc := new(Raycaster)
	rc.Ray = math32.Ray{
		Origin:    origin,
		Direction: direction,
	}
	rc.Near = 0
	rc.Far = math32.Inf(1)
	rc.LinePrecision = 0.1
	rc.PointPrecision = 0.1
	return rc
}

// IntersectObject checks intersections between this raycaster
// and the specified node. If recursive is true, it also checks
// the intersection with the node's children.
// Intersections are returned sorted by distance, closest first.
func (rc *Raycaster) IntersectObject(inode core.INode, recursive bool, intersections []Intersection) []Intersection {
	intersections = rc.intersectObject(inode, recursive, intersections)
	sort.Slice(intersections, func(i, j int) bool {
		return intersections[i].Distance < intersections[j].Distance
	})
	return intersections
}

// IntersectObjects checks intersections between this raycaster and
// the specified array of scene nodes. If recursive is true, it also checks
// the intersection with each nodes' children.
// Intersections are returned sorted by distance, closest first.
func (rc *Raycaster) IntersectObjects(inodes []core.INode, recursive bool, intersections []Intersection) []Intersection {
	for _, inode := range inodes {
		intersections = rc.intersectObject(inode, recursive, intersections)
	}
	sort.Slice(intersections, func(i, j int) bool {
		return intersections[i].Distance < intersections[j].Distance
	})
	return intersections
}

func (rc *Raycaster) intersectObject(inode core.INode, recursive bool, intersections []Intersection) []Intersection {
	node := inode.GetNode()
	if !node.Visible() {
		return intersections
	}

	switch in := inode.(type) {
	case *graphic.Sprite:
		if intersect, ok := rc.RaycastSprite(in); ok {
			intersections = append(intersections, intersect)
		}
	case *graphic.Points:
		intersections = rc.RaycastPoints(in, intersections)
	case *graphic.Mesh:
		intersections = rc.RaycastMesh(in, intersections)
	case *graphic.Lines:
		intersections = rc.RaycastLines(in, intersections)
	case *graphic.LineStrip:
		intersections = rc.RaycastLineStrip(in, intersections)
	}

	if recursive {
		for _, child := range node.Children() {
			intersections = rc.intersectObject(child, true, intersections)
		}
	}
	return intersections
}

// SetFromCamera sets the specified raycaster with this camera position in world coordinates
// pointing to the direction defined by the specified coordinates unprojected using this camera.
func (rc *Raycaster) SetFromCamera(cam *Camera, sx, sy float32) {
	matrixWorld := cam.MatrixWorld()
	rc.Origin.SetFromMatrixPosition(&matrixWorld)
	rc.Direction.Set(sx, sy, 0.5)
	cam.Unproject(&rc.Direction)
	rc.Direction.Sub(&rc.Origin).Normalize()
	cam.ViewMatrix(&rc.ViewMatrix)
}

// RaycastSprite checks intersections between the raycaster and the specified sprite
// and if any found appends it to the specified intersects array.
func (rc *Raycaster) RaycastSprite(s *graphic.Sprite) (Intersection, bool) {
	// Copy and convert ray to camera coordinates
	ray := rc.Ray
	ray.ApplyMatrix4(&rc.ViewMatrix)

	// Calculates ViewMatrix * MatrixWorld
	var mv math32.Matrix4
	matrixWorld := s.MatrixWorld()
	mv.MultiplyMatrices(&rc.ViewMatrix, &matrixWorld)

	// Decompose transformation matrix in its components
	var position math32.Vector3
	var quaternion math32.Quaternion
	var scale math32.Vector3
	mv.Decompose(&position, &quaternion, &scale)

	// Remove any rotation in X and Y axis and
	// compose new transformation matrix
	rotation := s.Rotation()
	rotation.X = 0
	rotation.Y = 0
	quaternion.SetFromEuler(&rotation)
	mv.Compose(&position, &quaternion, &scale)

	// Get buffer with vertices and uvs
	geom := s.GetGeometry()
	vboPos := geom.VBO(gls.VertexPosition)
	if vboPos == nil {
		panic("sprite.Raycast(): VertexPosition VBO not found")
	}
	// Get vertex positions, transform to camera coordinates and
	// checks intersection with ray
	buffer := vboPos.Buffer()
	indices := geom.Indices()
	var v1 math32.Vector3
	var v2 math32.Vector3
	var v3 math32.Vector3
	var point math32.Vector3
	intersect := false
	for i := 0; i < indices.Len(); i += 3 {
		pos := indices[i]
		buffer.GetVector3(int(pos*5), &v1)
		v1.ApplyMatrix4(&mv)
		pos = indices[i+1]
		buffer.GetVector3(int(pos*5), &v2)
		v2.ApplyMatrix4(&mv)
		pos = indices[i+2]
		buffer.GetVector3(int(pos*5), &v3)
		v3.ApplyMatrix4(&mv)
		if point, intersect = ray.IntersectTriangle(&v1, &v2, &v3, false); intersect {
			break
		}
	}
	if !intersect {
		return Intersection{}, false
	}
	// Get distance from intersection point
	distance := ray.Origin.DistanceTo(&point)

	// Checks if distance is between the bounds of the raycaster
	if distance < rc.Near || distance > rc.Far {
		return Intersection{}, false
	}

	// Appends intersection to received parameter.
	return Intersection{
		Distance: distance,
		Point:    point,
		Object:   s,
	}, true
}

// RaycastPoints checks for intersections with a graphic.Points.
func (rc *Raycaster) RaycastPoints(p *graphic.Points, intersections []Intersection) []Intersection {
	// Checks intersection with the bounding sphere transformed to world coordinates
	geom := p.GetGeometry()
	sphere := geom.BoundingSphere()
	matrixWorld := p.MatrixWorld()
	sphere.ApplyMatrix4(&matrixWorld)
	if !rc.IsIntersectionSphere(&sphere) {
		return intersections
	}

	// Copy ray and transforms to model coordinates
	var inverseMatrix math32.Matrix4
	_ = inverseMatrix.GetInverse(&matrixWorld)
	ray := rc.Ray
	ray.ApplyMatrix4(&inverseMatrix)

	// Checks intersection with all points
	scale := p.Scale()
	localThreshold := rc.PointPrecision / ((scale.X + scale.Y + scale.Z) / 3)
	localThresholdSq := localThreshold * localThreshold

	i := 0
	geom.ReadVertices(func(point math32.Vector3) bool {
		rayPointDistanceSq := ray.DistanceSqToPoint(&point)
		if rayPointDistanceSq >= localThresholdSq {
			return false
		}
		intersectPoint := ray.ClosestPointToPoint(&point)
		intersectPoint.ApplyMatrix4(&matrixWorld)
		distance := rc.Origin.DistanceTo(&intersectPoint)
		if distance < rc.Near || distance > rc.Far {
			return false
		}
		intersections = append(intersections, Intersection{
			Distance: distance,
			Point:    intersectPoint,
			Index:    uint32(i),
			Object:   p,
		})
		i++
		return false
	})
	return intersections
}

// RaycastMesh checks for intersections with a graphic.Mesh.
func (rc *Raycaster) RaycastMesh(m *graphic.Mesh, intersections []Intersection) []Intersection {
	// Transform this mesh geometry bounding sphere from model to world coordinates and checks intersection with raycaster
	geom := m.GetGeometry()
	sphere := geom.BoundingSphere()
	matrixWorld := m.MatrixWorld()
	sphere.ApplyMatrix4(&matrixWorld)
	if !rc.IsIntersectionSphere(&sphere) {
		return intersections
	}
	// Copy ray and transform to model coordinates
	// It is less expensive to transform the ray to model coordinates than the geometry to world coordinates
	// This ray will also be used to check intersects with the geometry
	var inverseMatrix math32.Matrix4
	_ = inverseMatrix.GetInverse(&matrixWorld)
	ray := rc.Ray
	ray.ApplyMatrix4(&inverseMatrix)
	bbox := geom.BoundingBox()
	if !ray.IsIntersectionBox(&bbox) {
		return intersections
	}
	i := 0
	geom.ReadFaces(func(pA, pB, pC math32.Vector3) bool {
		// Checks intersection of the ray with this face
		side := m.GetMaterial(i).GetMaterial().Side()
		var point math32.Vector3
		var intersect bool
		switch side {
		case material.SideBack:
			point, intersect = ray.IntersectTriangle(&pC, &pB, &pA, true)
		case material.SideFront:
			point, intersect = ray.IntersectTriangle(&pA, &pB, &pC, true)
		case material.SideDouble:
			point, intersect = ray.IntersectTriangle(&pA, &pB, &pC, false)
		}
		if intersect {
			// Transform intersection point from model to world coordinates
			point.ApplyMatrix4(&matrixWorld)
			// Calculates the distance from the ray origin to intersection point
			distance := rc.Ray.Origin.DistanceTo(&point)
			// Checks if distance is between the bounds of the raycaster
			if distance >= rc.Near && distance <= rc.Far {
				intersections = append(intersections,
					Intersection{
						Distance: distance,
						Point:    point,
						Object:   m,
						Index:    uint32(i),
					})
			}
		}
		i += 3
		return false
	})
	return intersections
}

// RaycastMeshFirst checks for intersections with a graphic.Mesh.
// This will return the first intersection found.
func (rc *Raycaster) RaycastMeshFirst(m *graphic.Mesh) (result Intersection, intersects bool) {
	// Transform this mesh geometry bounding sphere from model to world coordinates and checks intersection with raycaster
	geom := m.GetGeometry()
	sphere := geom.BoundingSphere()
	matrixWorld := m.MatrixWorld()
	sphere.ApplyMatrix4(&matrixWorld)
	if !rc.IsIntersectionSphere(&sphere) {
		return Intersection{}, false
	}
	// Copy ray and transform to model coordinates
	// It is less expensive to transform the ray to model coordinates than the geometry to world coordinates
	// This ray will also be used to check intersects with the geometry
	var inverseMatrix math32.Matrix4
	_ = inverseMatrix.GetInverse(&matrixWorld)
	ray := rc.Ray
	ray.ApplyMatrix4(&inverseMatrix)
	bbox := geom.BoundingBox()
	if !ray.IsIntersectionBox(&bbox) {
		return Intersection{}, false
	}
	i := 0
	geom.ReadFaces(func(pA, pB, pC math32.Vector3) bool {
		// Checks intersection of the ray with this face
		side := m.GetMaterial(i).GetMaterial().Side()
		var point math32.Vector3
		var intersect bool
		switch side {
		case material.SideBack:
			point, intersect = ray.IntersectTriangle(&pC, &pB, &pA, true)
		case material.SideFront:
			point, intersect = ray.IntersectTriangle(&pA, &pB, &pC, true)
		case material.SideDouble:
			point, intersect = ray.IntersectTriangle(&pA, &pB, &pC, false)
		}
		if intersect {
			// Transform intersection point from model to world coordinates
			point.ApplyMatrix4(&matrixWorld)
			// Calculates the distance from the ray origin to intersection point
			distance := rc.Ray.Origin.DistanceTo(&point)
			// Checks if distance is between the bounds of the raycaster
			if distance >= rc.Near && distance <= rc.Far {
				result = Intersection{
					Distance: distance,
					Point:    point,
					Object:   m,
					Index:    uint32(i),
				}
				intersects = true
			}
		}
		i += 3
		return intersects
	})
	return result, intersects
}

// RaycastLines checks for intersections with a graphic.Lines.
func (rc *Raycaster) RaycastLines(l *graphic.Lines, intersections []Intersection) []Intersection {
	return rc.lineRaycast(l, intersections, 2)
}

// RaycastLineStrip checks for intersections with a graphic.LineStrip.
func (rc *Raycaster) RaycastLineStrip(l *graphic.LineStrip, intersections []Intersection) []Intersection {
	return rc.lineRaycast(l, intersections, 1)
}

// Internal function used by raycasting for Lines and LineStrip.
func (rc *Raycaster) lineRaycast(igr graphic.IGraphic, intersections []Intersection, step int) []Intersection {
	// Get the bounding sphere
	gr := igr.GetGraphic()
	geom := igr.GetGeometry()
	sphere := geom.BoundingSphere()

	// Transform bounding sphere from model to world coordinates and
	// checks intersection with raycaster
	matrixWorld := gr.MatrixWorld()
	sphere.ApplyMatrix4(&matrixWorld)
	if !rc.IsIntersectionSphere(&sphere) {
		return intersections
	}

	// Copy ray and transform to model coordinates
	// This ray will also be used to check intersects with
	// the geometry, as is much less expensive to transform the
	// ray to model coordinates than the geometry to world coordinates.
	var inverseMatrix math32.Matrix4
	_ = inverseMatrix.GetInverse(&matrixWorld)
	ray := rc.Ray
	ray.ApplyMatrix4(&inverseMatrix)

	var vstart math32.Vector3
	var vend math32.Vector3
	var interSegment math32.Vector3
	var interRay math32.Vector3

	// Get geometry positions and indices buffers
	vboPos := geom.VBO(gls.VertexPosition)
	if vboPos == nil {
		return intersections
	}
	positions := vboPos.Buffer()
	indices := geom.Indices()
	precisionSq := rc.LinePrecision * rc.LinePrecision

	// Checks intersection with individual lines for indexed geometry
	if indices.Len() > 0 {
		for i := 0; i < indices.Len()-1; i += step {
			// Calculates distance from ray to this line segment
			a := indices[i]
			b := indices[i+1]
			positions.GetVector3(int(3*a), &vstart)
			positions.GetVector3(int(3*b), &vend)
			distSq := ray.DistanceSqToSegment(&vstart, &vend, &interRay, &interSegment)
			if distSq > precisionSq {
				continue
			}
			// Move back to world coordinates for distance calculation
			interRay.ApplyMatrix4(&matrixWorld)
			distance := rc.Origin.DistanceTo(&interRay)
			if distance < rc.Near || distance > rc.Far {
				continue
			}

			interSegment.ApplyMatrix4(&matrixWorld)
			intersections = append(intersections, Intersection{
				Distance: distance,
				Point:    interSegment,
				Index:    uint32(i),
				Object:   igr,
			})
		}
		// Checks intersection with individual lines for NON indexed geometry
	} else {
		for i := 0; i < positions.Len()/3-1; i += step {
			positions.GetVector3(3*i, &vstart)
			positions.GetVector3(3*i+3, &vend)
			distSq := ray.DistanceSqToSegment(&vstart, &vend, &interRay, &interSegment)
			if distSq > precisionSq {
				continue
			}

			// Move back to world coordinates for distance calculation
			interRay.ApplyMatrix4(&matrixWorld)
			distance := rc.Origin.DistanceTo(&interRay)
			if distance < rc.Near || distance > rc.Far {
				continue
			}

			interSegment.ApplyMatrix4(&matrixWorld)
			intersections = append(intersections, Intersection{
				Distance: distance,
				Point:    interSegment,
				Index:    uint32(i),
				Object:   igr,
			})
		}
	}
	return intersections
}
