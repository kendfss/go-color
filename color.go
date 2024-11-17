/*
	Package color implements some simple RGB/HSL color conversions for golang.

	By Brandon Thomson, Kenneth Sabalo

	Adapted from
	http://code.google.com/p/closure-library/source/browse/trunk/closure/goog/color/color.js
	and algorithms on easyrgb.com.

	To maintain accuracy between conversions we use floats in the color types.
	If you are storing lots of colors and care about memory use you might want
	to use something based on byte types instead.

	Also, color types don't verify their validity before converting. If you do
	something like RGB{10,20,30}.ToHSL() the results will be undefined. All
	values must be between 0 and 1.
*/

package color

import (
	"errors"
	"fmt"
	"image/color"
	"math/rand"

	"github.com/kendfss/oprs/math/real"
)

type RGB struct {
	R, G, B float64 // Red, Green, Blue values in [0, 1]
}

// Convert r, g, b values in the range [0, 255]^3
func (RGB) constructor(r, g, b uint8) RGB {
	return RGB{
		real.MapVal(float64(r), 0, 0xff, 0, 1),
		real.MapVal(float64(g), 0, 0xff, 0, 1),
		real.MapVal(float64(b), 0, 0xff, 0, 1),
	}
}

// Takes a string like '#123456' or 'ABCDEF' and returns an RGB
func HTMLToRGB(in string) (RGB, error) {
	if in[0] == '#' {
		in = in[1:]
	}

	if len(in) != 6 {
		return RGB{}, errors.New("Invalid string length")
	}

	var r, g, b byte
	if n, err := fmt.Sscanf(in, "%2x%2x%2x", &r, &g, &b); err != nil || n != 3 {
		return RGB{}, err
	}

	return RGB{float64(r) / 255, float64(g) / 255, float64(b) / 255}, nil
}

func (c RGB) ToHSL() HSL {
	var h, s, l float64

	r := c.R
	g := c.G
	b := c.B

	M := max(r, g, b)
	m := min(r, g, b)

	// Luminosity is the average of the max and min rgb color intensities.
	l = (M + m) / 2

	// saturation
	delta := M - m
	if delta == 0 {
		// it's gray
		return HSL{0, 0, l}
	}

	// it's not gray
	if l < 0.5 {
		s = delta / (M + m)
	} else {
		s = delta / (2 - M - m)
	}

	// hue
	r2 := (((M - r) / 6) + (delta / 2)) / delta
	g2 := (((M - g) / 6) + (delta / 2)) / delta
	b2 := (((M - b) / 6) + (delta / 2)) / delta
	switch {
	case r == M:
		h = b2 - g2
	case g == M:
		h = (1.0 / 3.0) + r2 - b2
	case b == M:
		h = (2.0 / 3.0) + g2 - r2
	}

	// fix wraparounds
	switch {
	case h < 0:
		h += 1
	case h > 1:
		h -= 1
	}

	return HSL{h, s, l}
}

// A nudge to make truncation round to nearest number instead of flooring
const delta = 1 / 512.0

func (c RGB) ToHTML() string {
	return fmt.Sprintf("%02x%02x%02x", byte((c.R+delta)*255), byte((c.G+delta)*255), byte((c.B+delta)*255))
}

func (c RGB) RGBA() (r, g, b, a uint32) {
	r = uint32(real.MapVal(c.R, 0, 1, 0, 0xffff))
	g = uint32(real.MapVal(c.G, 0, 1, 0, 0xffff))
	b = uint32(real.MapVal(c.B, 0, 1, 0, 0xffff))
	a = 0xffff
	return
}

var RGBModel color.Model = color.ModelFunc(rgbModel)

func rgbModel(c color.Color) color.Color {
	r, g, b, _ := c.RGBA()
	return RGB{
		real.MapVal(float64(r), 0, 0xffff, 0, 1),
		real.MapVal(float64(g), 0, 0xffff, 0, 1),
		real.MapVal(float64(b), 0, 0xffff, 0, 1),
	}
}

type HSL struct {
	H, S, L float64 // Hue, Saturation, Lightness values in [0, 1]
}

// Convert h, s, l values in the range [0, 255]^3
func (HSL) constructor(h, s, l uint8) HSL {
	return HSL{
		real.MapVal(float64(h), 0, 0xff, 0, 1),
		real.MapVal(float64(s), 0, 0xff, 0, 1),
		real.MapVal(float64(l), 0, 0xff, 0, 1),
	}
}

func hueToRGB(v1, v2, h float64) float64 {
	if h < 0 {
		h += 1
	}
	if h > 1 {
		h -= 1
	}
	switch {
	case 6*h < 1:
		return (v1 + (v2-v1)*6*h)
	case 2*h < 1:
		return v2
	case 3*h < 2:
		return v1 + (v2-v1)*((2.0/3.0)-h)*6
	}
	return v1
}

func (c HSL) RGBA() (r, g, b, a uint32) {
	return c.ToRGB().RGBA()
}

func (c HSL) ToRGB() RGB {
	h := c.H
	s := c.S
	l := c.L

	if s == 0 {
		// it's gray
		return RGB{l, l, l}
	}

	var v1, v2 float64
	if l < 0.5 {
		v2 = l * (1 + s)
	} else {
		v2 = (l + s) - (s * l)
	}

	v1 = 2*l - v2

	r := hueToRGB(v1, v2, h+(1.0/3.0))
	g := hueToRGB(v1, v2, h)
	b := hueToRGB(v1, v2, h-(1.0/3.0))

	return RGB{r, g, b}
}

func (c HSL) ToHTML() string {
	return c.ToRGB().ToHTML()
}

var HSLModel color.Model = color.ModelFunc(hslModel)

func hslModel(c color.Color) color.Color {
	return rgbModel(c).(RGB).ToHSL()
}

func New[T RGB | HSL](rh, gs, bl uint8) color.Color {
	switch any(new(T)).(type) {
	case *RGB:
		return RGB{}.constructor(rh, gs, bl)
	case *HSL:
		return HSL{}.constructor(rh, gs, bl)
	default:
		panic("impossible")
	}
}

func Random[T RGB | HSL]() color.Color {
	switch any(new(T)).(type) {
	case *RGB:
		return RGB{rand.Float64(), rand.Float64(), rand.Float64()}
	case *HSL:
		return HSL{rand.Float64(), rand.Float64(), rand.Float64()}
	default:
		panic("impossible")
	}
}
