package tv

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
)

const (
	screenshotMaxWidth  = 960
	screenshotQuality   = 80
	screenshotMaxSource = 3840
)

func (a *ADB) Screenshot(ctx context.Context) ([]byte, error) {
	raw, err := a.run(ctx, "exec-out", "screencap", "-p")
	if err != nil {
		return nil, err
	}
	return compressScreenshot(raw)
}

func compressScreenshot(raw []byte) ([]byte, error) {
	header, err := png.DecodeConfig(bytes.NewReader(raw))
	if err != nil || header.Width > screenshotMaxSource || header.Height > screenshotMaxSource {
		return nil, fmt.Errorf("%w: unexpected screenshot size", ErrCommandFailed)
	}
	screen, err := png.Decode(bytes.NewReader(raw))
	if err != nil {
		return nil, fmt.Errorf("%w: unreadable screenshot", ErrCommandFailed)
	}
	pixels := toRGBA(screen)
	for pixels.Rect.Dx() > screenshotMaxWidth {
		pixels = halve(pixels)
	}
	var encoded bytes.Buffer
	err = jpeg.Encode(&encoded, pixels, &jpeg.Options{Quality: screenshotQuality})
	if err != nil {
		return nil, err
	}
	return encoded.Bytes(), nil
}

func toRGBA(source image.Image) *image.RGBA {
	converted, ok := source.(*image.RGBA)
	if ok {
		return converted
	}
	bounds := source.Bounds()
	target := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(target, target.Rect, source, bounds.Min, draw.Src)
	return target
}

func halve(source *image.RGBA) *image.RGBA {
	width, height := source.Rect.Dx()/2, source.Rect.Dy()/2
	target := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := range height {
		top := source.Pix[2*y*source.Stride:]
		bottom := source.Pix[(2*y+1)*source.Stride:]
		row := target.Pix[y*target.Stride:]
		for x := range width {
			for channel := range 4 {
				left, right := 8*x+channel, 8*x+4+channel
				sum := int(top[left]) + int(top[right]) + int(bottom[left]) + int(bottom[right])
				row[4*x+channel] = uint8(sum / 4)
			}
		}
	}
	return target
}
