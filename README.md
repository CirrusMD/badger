# Badger

[![Build Status](https://travis-ci.org/CirrusMD/badger.svg?branch=master)](https://travis-ci.org/CirrusMD/badger)

A command line utility that adds a badge to your tvOS/iOS/Android app icon. Inspired by the Ruby gem [Badge](https://github.com/HazAT/badge) by Daniel Griesser. 

#### Example Badge:  

![alt](./examples/CirrusMD-Images.xcassets/AppIcon.appiconset/icon-76.png)
![alt](./examples/CirrusMD-Images.xcassets/AppIcon.appiconset/icon-120.png)
![alt](./examples/CirrusMD-Images.xcassets/AppIcon.appiconset/icon-180.png)

Unlike the Ruby gem, Badger has zero dependencies other than itself. Badger does not require a network connection, ImageMagick, or any sort of Ruby environment. It's a standalone executable available for macOS, linux, and windows.

The current API is not as flexible as the badge gem. Updates forthcoming.

To see a list of command line options run `badger -h` or `badger -help`

## Usage

**Warning**: *Badger modifies your icon PNGs in place.*

Options:
```
badger -h      // print help
badger -help   // print help
badger version // print current badger version
```

Flags:
```
  -alpha
    	Show alpha label image in lower right corner
  -b string
    	Build number
  -beta
    	Show beta label image in lower right corner
  -dark
    	Show dark beta/alpha image in lower right corner. Default is a light image.
  -glob string
    	Glob pattern to icon PNGs. (default "./**/*.appiconset")
  -mversion string
    	Marketing version (ex: 1.3.4)
```

## Installation

### (Recommended):

Download the latest release for your OS from [Releases](https://github.com/CirrusMD/badger/releases)


### From source:
```
go get github.com/CirrusMD/badger
cd $GOPATH/src/github.com/CirrusMD/badger
go build
```

### Using on CI

This bash one-liner will download the latest release and extract the executable to the current directory.
```
# Replace <your os> with either of these options macOS|linux|windows
curl -s -L $(curl -s https://api.github.com/repos/CirrusMD/badger/releases | grep browser_download_url | grep <your os> | head -n 1 | cut -d '"' -f 4) > badger.zip && unzip -o badger.zip

```

#### As Part of Fastlane
Here's an example to include badger in your fastlane scripts
```ruby
# Ruby
os = RbConfig::CONFIG['host_os'] =~ /darwin/ ? 'macOS' : 'linux'
`curl -s -L $(curl -s https://api.github.com/repos/CirrusMD/badger/releases | grep browser_download_url | grep #{os} | head -n 1 | cut -d '"' -f 4) > badger.zip && unzip -o badger.zip`
`./badger -beta`
```
