package main

import (
	"flag"
	"fmt"
	"github.com/hihebark/YuxeGo/lib"
)

var (
	url, output, quality string
	convert              bool
)

func init() {
	flag.StringVar(&url, "u", "", "URL for the Youtube video.")
	flag.BoolVar(&convert, "mp3", false, "Convert to mp3 format if set.")
	flag.StringVar(&output, "o", "Downloads/YuxeGo/", "Output Folder.")
	flag.StringVar(&quality, "q", "", "Quality 720,480,360,240,144 ...")
}

//Const: banner and verion of the app
const (
	BANNER  = "\033[92m  .----.\nt(\033[91m.\033[0m___\033[91m.\033[92mt) - Yuxe\n  `----\033[0m\n"
	VERSION = "0.1.0-dev"
)

func main() {

	fmt.Printf(BANNER)
	flag.Parse()
	if len(url) != 0 {
		videodata := lib.VideoFlag{
			URL:     url,
			Output:  output,
			Convert: convert,
			Quality: quality,
		}
		lib.DownloadVideo(videodata)
	} else {
		lib.Bad("No url provided.")
		flag.PrintDefaults()
	}

}
