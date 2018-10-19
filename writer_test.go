package winicon_test

import (
	"image"
	"image/color"
	"os"
	"testing"

	"github.com/tajtiattila/winicon"
)

func TestWrite(t *testing.T) {
	testWrite(t, "testdata/output_checkers.ico", winicon.LargePNG(true))
	testWrite(t, "testdata/output_checkers_uncompressed.ico", winicon.LargePNG(false))
	testWrite(t, "testdata/output_checkers_png.ico", winicon.PreferPNG(true))
}

func testWrite(t *testing.T, fn string, opts ...winicon.WriteOption) {
	var icon winicon.Icon
	for _, siz := range []int{256, 48, 32, 24, 16} {
		icon.Add(checkers(siz, siz/4,
			color.NRGBA{0, 255, 0, 255},
			color.NRGBA{255, 255, 0, 128}))
	}

	f, err := os.Create(fn)
	if err != nil {
		t.Fatal("open file:", err)
	}
	defer f.Close()

	if err := winicon.Write(f, icon, opts...); err != nil {
		t.Fatal("write file:", err)
	}
}

func checkers(dim, sdim int, c1, c2 color.Color) image.Image {
	m := image.NewNRGBA(image.Rect(0, 0, dim, dim))

	for y := 0; y < dim; y++ {
		yb := uint(y/sdim) & 1
		for x := 0; x < dim; x++ {
			xb := uint(x/sdim) & 1

			var c color.Color
			if xb^yb == 0 {
				c = c1
			} else {
				c = c2
			}

			m.Set(x, y, c)
		}
	}

	return m
}
