package lib

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

//Video information
const (
	VIDINFO = "https://youtube.com/get_video_info?video_id="
)

//VideoInformation all information about video
type VideoInformation struct {
	URL       string `json:"url"`
	VideoName string `json:"videoname"`
	Data      string `json:data`
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
func DownloadVideo(urlvid string) {
	//var videoParse map[string]string
	getVideoInfo, err := GetBody(VIDINFO + getVidID(urlvid))
	if err != nil {
		fmt.Printf("net:DownloadVideo:GetBody:%s\n", err)
	}
	videoData, err := url.ParseQuery(string(getVideoInfo))
	if err != nil {
		fmt.Printf("net:DownloadVideo:ParseQuery:%s\n", err)
	}
	Good(fmt.Sprintf("Downloading: %s", SayMe(LIGHTRED, videoData["title"][0])))
	format := formatSupported(videoData["fmt_list"][0])
	
	Run(fmt.Sprintf("Format supported: %s", SayMe(LIGHTRED, strings.Join(format, ", "))))
	

}

//getVidID get video id
func getVidID(urlvid string) string {
	//https://www.youtube.com/watch?v=XXXXXXXXXXX
	urldecode := strings.Split(urlvid, "?")
	return strings.Split(urldecode[1], "=")[1]
}

//formatSupported get format supported
func formatSupported(format string) []string {
	var arrayFormat []string
	sformat := strings.Split(format, ",")
	for _, val := range sformat {
		arrayFormat = append(arrayFormat, strings.Split(val,"x")[1])
	}
	return arrayFormat
}

//parseQuery to parse query
//func parseQuery(query string, key string) string{

//	var querymap string
//	if key == "" {
//		querymap, err := url.ParseQuery(string(query))
//		if err != nil {
//			fmt.Printf("%s\n", err)
//		}
//	} else {
//		querymap, err := url.ParseQuery(string(query[key]))
//		if err != nil {
//			fmt.Printf("%s\n", err)
//		}
//	}
//	return querymap

//}
