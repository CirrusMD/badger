package main

import (
	"log"

	"image"

	"flag"
	"strings"

	"bytes"

	"path/filepath"

	"math"

	"github.com/CirrusMD/badger/internal"
	"github.com/bmatcuk/doublestar"
	"github.com/nfnt/resize"
	"golang.org/x/image/font/inconsolata"
	gg "gopkg.in/fogleman/gg.v1"
)

const Version = "0.2.0"

var (
	mversion string
	buildNum string

	beta  bool
	alpha bool
	dark  bool

	globPattern string
)

func init() {
	log.SetFlags(0)

	flag.StringVar(&mversion, "mversion", "", "Marketing version (ex: 1.3.4)")
	flag.StringVar(&buildNum, "b", "", "Build number")

	flag.BoolVar(&beta, "beta", false, "Show beta label image in lower right corner")
	flag.BoolVar(&alpha, "alpha", false, "Show alpha label image in lower right corner")
	flag.BoolVar(&dark, "dark", false, "Show dark beta/alpha image in lower right corner. Default is a light image.")

	flag.StringVar(&globPattern, "glob", "./**/*.appiconset", "Glob pattern to icon PNGs.")
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

		if mversion != "" {
			drawTopText(parent, mversion, "#555555", 0)
		}
		if buildNum != "" {
			drawTopText(parent, buildNum, "f48041", parent.Width()/2)
		}
		overlayBadgeImage(parent)

		err = gg.SavePNG(imgPath, parent.Image())
		exitIf("could not save png", err)
	}
}

func findImages() []string {
	pattern := filepath.Join(globPattern, "*.png")
	pattern = filepath.Clean(pattern)
	images, err := doublestar.Glob(pattern)
	exitIf("glob pattern failed", err)
	if len(images) == 0 {
		log.Fatalf(`could not find any PNGs for pattern "%s"`, pattern)
	}

	return images
}

func drawTopText(parent *gg.Context, text string, hexColor string, x int) {
	temp := gg.NewContext(50, 100)
	temp.SetFontFace(inconsolata.Bold8x16)
	sw, sh := temp.MeasureString(text)

	vc := gg.NewContext(int(math.Ceil(sw)), int(math.Ceil(sh)))
	vc.SetHexColor(hexColor)
	vc.Clear()

	temp.SetFontFace(inconsolata.Bold8x16)
	vc.SetRGB(1, 1, 1)
	vc.DrawStringAnchored(text, sw/2, sh/2, 0.5, 0.5)

	img := vc.Image()
	w, h := versionDimensions(parent)
	img = resize.Resize(uint(w), uint(h), img, resize.NearestNeighbor)

	parent.DrawImage(img, x, 0)
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
