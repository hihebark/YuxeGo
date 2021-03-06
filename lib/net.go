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

type vidIn struct {
	URL       string
	Duration  string //`json:"duration"`
	Extension string //`json:"extension"`
	Size      int64  //`json:size`
	Quality   string
}

type vidInSlice struct {
	vidIn []vidIn //`json:"url"`
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

//DownloadVideo download video from url
func DownloadVideo(vf VideoFlag) {

	viSlice := vidInSlice{}
	vidID := getVidID(strings.Split(vf.URL, "?")[1])
	getVideoInfo, err := GetBody(VIDINFO + vidID)
	Printerr(err, "net:DownloadVideo:GetBody:")
	videoData, err := url.ParseQuery(string(getVideoInfo))
	Printerr(err, "net:DownloadVideo:ParseQuery:")
	name := videoData.Get("title")
	Good(fmt.Sprintf("Downloading: %s", SayMe(LIGHTRED, name)))
	pars, _ := url.ParseQuery(videoData["url_encoded_fmt_stream_map"][0])
	for k, v := range pars["url"] {
		vidinfo, _ := url.ParseQuery(v)
		size, _ := strconv.ParseInt(vidinfo["clen"][0], 10, 64)
		duration, _ := time.ParseDuration(fmt.Sprintf("%ss", vidinfo["dur"][0]))
		vi := vidIn{
			URL:       pars["url"][k],
			Duration:  duration.Round(time.Second).String(),
			Extension: vidinfo["mime"][0],
			Size:      size,
			Quality:   vidinfo["itag"][0],
		}
		viSlice.vidIn = append(viSlice.vidIn, vi)
	}
	if vf.Quality != "" {
		for _, v := range viSlice.vidIn {
			if strings.Split(getQualityinfo(v.Quality), ":")[0] == vf.Quality {
				Good(fmt.Sprintf("Downloading with the quality: %s - size: %s\n",
					SayMe(LIGHTCYAN, strings.Split(getQualityinfo(v.Quality), ":")[0]),
					SayMe(LIGHTCYAN, byteConverter(viSlice.vidIn[0].Size))))
				getVideo(v.URL,
					name,
					v.Size,
					vf.Output,
					vf.Convert,
					v.Quality)
				break
			}
		}
	} else {
		sort.Slice(viSlice.vidIn, func(i, j int) bool {
			return viSlice.vidIn[j].Quality < viSlice.vidIn[i].Quality
		})
		Good(fmt.Sprintf("Downloading with the quality: %s - size: %s\n",
			SayMe(LIGHTCYAN, strings.Split(getQualityinfo(viSlice.vidIn[0].Quality), ":")[0]),
			SayMe(LIGHTCYAN, byteConverter(viSlice.vidIn[0].Size))))
		getVideo(viSlice.vidIn[0].URL,
			name,
			viSlice.vidIn[0].Size,
			vf.Output,
			vf.Convert,
			viSlice.vidIn[0].Quality)
	}

}

func getQualityinfo(fmt string) string {
	// Format supported:
	// http://blog.sorlo.com/youtube-fmt-list/
	// https://pastebin.com/5hDj7kLj
	switch fmt {
	case "5", "6":
		return "small:flv"
	case "17", "13":
		return "small:3gp"
	case "18":
		return "medium:mp4"
	case "22":
		return "hd720:mp4"
	case "34":
		return "medium:flv"
	case "35":
		return "large:flv"
	case "37":
		return "hd1080:mp4"
	case "43":
		return "medium:vp8"
	case "44":
		return "large:vp8"
	case "45":
		return "hd720:vp8"
	case "46":
		return "hd1080:vp8"
	}
	return "small:flv"
}

//getVidID get video id
func getVidID(urlvid string) string {
	id, err := url.ParseQuery(string(urlvid))
	Printerr(err, fmt.Sprintf("getVidID:%s", urlvid))
	return id.Get("v")
}

func getVideo(url string, name string, size int64, output string, conv bool, q string) {

	response, err := http.Get(url)
	Printerr(err, "net:getVideo:http.Get:")
	if response.StatusCode != 200 {
		Bad(fmt.Sprintf("Response status check you're connection: %d", response.StatusCode))
		os.Exit(0)
	}
	user, err := user.Current()
	Printerr(err, "net:getVideo:user.Current:")
	homeFolder := user.HomeDir
	outputFolder := homeFolder + "/" + output
	err = os.MkdirAll(outputFolder, 0755)
	Printerr(err, "net:getVideo:MKdirAll:")
	ext := strings.Split(getQualityinfo(q), ":")[1]
	Que(fmt.Sprintf("Output: %s", outputFolder))
	os.Remove(fmt.Sprintf("%s%s.%s", output, name, ext))
	vfile, err := os.Create(fmt.Sprintf("%s%s.%s", outputFolder, name, ext))
	Printerr(err, "net:getVideo:os.Create:")
	counter := &writeCounter{Size: size}
	if _, err := io.Copy(vfile, io.TeeReader(response.Body, counter)); err != nil {
		Printerr(err, "net:getVideo:io.Copy")
	}
	fmt.Printf("\n")
	if conv {
		ConvertToMp3(fmt.Sprintf("%s%s", outputFolder, name), ext)
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
