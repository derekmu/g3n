// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package shaders contains the several shaders used by the engine
package shaders

import (
	"iter"
	"log"
	"maps"
	"regexp"
	"strings"
)

var includeMap = map[string]string{}
var shaderMap = map[string]string{}
var programMap = map[string]ProgramInfo{}

// ProgramInfo contains information for a registered shader program.
type ProgramInfo struct {
	Vertex   string // Vertex shader name
	Fragment string // Fragment shader name
	Geometry string // Geometry shader name (optional)
}

// AddInclude adds a shader include to the registry.
//
// Panics if the include name is already registered.
func AddInclude(name string, source string) {
	if _, ok := includeMap[name]; ok {
		log.Panicf("shader include already added %s", name)
	}
	includeMap[name] = source
}

// AddShader add a shader to registry.
//
// Panics if the shader name is already registered.
func AddShader(name string, source string) {
	if _, ok := shaderMap[name]; ok {
		log.Panicf("shader already added %s", name)
	}
	shaderMap[name] = source
}

// AddProgram adds a program to the registry.
//
// Vertex and fragment shaders are required.
// Geometry shader is optional, an empty string means it won't be used.
//
// Panics if the program name is already registered or if any of the shaders aren't registered.
func AddProgram(name string, vertex string, fragment string, geometry string) {
	if _, ok := programMap[name]; ok {
		log.Panicf("program already added %s", name)
	}
	if _, ok := shaderMap[vertex]; !ok {
		log.Panicf("shader %s not found for program %s", vertex, name)
	}
	if _, ok := shaderMap[fragment]; !ok {
		log.Panicf("shader %s missing for program %s", fragment, name)
	}
	if geometry != "" {
		if _, ok := shaderMap[geometry]; !ok {
			log.Panicf("shader %s missing for program %s", geometry, name)
		}
	}
	programMap[name] = ProgramInfo{
		Vertex:   vertex,
		Fragment: fragment,
		Geometry: geometry,
	}
}

// Shaders returns an iterator of all the shaders and their expanded source code.
func Shaders() iter.Seq2[string, string] {
	return func(yield func(string, string) bool) {
		for name, source := range shaderMap {
			if !yield(name, expandShaderSource(name, source)) {
				return
			}
		}
	}
}

func expandShaderSource(shaderName string, shaderSource string) string {
	includesFound := map[string]bool{}
	includeRegex := regexp.MustCompile(`\s*#include\s+<([^>]*)>`)
	expanded := true
	for expanded {
		expanded = false
		lines := strings.Split(shaderSource, "\n")
		for i, line := range lines {
			if matches := includeRegex.FindStringSubmatch(line); matches != nil {
				includeName := matches[1]
				if _, ok := includesFound[includeName]; ok {
					log.Panicf("circular include %s in shader %s", includeName, shaderName)
				} else {
					includesFound[includeName] = true
				}
				if includeSource, ok := includeMap[includeName]; ok {
					lines[i] = includeSource
					expanded = true
				} else {
					log.Panicf("include %s missing in shader %s", includeName, shaderName)
				}
			}
		}
		shaderSource = strings.Join(lines, "\n")
	}
	return shaderSource
}

// Programs returns an iterator of all the registered program names and their shaders.
func Programs() iter.Seq2[string, ProgramInfo] {
	return maps.All(programMap)
}
