package winicon_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/tajtiattila/winicon"
)

func TestRead(t *testing.T) {
	// test simple icon
	testRead(t, "testdata/red-circle-64.ico", 32)

	// test icon with PNG compressed 256x256 image
	testRead(t, "testdata/red-circle.ico", 256)
}

func testRead(t *testing.T, fn string, siz int) {
	t.Logf("testRead(%q)", fn)

	f, err := os.Open(fn)
	if err != nil {
		t.Fatal("open file:", err)
	}
	defer f.Close()

	icon, err := winicon.Read(f)
	if err != nil {
		t.Fatal("read error:", err)
	}

	im := icon.FindSize(siz, siz)
	dx, dy := im.Bounds().Dx(), im.Bounds().Dy()
	if dx != siz || dy != siz {
		t.Fatalf("got image size %dx%d, want %dx%d", dx, dy, siz, siz)
	}

	const wantColor = "rgba(255,0,0,255)"
	r, g, b, a := im.At(dx/2, dy/2).RGBA()
	color := fmt.Sprintf("rgba(%d,%d,%d,%d)", r/0x101, g/0x101, b/0x101, a/0x101)
	if color != wantColor {
		t.Fatalf("got color %s, want %s", color, wantColor)
	}
}
