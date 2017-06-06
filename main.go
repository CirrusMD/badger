package main

import (
	"log"

	"image"

	"flag"
	"strings"

	"bytes"

	"path/filepath"

	"github.com/CirrusMD/badger/internal"
	"github.com/nfnt/resize"
	"golang.org/x/image/font/basicfont"
	gg "gopkg.in/fogleman/gg.v1"
)

const Version = "0.1.1"

var (
	mversion string
	buildNum string

	beta  bool
	alpha bool
	dark  bool

	assetPath string
)

func init() {
	log.SetFlags(0)

	flag.StringVar(&mversion, "mversion", "", "Marketing version (ex: 1.3.4)")
	flag.StringVar(&buildNum, "b", "", "Build number")

	flag.BoolVar(&beta, "beta", false, "Show beta label image in lower right corner")
	flag.BoolVar(&alpha, "alpha", false, "Show alpha label image in lower right corner")
	flag.BoolVar(&dark, "dark", false, "Show dark beta/alpha image in lower right corner. Default is a light image.")

	flag.StringVar(&assetPath, "path", ".", "Path to your icon files")
}

func main() {
	flag.Parse()

	if strings.Contains(flag.Arg(0), "version") {
		log.Println(Version)
		return
	}

	for _, imgPath := range findImages() {
		log.Printf("Badging %s...", imgPath)
		img, err := gg.LoadImage(imgPath)
		exitIf("could not open file", err)
		parent := gg.NewContextForImage(img)
		drawMarketingVersion(parent)
		drawBuildNumber(parent)
		overlayBadgeImage(parent)

		err = gg.SavePNG(imgPath, parent.Image())
		exitIf("could not save png", err)
	}
}

func findImages() []string {
	path := filepath.Join(assetPath, "*.png")
	path = filepath.Clean(path)
	images, err := filepath.Glob(path)
	exitIf("could not find images", err)
	if len(images) == 0 {
		log.Fatalf(`could not find any PNGs in path "%s"`, path)
	}

	return images
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

func overlayBadgeImage(parent *gg.Context) {
	badge := findBadgeImage()
	if badge == nil {
		return
	}
	badge = resize.Resize(uint(parent.Width()), uint(parent.Height()), badge, resize.NearestNeighbor)
	parent.DrawImage(badge, 0, 0)
}

func findBadgeImage() image.Image {
	imgName := ""
	if alpha && dark {
		imgName = "alpha_badge_dark.png"
	} else if alpha {
		imgName = "alpha_badge_light.png"
	} else if beta && dark {
		imgName = "beta_badge_dark.png"
	} else if beta {
		imgName = "beta_badge_light.png"
	}
	if imgName == "" {
		return nil
	}
	raw, err := internal.Asset("assets/" + imgName)
	exitIf("could not load overlay image", err)

	img, _, err := image.Decode(bytes.NewReader(raw))
	exitIf("unable to decode overlay image", err)

	return img
}

func exitIf(mssg string, err error) {
	if err == nil {
		return
	}
	log.Fatalln(mssg+":", err)
}
