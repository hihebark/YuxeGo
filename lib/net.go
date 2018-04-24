package lib

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"strconv"
	"time"
)

//VideoFlag struct
type VideoFlag struct {
	URL     string
	Output  string
	Format  string
	Quality int
}

//Video information
const (
	VIDINFO = "https://youtube.com/get_video_info?video_id="
)

type videoInfo struct {
	URL       string
	Duration  string //`json:"duration"`
	Extension string //`json:"extension"`
	Size      int64 //`json:size`
}

type videoInfoSlice struct {
	videoInfoSlice []videoInfo //`json:"url"`
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

	//var videoParse map[string]string
	vidId := getVidID(strings.Split(videoflag.URL, "?")[1])
	getVideoInfo, err := GetBody(VIDINFO + vidId)
	if err != nil {
		fmt.Printf("net:DownloadVideo:GetBody:%s\n", err)
	}
	videoData, err := url.ParseQuery(string(getVideoInfo))
	if err != nil {
		fmt.Printf("net:DownloadVideo:ParseQuery:%s\n", err)
	}
	Good(fmt.Sprintf("Downloading: %s", SayMe(LIGHTRED, videoData.Get("title"))))
	format := strings.Join(formatSupported(videoData.Get("fmt_list")), ", ")
	Run(fmt.Sprintf("Format supported: %s", SayMe(LIGHTRED, format)))
	pars, _ := url.ParseQuery(videoData["url_encoded_fmt_stream_map"][0])
	videoinfoSlice := videoInfoSlice{}
	for k, v := range pars["url"] {
		vidinfo, _ := url.ParseQuery(v)
		size, _ := strconv.ParseInt(vidinfo["clen"][0], 10, 64)
		duration, _ := time.ParseDuration(fmt.Sprintf("%ss",vidinfo["dur"][0]))
		//fmt.Printf("%s\n", fmt.Sprintf("%ss",vidinfo["dur"][0]))
		//fmt.Printf("%v - %v - %v\n", duration, duration.Minutes(), duration.String())
		vi := videoInfo{
			URL:       pars["url"][k],
			Duration:  duration.Round(time.Second).String(),
			Extension: vidinfo["mime"][0],
			Size:      size,
		}
		videoinfoSlice.videoInfoSlice = append(videoinfoSlice.videoInfoSlice, vi)

	}
//	fmt.Printf("Information about video \n")
//	for _, v := range videoinfoSlice.videoInfoSlice {
//		fmt.Printf("Duration: %s - Extension: %-10s - Size: %10s\n",
//			v.Duration, v.Extension, byteConverter(v.Size))
//	}
	getVideo(videoinfoSlice.videoInfoSlice[0].URL, videoData["title"][0])
	//	duration, _ := time.ParseDuration("336.735s")
	//	fmt.Printf("%s\n", duration)
	//	content, err = GetBody()
	//	if err != nil {
	//		Bad(fmt.Sprintf("net:DownloadVideo:GetBody%s\n", err))
	//	}

}

func parsItForMe(value string, key string) string {

	data, err := url.ParseQuery(string(value))
	if err != nil {
		fmt.Printf("getVidID:%s %v\n", value, err)
	}
	return data.Get(key)

}

//getVidID get video id
func getVidID(urlvid string) string {
	id, err := url.ParseQuery(string(urlvid))
	if err != nil {
		fmt.Printf("getVidID:%s %v\n", urlvid, err)
	}
	return id.Get("v")
}

//formatSupported get format supported
func formatSupported(format string) []string {
	var arrayFormat []string
	sformat := strings.Split(format, ",")
	for _, val := range sformat {
		arrayFormat = append(arrayFormat, strings.Split(val, "x")[1])
	}
	return arrayFormat
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
