package main

import (
	"flag"
	"fmt"
	"github.com/hihebark/YuxeGo/lib"
	"os"
	"strings"
)

var (
	url, output, quality string
	convert              bool
)

func init() {
	flag.StringVar(&url, "u", "", "URL for the Youtube video.")
	flag.BoolVar(&convert, "mp3", false, "Convert to mp3 format if set.")
	flag.StringVar(&output, "o", "Downloads/YuxeGo/", "Output Folder.")
	flag.StringVar(&quality, "q", "",
		"Quality: small (240p),\n\t\tmedium (360p),\n\t\tlarge (480p),\n\t\thd720 (720p),\n\t\thd1080 (1080p).")
}

//Const: banner and verion of the app
const (
	BANNER  = "\033[92m  .----.\nt(\033[91m.\033[0m___\033[91m.\033[92mt) - Yuxe\n  `----\033[0m\n"
	VERSION = "0.2.0-dev"
)

func main() {

	fmt.Printf("%s\t%s\n\n", BANNER, VERSION)
	flag.Parse()
	if len(url) != 0 {
		if !strings.Contains(url, "youtube.com") {
			lib.Bad("Must be a youtube video.")
			flag.PrintDefaults()
			os.Exit(1)
		}
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
