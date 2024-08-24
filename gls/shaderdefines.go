// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gls

import (
	"maps"
	"strconv"
)

type MaterialDefines struct {
	HAS_BASECOLORMAP      bool
	HAS_METALROUGHNESSMAP bool
	HAS_NORMALMAP         bool
	HAS_OCCLUSIONMAP      bool
	HAS_EMISSIVEMAP       bool
	BLINN                 bool
}

func (d MaterialDefines) AddToMap(defines map[string]string) {
	if d.HAS_BASECOLORMAP {
		defines["HAS_BASECOLORMAP"] = ""
	}
	if d.HAS_METALROUGHNESSMAP {
		defines["HAS_METALROUGHNESSMAP"] = ""
	}
	if d.HAS_NORMALMAP {
		defines["HAS_NORMALMAP"] = ""
	}
	if d.HAS_OCCLUSIONMAP {
		defines["HAS_OCCLUSIONMAP"] = ""
	}
	if d.HAS_EMISSIVEMAP {
		defines["HAS_EMISSIVEMAP"] = ""
	}
	if d.BLINN {
		defines["BLINN"] = ""
	}
}

type GeometryDefines struct {
	MORPHTARGETS int
}

func (d GeometryDefines) AddToMap(defines map[string]string) {
	if d.MORPHTARGETS > 0 {
		defines["MORPHTARGETS"] = strconv.Itoa(d.MORPHTARGETS)
	}
}

type GraphicDefines struct {
	BONE_INFLUENCERS int
	TOTAL_BONES      int
}

func (d GraphicDefines) AddToMap(defines map[string]string) {
	if d.BONE_INFLUENCERS > 0 {
		defines["BONE_INFLUENCERS"] = strconv.Itoa(d.BONE_INFLUENCERS)
		defines["TOTAL_BONES"] = strconv.Itoa(d.TOTAL_BONES)
	}
}

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
	} else if sd == nil || other == nil {
		return false
	}
	return maps.Equal(sd, other)
}
