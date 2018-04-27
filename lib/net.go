package lib

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"sort"
	"strconv"
	"strings"
	"time"
)

//VideoFlag struct
type VideoFlag struct {
	URL     string
	Output  string
	Convert bool
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

type writeCounter struct {
	Size  int64 //size of the file
	Total int64 // Total # of bytes written
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
func DownloadVideo(vf VideoFlag) {

	viSlice := videoInfoSlice{}
	vidID := getVidID(strings.Split(vf.URL, "?")[1])
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
	if vf.Quality != "" {
		for _, v := range viSlice.videoInfo {
			if v.Quality == vf.Quality {
				Good(fmt.Sprintf("Downloading with the quality: %s - size: %s\n",
					vf.Quality, byteConverter(viSlice.videoInfo[0].Size)))

				getVideo(v.URL, name, v.Size, vf.Output, vf.Convert)
				break
			}
		}
	} else {
		sort.Slice(viSlice.videoInfo,
			func(i, j int) bool {
				return viSlice.videoInfo[j].Quality < viSlice.videoInfo[i].Quality
			})
		Good(fmt.Sprintf("Downloading with the quality: %s - size: %s\n",
			viSlice.videoInfo[0].Quality, byteConverter(viSlice.videoInfo[0].Size)))

		getVideo(viSlice.videoInfo[0].URL, name, viSlice.videoInfo[0].Size, vf.Output, vf.Convert)
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

}

//getVidID get video id
func getVidID(urlvid string) string {
	id, err := url.ParseQuery(string(urlvid))
	if err != nil {
		fmt.Printf("getVidID:%s %v\n", urlvid, err)
	}
	return id.Get("v")
}

func getVideo(url string, name string, size int64, output string, conv bool) {

	response, err := http.Get(url)
	if err != nil {
		Bad(fmt.Sprintf("getVideo:%v", err))
	}
	if response.StatusCode != 200 {
		Bad(fmt.Sprintf("response status: %d", response.StatusCode))
		os.Exit(0)
	}
	user, err := user.Current()
	Printerr(err, "net:getVideo:user.Current:")
	homeFolder := user.HomeDir
	outputFolder := homeFolder + "/" + output
	err = os.MkdirAll(outputFolder, 0755)
	if err != nil {
		Bad(fmt.Sprintf("getVideo:MKdirAll:%v", err))
	}
	Que(fmt.Sprintf("Output:%s", outputFolder))
	os.Remove(fmt.Sprintf("%s%s.flv", output, name))
	vfile, err := os.Create(fmt.Sprintf("%s%s.flv", outputFolder, name))
	if err != nil {
		Bad(fmt.Sprintf("os.Create:%v", err))
	}
	counter := &writeCounter{Size: size}
	if _, err := io.Copy(vfile, io.TeeReader(response.Body, counter)); err != nil {
		Bad(fmt.Sprintf("io.Copy:%v", err))
	}
	fmt.Printf("\n")
	if conv {
		ConvertToMp3(fmt.Sprintf("%s%s", outputFolder, name))
	}

}

func (wc *writeCounter) Write(p []byte) (int, error) {

	n := len(p)
	wc.Total += int64(n)
	per := wc.Total * 100 / wc.Size
	r := math.Ceil(float64(per / 2))
	fmt.Printf(" [%s%s] %d%% \r",
		SayMe(RED, strings.Repeat("#", int(r))), strings.Repeat("+", int(50-r)), per)
	return n, nil

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
