package cli

import (
	"io/ioutil"
	"os"
	"os/exec"
)

var defaultEditor = "vi"

func OpenInEditor(filename string) error {
	executable, err := exec.LookPath(defaultEditor)
	if err != nil {
		return err
	}

	command := exec.Command(executable, filename)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	return command.Run()
}

func CaptureInput() ([]byte, error) {
	var input []byte

	file, err := ioutil.TempFile(os.TempDir(), "gongcommit")
	if err != nil {
		return input, err
	}

	filename := file.Name()
	defer os.Remove(filename)

	if err := file.Close(); err != nil {
		return input, err
	}

	if err := OpenInEditor(filename); err != nil {
		return input, err
	}

	return ioutil.ReadFile(filename)
}
