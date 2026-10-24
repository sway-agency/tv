package tv

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"testing"
)

func TestHalveAveragesBlocks(t *testing.T) {
	source := image.NewRGBA(image.Rect(0, 0, 4, 2))
	source.SetRGBA(0, 0, color.RGBA{R: 200, A: 255})
	source.SetRGBA(1, 1, color.RGBA{R: 100, A: 255})
	half := halve(source)
	if half.Rect.Dx() != 2 || half.Rect.Dy() != 1 {
		t.Fatalf("unexpected size %v", half.Rect)
	}
	red := half.RGBAAt(0, 0).R
	if red != 75 {
		t.Fatalf("expected red average 75, got %d", red)
	}
}

func BenchmarkHalve(b *testing.B) {
	source := image.NewRGBA(image.Rect(0, 0, 1920, 1080))
	b.ReportAllocs()
	for b.Loop() {
		halve(source)
	}
}

func TestCompressScreenshotRejectsOversizedImages(t *testing.T) {
	var encoded bytes.Buffer
	err := png.Encode(&encoded, image.NewRGBA(image.Rect(0, 0, screenshotMaxSource+1, 1)))
	if err != nil {
		t.Fatal(err)
	}
	_, err = compressScreenshot(encoded.Bytes())
	if err == nil {
		t.Fatal("oversized screenshot accepted")
	}
}

func TestCompressScreenshotShrinksToJPEG(t *testing.T) {
	var encoded bytes.Buffer
	err := png.Encode(&encoded, image.NewRGBA(image.Rect(0, 0, 1920, 1080)))
	if err != nil {
		t.Fatal(err)
	}
	jpegData, err := compressScreenshot(encoded.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	config, err := jpeg.DecodeConfig(bytes.NewReader(jpegData))
	if err != nil || config.Width != screenshotMaxWidth {
		t.Fatalf("expected a %d px JPEG, got %+v %v", screenshotMaxWidth, config, err)
	}
}
