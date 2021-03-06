package lib

import (
	"fmt"
	"os"
	"os/exec"
)

// Execute a shell command.
// return cmd (output) or emty string if error.
func Execute(pathExec string, args []string) (string, error) {

	path, err := exec.LookPath(pathExec)
	if err != nil {
		return "", err
	}
	cmd, err := exec.Command(path, args...).CombinedOutput()
	if err != nil {
		return string(cmd), err
	}
	return string(cmd), nil

}

//ConvertToMp3 convert the video given to mp3 format using ffmpeg
func ConvertToMp3(video string, ext string) {
	if existe(video + "." + ext) {
		args := []string{
			"-c",
			fmt.Sprintf("ffmpeg -i '%s.%s' -b:a 192K -vn '%s.mp3'", video, ext, video),
		}
		Good("Converting ...")
		_, err := Execute("/bin/sh", args)
		Printerr(err, "utils:convertToMp3:Execute "+video)
		Good("Done.")
	} else {
		Bad("No such a file")
	}
}

func existe(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
