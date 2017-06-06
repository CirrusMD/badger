package main

import (
	"log"

	"image"

	"flag"
	"strings"

	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/font/inconsolata"
	gg "gopkg.in/fogleman/gg.v1"
)

const Version = "0.1.0"

var (
	mversion string
	buildNum string
)

func init() {
	log.SetFlags(0)

	flag.StringVar(&mversion, "mversion", "", "Marketing version (ex: 1.3.4)")
	flag.StringVar(&buildNum, "b", "", "Build number")
}

func main() {
	flag.Parse()

	if strings.Contains(flag.Arg(0), "version") {
		log.Println(Version)
		return
	}

	img, err := gg.LoadImage("/users/david/desktop/icon-180.png")
	exitIf("could not open file", err)
	dc := gg.NewContextForImage(img)
	drawMarketingVersion(dc)
	drawBuildNumber(dc)

	err = gg.SavePNG("/users/david/desktop/icon-transformed.png", dc.Image())
	exitIf("could not save png", err)
}

func drawMarketingVersion(parent *gg.Context) {
	w, h := versionDimensions(parent)
	vc := gg.NewContext(w, h)
	vc.SetHexColor("#555555")
	vc.Clear()

	vc.SetFontFace(basicfont.Face7x13)
	vc.SetRGB(1, 1, 1)
	vc.DrawStringAnchored(mversion, float64(w/2), float64(h/2), 0.5, 0.5)

	parent.DrawImage(vc.Image(), 0, 0)
}

func drawBuildNumber(parent *gg.Context) {
	w, h := versionDimensions(parent)
	vc := gg.NewContext(w, h)
	vc.SetHexColor("#f48041")
	vc.Clear()

	vc.SetFontFace(basicfont.Face7x13)
	vc.SetRGB(1, 1, 1)
	vc.DrawStringAnchored(buildNum, float64(w/2), float64(h/2), 0.5, 0.5)

	parent.DrawImage(vc.Image(), parent.Width()/2, 0)
}

func versionDimensions(dc *gg.Context) (int, int) {
	width := dc.Image().Bounds().Max.X / 2
	height := int(float64(dc.Image().Bounds().Max.Y)*0.2) - 1
	return width, height
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
