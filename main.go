package main

import (
	"log"

	"image"

	"golang.org/x/image/font/inconsolata"
	gg "gopkg.in/fogleman/gg.v1"
)

func main() {
	img, err := gg.LoadImage("/users/david/desktop/icon-180.png")
	exitIf("could not open file", err)

	betaText := betaContext(img)

	dc := gg.NewContextForImage(img)
	dc.DrawImage(betaText.Image(), 0, img.Bounds().Max.Y-betaText.Height())
	err = gg.SavePNG("/users/david/desktop/icon-transformed.png", dc.Image())

	exitIf("could not save png", err)
}

func betaContext(img image.Image) *gg.Context {
	betaX := img.Bounds().Max.X
	betaY := int(float64(img.Bounds().Max.Y)*0.2) - 1
	dc := gg.NewContext(betaX, betaY)
	dc.SetRGBA(0, 0, 0, 0.4)
	dc.Clear()

	dc.SetFontFace(inconsolata.Bold8x16)
	dc.SetRGB(1, 1, 1)
	text := "BETA"
	dc.DrawStringAnchored(text, float64(betaX/2), float64(betaY/2), 0.5, 0.5)

	return dc
}

func exitIf(mssg string, err error) {
	if err == nil {
		return
	}
	log.Fatalln(mssg+":", err)
}
