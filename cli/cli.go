package cli

import (
	"io/ioutil"
	"os"
	"os/exec"
)

var defaultEditor = "vi"

func OpenInEditor(filename string) (err error) {
	executable, err := exec.LookPath(defaultEditor)
	if err != nil {
		return
	}

	command := exec.Command(executable, filename)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	return command.Run()
}

func CaptureInput() (input []byte, err error) {
	file, err := ioutil.TempFile(os.TempDir(), "gongcommit")
	if err != nil {
		return
	}

	filename := file.Name()

	defer os.Remove(filename)

	if err = file.Close(); err != nil {
		return
	}

	if err = OpenInEditor(filename); err != nil {
		return
	}

	return ioutil.ReadFile(filename)
}
