// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math32

// Sphere represents a 3D sphere defined by its center point and a radius
type Sphere struct {
	Center Vector3 // center of the sphere
	Radius float32 // radius of the sphere
}

// NewSphere creates and returns a pointer to a new sphere with
// the specified center and radius.
func NewSphere(center *Vector3, radius float32) *Sphere {
	s := new(Sphere)
	s.Center = *center
	s.Radius = radius
	return s
}

// Set sets the center and radius of this sphere.
// Returns pointer to this update sphere.
func (s *Sphere) Set(center *Vector3, radius float32) *Sphere {
	s.Center = *center
	s.Radius = radius
	return s
}

// SetFromPoints sets this sphere from the specified points array and optional center.
// Returns pointer to this update sphere.
func (s *Sphere) SetFromPoints(points []Vector3, optionalCenter *Vector3) *Sphere {
	box := Box3{}

	if optionalCenter != nil {
		s.Center.Copy(optionalCenter)
	} else {
		s.Center = box.SetFromPoints(points).Center()
	}
	var maxRadiusSq float32
	for i := 0; i < len(points); i++ {
		maxRadiusSq = Max(maxRadiusSq, s.Center.DistanceToSquared(&points[i]))
	}
	s.Radius = Sqrt(maxRadiusSq)
	return s
}

// Copy copies other sphere to this one.
// Returns pointer to this update sphere.
func (s *Sphere) Copy(other *Sphere) *Sphere {
	*s = *other
	return s
}

// Empty checks if this sphere is empty (radius <= 0)
func (s *Sphere) Empty() bool {
	return s.Radius <= 0
}

// ContainsPoint returns if this sphere contains the specified point.
func (s *Sphere) ContainsPoint(point *Vector3) bool {
	return point.DistanceToSquared(&s.Center) <= (s.Radius * s.Radius)
}

// DistanceToPoint returns the distance from the sphere surface to the specified point.
func (s *Sphere) DistanceToPoint(point *Vector3) float32 {
	return point.DistanceTo(&s.Center) - s.Radius
}

// IntersectSphere returns if other sphere intersects this one.
func (s *Sphere) IntersectSphere(other *Sphere) bool {
	radiusSum := s.Radius + other.Radius
	return other.Center.DistanceToSquared(&s.Center) <= (radiusSum * radiusSum)
}

// ClampPoint clamps the specified point inside the sphere.
// If the specified point is inside the sphere, it is the clamped point.
// Otherwise, the clamped point is the  point in the sphere surface in the nearest to the specified point.
func (s *Sphere) ClampPoint(point *Vector3) (result Vector3) {
	deltaLengthSq := s.Center.DistanceToSquared(point)
	result.Copy(point)
	if deltaLengthSq > (s.Radius * s.Radius) {
		result.Sub(&s.Center).Normalize()
		result.MultiplyScalar(s.Radius).Add(&s.Center)
	}
	return result
}

// GetBoundingBox calculates a Box3 which bounds this sphere.
func (s *Sphere) GetBoundingBox() (result Box3) {
	result = Box3{s.Center, s.Center}
	result.ExpandByScalar(s.Radius)
	return result
}

// ApplyMatrix4 applies the specified matrix transform to this sphere.
// Returns pointer to this updated sphere.
func (s *Sphere) ApplyMatrix4(matrix *Matrix4) *Sphere {
	s.Center.ApplyMatrix4(matrix)
	s.Radius = s.Radius * matrix.GetMaxScaleOnAxis()
	return s
}

// Translate translates this sphere by the specified offset.
// Returns pointer to this updated sphere.
func (s *Sphere) Translate(offset *Vector3) *Sphere {
	s.Center.Add(offset)
	return s
}
