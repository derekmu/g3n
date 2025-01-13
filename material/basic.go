// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package material

// Basic is a simple material that uses the 'basic' shader.
type Basic struct {
	Material
}

// NewBasic creates a new Basic material.
func NewBasic() *Basic {
	m := new(Basic)
	m.InitMaterial()
	m.SetShader("basic")
	return m
}
