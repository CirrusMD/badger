package main

import (
	"flag"
	"log"
	"strings"

	"github.com/CirrusMD/badger"
)

func main() {
	log.SetFlags(0)

	var (
		mversion string
		buildNum string

		alpha bool
		dark  bool

		globPattern string
	)

	flag.StringVar(&mversion, "mversion", "", "Marketing version (ex: 1.3.4)")
	flag.StringVar(&buildNum, "b", "", "Build number")

	flag.BoolVar(&alpha, "alpha", false, "Show alpha label image in lower right corner (default is beta image)")
	flag.BoolVar(&dark, "dark", false, "Show dark beta/alpha image in lower right corner. Default is a light image.")

	flag.StringVar(&globPattern, "glob", "./**/*.appiconset", "Glob pattern to icon PNGs.")

	flag.Parse()

	if strings.Contains(flag.Arg(0), "version") {
		log.Println(badger.Version)
		return
	}

	opts := badger.Options{
		MarketingVersion: mversion,
		BuildNumber:      buildNum,
		Alpha:            alpha,
		Dark:             dark,
		Glob:             globPattern,
	}
	err := badger.Badge(opts)
	if err != nil {
		log.Fatal(err)
	}
}
