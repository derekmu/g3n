// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math32

// Ray represents an oriented 3D line segment defined by an origin point and a direction vector.
type Ray struct {
	Origin    Vector3
	Direction Vector3
}

// At calculates the point in the ray which is at the specified distance from the origin along its direction.
func (ray *Ray) At(distance float32) (result Vector3) {
	result = ray.Direction
	result.MultiplyScalar(distance)
	result.Add(&ray.Origin)
	return result
}

// Recast sets the new origin of the ray at the specified distance from its origin along its direction.
func (ray *Ray) Recast(distance float32) *Ray {
	ray.Origin = ray.At(distance)
	return ray
}

// ClosestPointToPoint calculates the point in the ray which is closest to the specified point.
func (ray *Ray) ClosestPointToPoint(point *Vector3) (result Vector3) {
	result.SubVectors(point, &ray.Origin)
	directionDistance := result.Dot(&ray.Direction)
	if directionDistance < 0 {
		return ray.Origin
	}
	result = ray.Direction
	result.MultiplyScalar(directionDistance)
	result.Add(&ray.Origin)
	return result
}

// DistanceToPoint returns the smallest distance
// from the ray direction vector to the specified point.
func (ray *Ray) DistanceToPoint(point *Vector3) float32 {
	return Sqrt(ray.DistanceSqToPoint(point))
}

// DistanceSqToPoint returns the smallest squared distance from the ray direction vector to the specified point.
// If the ray was pointed directly at the point this distance would be 0.
func (ray *Ray) DistanceSqToPoint(point *Vector3) float32 {
	var v1 Vector3
	directionDistance := v1.SubVectors(point, &ray.Origin).Dot(&ray.Direction)
	// point behind the ray
	if directionDistance < 0 {
		return ray.Origin.DistanceTo(point)
	}
	v1.Copy(&ray.Direction).MultiplyScalar(directionDistance).Add(&ray.Origin)
	return v1.DistanceToSquared(point)
}

// DistanceSqToSegment returns the smallest squared distance
// from this ray to the line segment from v0 to v1.
// If optionalPointOnRay Vector3 is not nil,
// it is set with the coordinates of the point on the ray.
// if optionalPointOnSegment Vector3 is not nil,
// it is set with the coordinates of the point on the segment.
func (ray *Ray) DistanceSqToSegment(v0, v1, optionalPointOnRay, optionalPointOnSegment *Vector3) float32 {
	var segCenter Vector3
	var segDir Vector3
	var diff Vector3

	segCenter.Copy(v0).Add(v1).MultiplyScalar(0.5)
	segDir.Copy(v1).Sub(v0).Normalize()
	diff.Copy(&ray.Origin).Sub(&segCenter)

	segExtent := v0.DistanceTo(v1) * 0.5
	a01 := -ray.Direction.Dot(&segDir)
	b0 := diff.Dot(&ray.Direction)
	b1 := -diff.Dot(&segDir)
	c := diff.LengthSq()
	det := Abs(1 - a01*a01)

	var s0, s1, sqrDist, extDet float32

	if det > 0 {
		// The ray and segment are not parallel.
		s0 = a01*b1 - b0
		s1 = a01*b0 - b1
		extDet = segExtent * det

		if s0 >= 0 {
			if s1 >= -extDet {
				if s1 <= extDet {
					// region 0
					// Minimum at interior points of ray and segment.
					invDet := 1 / det
					s0 *= invDet
					s1 *= invDet
					sqrDist = s0*(s0+a01*s1+2*b0) + s1*(a01*s0+s1+2*b1) + c

				} else {
					// region 1
					s1 = segExtent
					s0 = Max(0, -(a01*s1 + b0))
					sqrDist = -s0*s0 + s1*(s1+2*b1) + c
				}

			} else {
				// region 5
				s1 = -segExtent
				s0 = Max(0, -(a01*s1 + b0))
				sqrDist = -s0*s0 + s1*(s1+2*b1) + c

			}

		} else {
			if s1 <= -extDet {
				// region 4
				s0 = Max(0, -(-a01*segExtent + b0))
				if s0 > 0 {
					s1 = -segExtent
				} else {
					s1 = Min(Max(-segExtent, -b1), segExtent)
				}
				sqrDist = -s0*s0 + s1*(s1+2*b1) + c

			} else if s1 <= extDet {
				// region 3
				s0 = 0
				s1 = Min(Max(-segExtent, -b1), segExtent)
				sqrDist = s1*(s1+2*b1) + c

			} else {
				// region 2
				s0 = Max(0, -(a01*segExtent + b0))
				if s0 > 0 {
					s1 = segExtent
				} else {
					s1 = Min(Max(-segExtent, -b1), segExtent)
				}
				sqrDist = -s0*s0 + s1*(s1+2*b1) + c
			}
		}
	} else {
		// Ray and segment are parallel.
		if a01 > 0 {
			s1 = -segExtent
		} else {
			s1 = segExtent
		}
		s0 = Max(0, -(a01*s1 + b0))
		sqrDist = -s0*s0 + s1*(s1+2*b1) + c

	}

	if optionalPointOnRay != nil {
		optionalPointOnRay.Copy(&ray.Direction).MultiplyScalar(s0).Add(&ray.Origin)
	}

	if optionalPointOnSegment != nil {
		optionalPointOnSegment.Copy(&segDir).MultiplyScalar(s1).Add(&segCenter)
	}
	return sqrDist
}

