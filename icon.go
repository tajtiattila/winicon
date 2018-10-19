// Package winicon reads and writes 32 bit per pixel Windows ICO files.
package winicon

import "image"

// Icon is a set of images.
type Icon []image.Image

// Len returns the number of images in the icon.
func (icon Icon) Len() int {
	return len(icon)
}

// Add adds m to the icon.
func (icon *Icon) Add(m image.Image) {
	*icon = append(*icon, m)
}

// FindSize returns the image with the specified dimensions.
//
// If there is no exact match, the next larger image is returned.
// If there is no larger image, the largest image is returned.
//
// It returns nil only if the icon is empty.
func (icon Icon) FindSize(dx, dy int) image.Image {
	var largest, bestMatch image.Image
	for _, m := range icon {
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
