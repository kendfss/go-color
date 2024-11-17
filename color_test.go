package color

import (
	"image/color"
	"math/rand"
	"strconv"
	"testing"

	"github.com/kendfss/but"
	"github.com/kendfss/iters/slices"
	"github.com/kendfss/oprs"
	"github.com/kendfss/oprs/math/real"
)

const (
	// nTrials  = 2_000_000
	nTrials  = 2
	epsilonU = 129
	epsilonF = 10e-10
)

var (
	_ color.Color = RGB{}
	_ color.Color = HSL{}
)

func TestRGBtoHSLtoRGB(t *testing.T) {
	eq := func(l, r float64) bool {
		return real.Diff(r, l) <= epsilonF
	}
	for i := range nTrials {
		want := RGB{rand.Float64(), rand.Float64(), rand.Float64()}
		rgb := want.ToHTML()
		t.Run(rgb, func(t *testing.T) {
			have := want.ToHSL().ToRGB()
			rw, gw, bw := want.R, want.G, want.B
			rh, gh, bh := have.R, have.G, have.B

			if !eq(rh, rw) {
				t.Errorf("%2d   red: have %f, want %f, delta %f", i, rh, rw, real.Diff(rh, rw))
			}
			if !eq(gh, gw) {
				t.Errorf("%2d green: have %f, want %f, delta %f", i, gh, gw, real.Diff(gh, gw))
			}
			if !eq(bh, bw) {
				t.Errorf("%2d  blue: have %f, want %f, delta %f", i, bh, bw, real.Diff(bh, bw))
			}
		})
	}
}

func TestRGBtoRGBA(t *testing.T) {
	eq := func(l, r uint32) bool {
		return real.Diff(r, l) <= epsilonU
	}
	for i := range nTrials {
		c := RGB{rand.Float64(), rand.Float64(), rand.Float64()}
		rgb := c.ToHTML()
		t.Run(rgb, func(t *testing.T) {
			R := but.Mustv(strconv.ParseUint(rgb[:2], 16, 8))
			G := but.Mustv(strconv.ParseUint(rgb[2:4], 16, 8))
			B := but.Mustv(strconv.ParseUint(rgb[4:], 16, 8))

			rw := uint32(real.MapVal(float64(R), 0, 0xff, 0, 0xffff))
			gw := uint32(real.MapVal(float64(G), 0, 0xff, 0, 0xffff))
			bw := uint32(real.MapVal(float64(B), 0, 0xff, 0, 0xffff))
			aw := uint32(0xffff)
			rh, gh, bh, ah := c.RGBA()

			if !eq(rh, rw) {
				t.Errorf("%2d   red: have %6d, want %6d, delta %6d", i, rh, rw, real.Diff(rh, rw))
			}
			if !eq(gh, gw) {
				t.Errorf("%2d green: have %6d, want %6d, delta %6d", i, gh, gw, real.Diff(gh, gw))
			}
			if !eq(bh, bw) {
				t.Errorf("%2d  blue: have %6d, want %6d, delta %6d", i, bh, bw, real.Diff(bh, bw))
			}
			if !eq(ah, aw) {
				t.Errorf("%2d alpha: have %6d, want %6d, delta %6d", i, ah, aw, real.Diff(ah, aw))
			}
		})
	}
}

func TestHSLModel(t *testing.T) {
	eq := func(l, r uint32) bool {
		return real.Diff(r, l) <= epsilonU
	}
	errs := []uint32{}
	for i := range nTrials {
		c := color.RGBA{uint8(rand.Uint32() % 0xff), uint8(rand.Uint32() % 0xff), uint8(rand.Uint32() % 0xff), 255}
		hsl := HSLModel.Convert(c).(HSL)
		t.Run(hsl.ToHTML(), func(t *testing.T) {
			rw, gw, bw, aw := c.RGBA()
			rh, gh, bh, ah := hsl.RGBA()

			if !eq(rh, rw) {
				t.Errorf("%2d   red: have %6d, want %6d, delta %6d", i, rh, rw, real.Diff(rh, rw))
				errs = append(errs, real.Diff(rh, rw))
			}
			if !eq(gh, gw) {
				t.Errorf("%2d green: have %6d, want %6d, delta %6d", i, gh, gw, real.Diff(gh, gw))
				errs = append(errs, real.Diff(gh, gw))
			}
			if !eq(bh, bw) {
				t.Errorf("%2d  blue: have %6d, want %6d, delta %6d", i, bh, bw, real.Diff(bh, bw))
				errs = append(errs, real.Diff(bh, bw))
			}
			if !eq(ah, aw) {
				t.Errorf("%2d alpha: have %6d, want %6d, delta %6d", i, ah, aw, real.Diff(ah, aw))
				errs = append(errs, real.Diff(ah, aw))
			}
		})
	}
	if len(errs) > 0 {
		er := make([]float64, len(errs))
		for i, e := range errs {
			er[i] = float64(e) / float64(len(errs))
		}
		t.Errorf("mean %.0f, errors %6d", slices.Reduce(oprs.Add[float64], er), len(er))
	}
}
