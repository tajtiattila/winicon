package main

import (
	"flag"
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/tajtiattila/winicon"
)

func main() {
	destdir := flag.String("dir", "", "destination directory")
	flag.Parse()

	if *destdir == "" {
		fmt.Fprintln(os.Stderr, "-dir missing")
	}

	for _, a := range flag.Args() {
		processFile(*destdir, a)
	}
}

func processFile(destdir, fn string) {
	f, err := os.Open(fn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", fn, err)
		return
	}
	defer f.Close()

	icon, err := winicon.Read(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", fn, err)
		return
	}

	base := strings.TrimSuffix(filepath.Base(fn), filepath.Ext(fn))
	for _, im := range icon.Image {
		siz := im.Bounds().Size()
		ofn := filepath.Join(destdir, fmt.Sprintf("%s_%dx%d.png", base, siz.X, siz.Y))
		f, err := os.Create(ofn)
		if err != nil {
			fmt.Fprintf(os.Stderr, " %s: %v\n", ofn, err)
			return
		}
		defer f.Close()
		png.Encode(f, im)
	}
}
