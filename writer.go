package winicon

import (
	"encoding/binary"
	"errors"
	"image"
	"image/draw"
	"io"
)

// ErrImageTooBig is reported when writing an Icon
// with an image larger than 256x256 pixels.
var ErrImageTooBig = errors.New("winicon image too big")

const headerSize = 6
const ideSize = 16
const bihSize = 40

func Write(w io.Writer, icon *Icon) error {
	for _, im := range icon.Image {
		if im.Bounds().Dx() > 256 || im.Bounds().Dy() > 256 {
			return ErrImageTooBig
		}
	}

	// header
	header := make([]byte, headerSize)

	le := binary.LittleEndian
	le.PutUint16(header[0:2], 0)
	le.PutUint16(header[2:4], 1) // ICO format
	le.PutUint16(header[4:6], uint16(len(icon.Image)))

	_, err := w.Write(header)
	if err != nil {
		return err
	}

	// icon directory
	ide := make([]byte, ideSize)
	offset := headerSize + len(icon.Image)*ideSize
	for _, im := range icon.Image {
		var dx, dy uint8
		if im.Bounds().Dx() < 256 {
			dx = uint8(im.Bounds().Dx())
		}
		if im.Bounds().Dy() < 256 {
			dy = uint8(im.Bounds().Dy())
		}
		ide[0] = dx
		ide[1] = dy
		le.PutUint16(ide[4:6], 1)  // color planes
		le.PutUint16(ide[6:8], 32) // bits per pixel
		nbytes := calcBMPsize(im)
		le.PutUint32(ide[8:12], uint32(nbytes))
		le.PutUint32(ide[12:16], uint32(offset))
		offset += nbytes
		_, err := w.Write(ide)
		if err != nil {
			return err
		}
	}

	// bitmap image data
	for _, im := range icon.Image {
		_, err := w.Write(encodeBMPbits(im))
		if err != nil {
			return err
		}
	}

	return nil
}

func calcBMPsize(im image.Image) int {
	m := asNRGBA(im)

	dx := m.Bounds().Dx()
	dy := m.Bounds().Dy()

	imageRowBytes := dx * 4 // 32 bpp
	maskRowBytes := ((dx + 31) / 32) * 4

	nbytes := imageRowBytes*dy + maskRowBytes*dy
	return bihSize + nbytes
}

func encodeBMPbits(im image.Image) []byte {
	m := asNRGBA(im)

	dx := m.Bounds().Dx()
	dy := m.Bounds().Dy()

	imageRowBytes := dx * 4 // 32 bpp
	maskRowBytes := ((dx + 31) / 32) * 4

	nbytes := imageRowBytes*dy + maskRowBytes*dy
	p := make([]byte, bihSize+nbytes)

	maskOffset := bihSize + imageRowBytes*dy

	le := binary.LittleEndian

	// BITMAPINFOHEADER
	le.PutUint32(p[0:4], bihSize)
	le.PutUint32(p[4:8], uint32(dx))       // width
	le.PutUint32(p[8:12], 2*uint32(dy))    // 2*height: image+mask
	le.PutUint16(p[12:14], 1)              // planes
	le.PutUint16(p[14:16], 32)             // bits per pixel
	le.PutUint32(p[20:24], uint32(nbytes)) // sizeImage

	for y := 0; y < dy; y++ {
		so := y * m.Stride
		yy := (dy - 1 - y) // up from bottom
		do := bihSize + yy*imageRowBytes
		mo := maskOffset + yy*maskRowBytes
		for x := 0; x < dx; x++ {
			r, g, b, a := m.Pix[so], m.Pix[so+1], m.Pix[so+2], m.Pix[so+3]
			if a > 128 {
				p[mo+x/8] |= byte(1) << uint(7-(x%8))
			}
			p[do], p[do+1], p[do+2], p[do+3] = b, g, r, a
			so += 4
			do += 4
		}
	}

	return p
}

func asNRGBA(im image.Image) *image.NRGBA {
	rgba, ok := im.(*image.NRGBA)
	if ok {
		return rgba
	}

	rgba = image.NewNRGBA(im.Bounds())
	draw.Draw(rgba, rgba.Bounds(), im, im.Bounds().Min, draw.Src)
	return rgba
}
