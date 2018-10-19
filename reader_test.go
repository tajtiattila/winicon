package winicon_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/tajtiattila/winicon"
)

func TestRead(t *testing.T) {
	f, err := os.Open("testdata/red-circle-64.ico")
	if err != nil {
		log.Fatal("can't open file:", err)
	}
	defer f.Close()

	icon, err := winicon.Read(f)
	if err != nil {
		log.Fatal("read error:", err)
	}

	const siz = 32
	im := icon.FindSize(siz, siz)
	dx, dy := im.Bounds().Dx(), im.Bounds().Dy()
	if dx != siz || dy != siz {
		log.Fatalf("got image size %dx%d, want %dx%d", dx, dy, siz, siz)
	}

	const wantColor = "rgba(255,0,0,255)"
	r, g, b, a := im.At(dx/2, dy/2).RGBA()
	color := fmt.Sprintf("rgba(%d,%d,%d,%d)", r/0x101, g/0x101, b/0x101, a/0x101)
	if color != wantColor {
		log.Fatalf("got color %s, want %s", color, wantColor)
	}
}
