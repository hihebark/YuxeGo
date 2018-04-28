package lib

import (
	"fmt"
	"os"
	"os/exec"
)

// Execute a shell command.
// return cmd (output) or emty string if error.
func Execute(pathExec string, args []string) (string, error) {

	fmt.Printf("%v\n", args)
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

func ConvertToMp3(video string) {
	if existe(video + ".flv") {
		args := []string{
			"-c",
			//"ffmpeg",
			//"-i",
			//fmt.Sprintf("%s.flv", video),
			//"-b:a 192K",
			//"-vn",
			fmt.Sprintf("ffmpeg -i '%s.flv' -b:a 192K -vn '%s.mp3'", video, video),
		}
		out, err := Execute("/bin/sh", args)
		Printerr(err, "utils:convertToMp3:Execute "+video)
		fmt.Printf("%v\n", out)
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
