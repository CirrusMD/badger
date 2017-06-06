# Badger

A command line utility that adds a badge to your tvOS/iOS/Android app icon. It's inspired by the Ruby gem [Badge](https://github.com/HazAT/badge) by Daniel Griesser. 

![alt](./example_icon.png)

Unlike the ruby gem, Badger has zero dependencies (i.e. bye bye ImageMagick). Badger doesn't require a network connection either.

The current API is not as flexible as the badge gem. Updates forthcoming.

To see a list of command line options run `badger -h` or `badger -help`

**Warning**: *Badger modifies your icon PNGs in place.*

## Installation

(Recommended): Download the latest release for your OS from [Releases](https://github.com/CirrusMD/badger/releases)


From source:
```
go get github.com/CirrusMD/badger
cd $GOPATH/github.com/CirrusMD/badger
go build
```

## TODO
* Add to travis-ci
* Automated releases via travis-ci
* Add to homebrew to make CI Fastlane workflows easier. (At CirrusMD, our Fastlane scripts are currently only run on a developer's local machine.)