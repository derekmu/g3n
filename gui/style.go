// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"github.com/derekmu/g3n/gui/assets"
	"github.com/derekmu/g3n/text"
)

// Style contains the styles for all GUI elements
type Style struct {
	Font *text.Font
}

var defaultStyle = Style{
	Font: assets.NewFreeSansFont(),
}

// StyleDefault returns the current default style.
func StyleDefault() Style {
	return defaultStyle
}

// SetStyleDefault sets the current default style.
func SetStyleDefault(style Style) {
	defaultStyle = style
}
