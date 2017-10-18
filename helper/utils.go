package helper

import (
	"io/ioutil"
	"os"

	"github.com/Songmu/prompter"
)

// WriteFile writes buf to a file whose path is indicated by filename.
func WriteFile(filename string, buf []byte, perm os.FileMode) error {
	write := true
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		if !prompter.YN("The file already exists, do you want to overwrite it ?", false) {
			write = false
		}
	}
	if !write {
		os.Exit(0)
	}
	err := ioutil.WriteFile(filename, buf, perm)
	return err
}
