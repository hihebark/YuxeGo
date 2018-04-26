package lib

import (
	"os/exec"
	"strings"
	"fmt"
)

//GetSizeFile get the size of specifiq file
func GetSizeFile(path string) string {
	sizefile, err := execute("/bin/sh", []string{"-c", fmt.Sprintf("ls -s %s", path)})
	Printerr(err, "utils:GetListFile:Execute"+path)
	return strings.Split(sizefile, " ")[0]
}

func execute(pathExec string, args []string) (string, error) {

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