// IsIntersectionSphere returns if this ray intersects with the specified sphere.
func (ray *Ray) IsIntersectionSphere(sphere *Sphere) bool {
	if ray.DistanceToPoint(&sphere.Center) <= sphere.Radius {
		return true
	}
	return false
}

// IntersectSphere calculates the point which is the intersection of this ray with the specified sphere.
func (ray *Ray) IntersectSphere(sphere *Sphere) (Vector3, bool) {
	var v1 Vector3

	v1.SubVectors(&sphere.Center, &ray.Origin)
	tca := v1.Dot(&ray.Direction)
	d2 := v1.Dot(&v1) - tca*tca
	radius2 := sphere.Radius * sphere.Radius

	if d2 > radius2 {
		return Vector3{}, false
	}

	thc := Sqrt(radius2 - d2)

	// t0 = first intersect point - entrance on front of sphere
	t0 := tca - thc

	// t1 = second intersect point - exit point on back of sphere
	t1 := tca + thc

	// test to see if both t0 and t1 are behind the ray - if so, return null
	if t0 < 0 && t1 < 0 {
		return Vector3{}, false
	}

	// test to see if t0 is behind the ray:
	// if it is, the ray is inside the sphere, so return the second exit point scaled by t1,
	// in order to always return an intersect point that is in front of the ray.
	if t0 < 0 {
		return ray.At(t1), true
	}

	// else t0 is in front of the ray, so return the first collision point scaled by t0
	return ray.At(t0), true
}

// IsIntersectPlane returns if this ray intersects the specified plane.
func (ray *Ray) IsIntersectPlane(plane *Plane) bool {
	distToPoint := plane.DistanceToPoint(&ray.Origin)
	if distToPoint == 0 {
		return true
	}
	denominator := plane.normal.Dot(&ray.Direction)
	if denominator*distToPoint < 0 {
		return true
	}
	// ray origin is behind the plane (and is pointing behind it)
	return false
}

// DistanceToPlane returns the distance of this ray origin to its intersection point in the plane.
// If the ray does not intersect the plane, returns NaN.
func (ray *Ray) DistanceToPlane(plane *Plane) float32 {
	denominator := plane.normal.Dot(&ray.Direction)
	if denominator == 0 {
		// line is coplanar, return origin
		if plane.DistanceToPoint(&ray.Origin) == 0 {
			return 0
		}
		return NaN()
	}
	t := -(ray.Origin.Dot(&plane.normal) + plane.constant) / denominator
	// Return if the ray never intersects the plane
	if t >= 0 {
		return t
	}
	return NaN()
}

// IntersectPlane calculates the point which is the intersection of this ray with the specified plane.
func (ray *Ray) IntersectPlane(plane *Plane) (Vector3, bool) {
	t := ray.DistanceToPlane(plane)
	if IsNaN(t) {
		return Vector3{}, false
	}
	return ray.At(t), true
}

// IsIntersectionBox returns if this ray intersects the specified box.
func (ray *Ray) IsIntersectionBox(box *Box3) bool {
	_, ok := ray.IntersectBox(box)
	return ok
}

