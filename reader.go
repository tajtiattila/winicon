package winicon

import (
	"bytes"
	"encoding/binary"
	"errors"
	"image"
	"image/png"
	"io"
	"io/ioutil"
)

var ErrFormat = errors.New("bad format")
var ErrUnsupported = errors.New("unsupported BMP format")

func Read(r io.Reader) (*Icon, error) {
	var header [6]byte

	p := header[:]
	if _, err := io.ReadFull(r, p); err != nil {
		return nil, err
	}

	le := binary.LittleEndian
	zero, format, nimages := le.Uint16(p[0:2]), le.Uint16(p[2:4]), int(le.Uint16(p[4:6]))
	if zero != 0 || format != 1 {
		return nil, ErrFormat
	}

	ico, err := ioutil.ReadAll(&io.LimitedReader{
		R: r,
		N: 1 << 20, // 1M
	})
	if err != nil {
		return nil, err
	}

	if len(ico) < nimages*ideSize {
		return nil, io.ErrUnexpectedEOF
	}

	icon := new(Icon)
	for i := 0; i < nimages; i++ {
		p := ico[i*ideSize:]
		siz, ofs := int(le.Uint32(p[8:12])), int(le.Uint32(p[12:16]))

		ofs -= headerSize

		if ofs < 0 || len(ico) < ofs+siz {
			// pixel data outside file
			return nil, ErrFormat
		}

		im, err := decodeImage(ico[ofs : ofs+siz])
		if err != nil {
			return nil, err
		}

		icon.Image = append(icon.Image, im)
	}

	return icon, nil
}

func decodeImage(p []byte) (image.Image, error) {
	if len(p) < 8 {
		return nil, ErrFormat
	}

	if p[0] == 0x89 && string(p[1:4]) == "PNG" {
		return png.Decode(bytes.NewReader(p))
	}

	// BITMAPINFOHEADER
	if len(p) < bihSize {
		return nil, ErrFormat
	}

	le := binary.LittleEndian

	bihs, dx, dy2 := le.Uint32(p[0:4]), int(le.Uint32(p[4:8])), int(le.Uint32(p[8:12]))
	if bihs != bihSize {
		return nil, ErrFormat
	}

	dy := dy2 / 2

	var topDown bool
	if dy < 0 {
		dy = -dy
		topDown = true
	}

	planes, bpp, compression := le.Uint16(p[12:14]), le.Uint16(p[14:16]), le.Uint32(p[16:20])
	if planes != 1 || bpp != 32 || compression != 0 {
		return nil, ErrUnsupported
	}

	nbytes := dx * dy * 4
	if len(p) < bihSize+nbytes {
		return nil, ErrFormat
	}

	m := image.NewNRGBA(image.Rect(0, 0, dx, dy))

	if dx == 0 || dy == 0 {
		return m, nil
	}

	imageRowBytes := dx * 4 // 32 bpp
	for y := 0; y < dy; y++ {
		so := bihSize + y*imageRowBytes
		var do int
		if topDown {
			do = y * m.Stride
		} else {
			do = (dy - 1 - y) * m.Stride
		}
		for x := 0; x < dx; x++ {
			b, g, r, a := p[so], p[so+1], p[so+2], p[so+3]
			m.Pix[do], m.Pix[do+1], m.Pix[do+2], m.Pix[do+3] = r, g, b, a
			so += 4
			do += 4
		}
	}

	return m, nil
}
