// Package winicon reads and writes 32 bit per pixel Windows ICO files.
package winicon

import "image"

// Icon represents a Windows ICO file.
// It is a set (directory) of images.
type Icon struct {
	Image []image.Image
}

func (icon *Icon) Add(m image.Image) {
	icon.Image = append(icon.Image, m)
}

// FindSize returns the image with the specified dimensions.
// If there is no exact match, a larger image is returned.
func (icon *Icon) FindSize(dx, dy int) image.Image {
	var largest, bestMatch image.Image
	for _, m := range icon.Image {
		ix, iy := m.Bounds().Dx(), m.Bounds().Dy()
		if ix == dx && iy == dy {
			return m
		}
		if ix >= dx && iy >= dy {
			if bestMatch == nil || area(bestMatch) > area(m) {
				bestMatch = m
			}
		}
		if largest == nil || area(largest) < area(m) {
			largest = m
		}
	}
	if bestMatch != nil {
		return bestMatch
	}
	return largest
}

func area(m image.Image) int64 {
	ix, iy := m.Bounds().Dx(), m.Bounds().Dy()
	return int64(ix) * int64(iy)
}