// IntersectBox calculates the point which is the intersection of this ray with the specified box.
func (ray *Ray) IntersectBox(b *Box3) (Vector3, bool) {
	// Calculate intersection parameters for the X-axis by finding the distance
	// (t) at which the ray intersects the minimum and maximum X planes of the box
	tMin := (b.Min.X - ray.Origin.X) / ray.Direction.X
	tMax := (b.Max.X - ray.Origin.X) / ray.Direction.X
	// Ensure tMin is the smaller of the two values
	if tMin > tMax {
		tMin, tMax = tMax, tMin
	}

	// Repeat the process for the Y-axis
	tyMin := (b.Min.Y - ray.Origin.Y) / ray.Direction.Y
	tyMax := (b.Max.Y - ray.Origin.Y) / ray.Direction.Y
	// Ensure tyMin is the smaller of the two values
	if tyMin > tyMax {
		tyMin, tyMax = tyMax, tyMin
	}

	// Check if the ray's X and Y distances overlap
	if (tMin > tyMax) || (tyMin > tMax) {
		return Vector3{}, false
	}
	// Update tMin and tMax to include the Y-axis range
	// Compare tMin to handle NaN from X axis
	if tyMin > tMin || tMin != tMin {
		tMin = tyMin
	}
	if tyMax < tMax || tMax != tMax {
		tMax = tyMax
	}

	// Repeat the process for the Z-axis
	tzMin := (b.Min.Z - ray.Origin.Z) / ray.Direction.Z
	tzMax := (b.Max.Z - ray.Origin.Z) / ray.Direction.Z
	// Ensure tzMin is the smaller of the two values.
	if tzMin > tzMax {
		tzMin, tzMax = tzMax, tzMin
	}

	// Check if the ray's X, Y, and Z distances overlap
	if (tMin > tzMax) || (tzMin > tMax) {
		return Vector3{}, false
	}
	// Update tMin and tMax based on Z-axis range
	// Compare tMin to handle NaN from X and Y axes
	if tzMin > tMin || tMin != tMin {
		tMin = tzMin
	}
	if tzMax < tMax || tMax != tMax {
		tMax = tzMax
	}

	// The ray intersects the box
	// Pick the closest intersection
	if tMin >= 0 {
		return ray.At(tMin), true
	}
	return ray.At(tMax), true
}

// IntersectTriangle returns if this ray intersects the triangle with the face defined by points a, b, c.
// If backfaceCulling is false it ignores the intersection if the face is not oriented in the ray direction.
func (ray *Ray) IntersectTriangle(a, b, c *Vector3, backfaceCulling bool) (Vector3, bool) {
	var diff Vector3
	var edge1 Vector3
	var edge2 Vector3
	var normal Vector3

	edge1.SubVectors(b, a)
	edge2.SubVectors(c, a)
	normal.CrossVectors(&edge1, &edge2)

	// Solve Q + t*D = b1*E1 + b2*E2 (Q = kDiff, D = ray direction,
	// E1 = kEdge1, E2 = kEdge2, N = Cross(E1,E2)) by
	//   |Dot(D,N)|*b1 = sign(Dot(D,N))*Dot(D,Cross(Q,E2))
	//   |Dot(D,N)|*b2 = sign(Dot(D,N))*Dot(D,Cross(E1,Q))
	//   |Dot(D,N)|*t = -sign(Dot(D,N))*Dot(Q,N)
	DdN := ray.Direction.Dot(&normal)
	var sign float32

	if DdN > 0 {
		if backfaceCulling {
			return Vector3{}, false
		}
		sign = 1
	} else if DdN < 0 {
		sign = -1
		DdN = -DdN
	} else {
		return Vector3{}, false
	}

	diff.SubVectors(&ray.Origin, a)
	DdQxE2 := sign * ray.Direction.Dot(edge2.CrossVectors(&diff, &edge2))

	// b1 < 0, no intersection
	if DdQxE2 < 0 {
		return Vector3{}, false
	}

	DdE1xQ := sign * ray.Direction.Dot(edge1.Cross(&diff))
	// b2 < 0, no intersection
	if DdE1xQ < 0 {
		return Vector3{}, false
	}

	// b1+b2 > 1, no intersection
	if DdQxE2+DdE1xQ > DdN {
		return Vector3{}, false
	}

	// Line intersects triangle, check if ray does.
	QdN := -sign * diff.Dot(&normal)

	// t < 0, no intersection
	if QdN < 0 {
		return Vector3{}, false
	}

	// Ray intersects triangle.
	return ray.At(QdN / DdN), true
}

// ApplyMatrix4 multiplies this ray origin and direction
// by the specified matrix4, basically transforming this ray coordinates.
func (ray *Ray) ApplyMatrix4(matrix4 *Matrix4) *Ray {
	ray.Direction.Add(&ray.Origin).ApplyMatrix4(matrix4)
	ray.Origin.ApplyMatrix4(matrix4)
	ray.Direction.Sub(&ray.Origin)
	ray.Direction.Normalize()
	return ray
}
