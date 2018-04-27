package main

import (
	"flag"
	"fmt"
	"github.com/hihebark/YuxeGo/lib"
)

var (
	url, output, quality *string
	convert              *bool
)

func init() {
	url = flag.String("u", "", "URL for the Youtube video.")
	convert = flag.Bool("mp3", false, "Convert to mp3 format if set.")
	output = flag.String("o", "Downloads/YuxeGo/", "Output Folder.")
	quality = flag.String("q", "", "Quality 720,480,360,240,144 ...")
}

//Const: banner and verion of the app
const (
	BANNER  = "\033[92m  .----.\nt(\033[91m.\033[0m___\033[91m.\033[92mt) - Yuxe\n  `----\033[0m\n"
	VERSION = "0.1.0-dev"
)

func main() {

	fmt.Printf(BANNER)
	flag.Parse()
	if *url != "" {
		videodata := lib.VideoFlag{
			URL:     *url,
			Output:  *output,
			Convert: *convert,
			Quality: *quality,
		}
		lib.DownloadVideo(videodata)
	} else {
		lib.Bad("No url provided.")
		flag.PrintDefaults()
	}

}
