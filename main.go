package main

import (
	"flag"
	"fmt"
	"github.com/hihebark/YuxeGo/lib"
)

var (
	url, format, output *string
	quality             *int
)

func init() {
	url = flag.String("u", "", "URL for the Youtube video.")
	format = flag.String("f", "", "Format of the output mp3,mp4,flv ...")
	output = flag.String("o", "~/Download/YuxeGo/", "Output Folder.")
	quality = flag.Int("q", 720, "Quality 720,480,360,240,144 ...")
}

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
			Format:  *format,
			Quality: *quality,
		}
		lib.DownloadVideo(videodata)
	} else {
		lib.Bad("No url provided.")
		flag.PrintDefaults()
	}

}
