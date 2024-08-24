// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gls

import "maps"

// ShaderDefines is a store of shader defines ("#define <key> <value>" in GLSL).
type ShaderDefines map[string]string

// NewShaderDefines creates and returns a pointer to a ShaderDefines object.
func NewShaderDefines() ShaderDefines {
	return make(ShaderDefines)
}

// Add adds to this ShaderDefines all the key-value pairs in the specified ShaderDefines.
func (sd ShaderDefines) Add(other ShaderDefines) {
	for k, v := range other {
		sd[k] = v
	}
}

// Equals compares two ShaderDefines and return true if they contain the same key-value pairs.
func (sd ShaderDefines) Equals(other ShaderDefines) bool {
	if sd == nil && other == nil {
		return true
	} else if sd != nil || other != nil {
		return false
	}
	return maps.Equal(sd, other)
}
