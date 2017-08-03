package badger

import (
	"image"

	"bytes"

	"path/filepath"

	"math"

	"fmt"

	"io"

	"io/ioutil"

	"github.com/CirrusMD/badger/internal"
	"github.com/bmatcuk/doublestar"
	"github.com/nfnt/resize"
	"golang.org/x/image/font/inconsolata"
	gg "gopkg.in/fogleman/gg.v1"
)

const Version = "0.3.2"

type Options struct {
	// Semantic version, example: 3.0.1
	MarketingVersion string
	BuildNumber      string

	// Flag to use an alpha overlay image (default is beta)
	Alpha bool

	// Flag to use a dark alpha/beta overlay image (default is light)
	Dark bool

	// Glob pattern (default is ./**/*.appiconset)
	Glob string

	Logger io.Writer
}

func (o Options) isBeta() bool {
	return !o.Alpha
}

// Badge applies badge to all images found with the glob pattern given a set of Options
func Badge(opts Options) error {
	opts = validOptions(opts)
	images, err := findImages(opts.Glob)
	if err != nil {
		return err
	}

	for _, imgPath := range images {
		fmt.Fprintf(opts.Logger, "Badging %s...\n", imgPath)
		img, err := gg.LoadImage(imgPath)
		if err != nil {
			return fmt.Errorf("could not open file: %v", err)
		}
		parent := gg.NewContextForImage(img)

		if opts.MarketingVersion != "" {
			drawTopText(parent, opts.MarketingVersion, "#555555", 0)
		}
		if opts.BuildNumber != "" {
			drawTopText(parent, opts.BuildNumber, "f48041", parent.Width()/2)
		}
		if err := overlayBadgeImage(opts, parent); err != nil {
			return err
		}
		err = gg.SavePNG(imgPath, parent.Image())
		if err != nil {
			return fmt.Errorf("could not save png: %v", err)
		}
	}
	return nil
}

func validOptions(opts Options) Options {
	if opts.Glob == "" {
		opts.Glob = "./**/*.appiconset"
	}
	if opts.Logger == nil {
		opts.Logger = ioutil.Discard
	}
	return opts
}

func findImages(glob string) ([]string, error) {
	pattern := filepath.Join(glob, "*.png")
	pattern = filepath.Clean(pattern)
	images, err := doublestar.Glob(pattern)
	if err != nil {
		return images, fmt.Errorf("unable to parse glob pattern: %s, error: %v", pattern, err)
	}
	if len(images) == 0 {
		return images, fmt.Errorf("unable to find images for pattern: %s", pattern)
	}
	return images, nil
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

func overlayBadgeImage(opts Options, parent *gg.Context) error {
	badge, err := findBadgeImage(opts)
	if err != nil {
		return err
	}
	badge = resize.Resize(uint(parent.Width()), uint(parent.Height()), badge, resize.NearestNeighbor)
	parent.DrawImage(badge, 0, 0)
	return nil
}

func findBadgeImage(opts Options) (image.Image, error) {
	imgName := ""
	if opts.Alpha && opts.Dark {
		imgName = "alpha_badge_dark.png"
	} else if opts.Alpha {
		imgName = "alpha_badge_light.png"
	} else if opts.isBeta() && opts.Dark {
		imgName = "beta_badge_dark.png"
	} else if opts.isBeta() {
		imgName = "beta_badge_light.png"
	}
	if imgName == "" {
		return nil, nil
	}
	raw, err := internal.Asset("assets/" + imgName)
	if err != nil {
		return nil, fmt.Errorf("could not load overlay image: %v", err)
	}

	img, _, err := image.Decode(bytes.NewReader(raw))
	if err != nil {
		return nil, fmt.Errorf("unable to decode overlay image: %v", err)
	}

	return img, nil
}
