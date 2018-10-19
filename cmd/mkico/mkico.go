package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/png"
	"os"

	"github.com/tajtiattila/winicon"
)

func main() {
	iconfn := flag.String("icon", "", "icon file to create")
	flag.Parse()

	if *iconfn == "" {
		fmt.Fprintln(os.Stderr, "-icon missing")
		return
	}

	icon := new(winicon.Icon)

	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "need image argument(s) to add to icon")
		return
	}

	for _, a := range flag.Args() {
		m, err := readImage(a)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		icon.Add(m)
	}
}

func readImage(fn string) (image.Image, error) {
	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	m, _, err := image.Decode(f)
	return m, err
}
