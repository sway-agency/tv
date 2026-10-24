package tv

import (
	"image"
	"image/color"
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
