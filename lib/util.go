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
	cmd, err := exec.Command(path, args...).Output()
	if err != nil {
		return string(cmd), err
	}
	return string(cmd), nil

}

func ConvertToMp3(video string) {
	if existe(video + ".flv") {
		args := []string{
			"-c",
			"ffmpeg",
			"-i",
			fmt.Sprintf("%s.flv", video),
			"-b:a 192K",
			"-vn",
			fmt.Sprintf("%s.mp3", video),
		}
		_, err := Execute("/bin/sh", args)
		Printerr(err, "utils:convertToMp3:Execute "+video)
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
