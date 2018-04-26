package lib

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

//VideoFlag struct
type VideoFlag struct {
	URL     string
	Output  string
	Format  string
	Quality string
}

//Video information
const (
	VIDINFO = "https://youtube.com/get_video_info?video_id="
)

type videoInfo struct {
	URL       string
	Duration  string //`json:"duration"`
	Extension string //`json:"extension"`
	Size      int64  //`json:size`
	Quality   string
}

type videoInfoSlice struct {
	videoInfo []videoInfo //`json:"url"`
}

//GetBody fetch the body
func GetBody(urlVideo string) (string, error) {

	client := &http.Client{}
	response, err := client.Get(urlVideo)
	if err != nil {
		return "", err
	}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil

}

//DownloadVideo donwload video from url
func DownloadVideo(videoflag VideoFlag) {

	viSlice := videoInfoSlice{}
	vidID := getVidID(strings.Split(videoflag.URL, "?")[1])
	getVideoInfo, err := GetBody(VIDINFO + vidID)
	if err != nil {
		fmt.Printf("net:DownloadVideo:GetBody:%s\n", err)
	}
	videoData, err := url.ParseQuery(string(getVideoInfo))
	if err != nil {
		fmt.Printf("net:DownloadVideo:ParseQuery:%s\n", err)
	}
	name := videoData.Get("title")
	Good(fmt.Sprintf("Downloading: %s", SayMe(LIGHTRED, name)))
	pars, _ := url.ParseQuery(videoData["url_encoded_fmt_stream_map"][0])
	for k, v := range pars["url"] {
		vidinfo, _ := url.ParseQuery(v)
		size, _ := strconv.ParseInt(vidinfo["clen"][0], 10, 64)
		duration, _ := time.ParseDuration(fmt.Sprintf("%ss", vidinfo["dur"][0]))
		vi := videoInfo{
			URL:       pars["url"][k],
			Duration:  duration.Round(time.Second).String(),
			Extension: vidinfo["mime"][0],
			Size:      size,
			Quality:   vidinfo["itag"][0],
		}
		viSlice.videoInfo = append(viSlice.videoInfo, vi)

	}
	if videoflag.Quality != "" {
		for _, v := range viSlice.videoInfo {
			if v.Quality == videoflag.Quality {
				fmt.Printf("Downloading with the quality: %s\n", videoflag.Quality)
				getVideo(v.URL, name)
				break
			}
		}
	} else {
		sort.Slice(viSlice.videoInfo,
			func(i, j int) bool {
				return viSlice.videoInfo[j].Quality < viSlice.videoInfo[i].Quality
			})
		fmt.Printf("Downloading with the quality: %s\n", viSlice.videoInfo[0].Quality)
		getVideo(viSlice.videoInfo[0].URL, name)
	}
	//http://blog.sorlo.com/youtube-fmt-list/
	//formatSupported
	// From here https://pastebin.com/5hDj7kLj
	// fmt=5    240p          vq=small     flv  mp3
	// fmt=18   360p          vq=medium    mp4  aac
	// fmt=34   360p          vq=medium    flv  aac
	// fmt=43   360p          vq=medium    vp8  vorbis
	// fmt=35   480p          vq=large     flv  aac
	// fmt=44   480p          vq=large     vp8  vorbis
	// fmt=22   720p          vq=hd720     mp4  aac
	// fmt=45   720p          vq=hd720     vp8  vorbis
	// fmt=37  1080p          vq=hd1080    mp4  aac
	// fmt=46  1080p          vq=hd1080    vp8  vorbis
	//getVideo(viSlice.videoInfo[0].URL, videoData["title"][0])
	//	Good("Information about video")
	//	for _, v := range viSlice.videoInfo {
	//		fmt.Printf("Duration: %s - Extension: %-10s - Size: %10s\n",
	//			v.Duration, v.Extension, byteConverter(v.Size))
	//} //maybe add scanner here to let the user choose.

}

//getVidID get video id
func getVidID(urlvid string) string {
	id, err := url.ParseQuery(string(urlvid))
	if err != nil {
		fmt.Printf("getVidID:%s %v\n", urlvid, err)
	}
	return id.Get("v")
}

func getVideo(path string, name string) {

	response, err := http.Get(path)
	if err != nil {
		Bad(fmt.Sprintf("getVideo:%v", err))
	}
	if response.StatusCode != 200 {
		Bad(fmt.Sprintf("response status: %d", response.StatusCode))
		os.Exit(0)
	}
	err = os.MkdirAll("data/", 0755)
	if err != nil {
		Bad(fmt.Sprintf("getVideo:MKdirAll:%v", err))
	}
	os.Remove(fmt.Sprintf("data/%s.flv", name))
	vfile, err := os.Create(fmt.Sprintf("data/%s.flv", name))
	if err != nil {
		Bad(fmt.Sprintf("os.Create:%v", err))
	}
	if _, err := io.Copy(vfile, response.Body); err != nil {
		Bad(fmt.Sprintf("io.Copy:%v", err))
	}

}

func byteConverter(length int64) string {
	mbyte := []string{"bytes", "KB", "MB", "GB", "TB"}
	if length == -1 {
		return "0 byte"
	}
	for _, x := range mbyte {
		if length < 1024.0 {
			return fmt.Sprintf("%3.1d %s", length, x)
		}
		length = length / 1024.0
	}
	return ""
}

//progressBar
func progressBar(size int, file string) {

	//ls -s main.go | awk '{print $1}'
	var percentage int
	style := "#"
	//var oldsize int
	var sfilenow int
	for {
		percentage = sfilenow * 100 / size
		if percentage == 100 {
			break
		}
		fmt.Printf("[%20s] %d\r", style, percentage)
		style += "#"
	}

}
