package vim

import (
	"errors"
	"os"
	"os/exec"
)

// https://qiita.com/lighttiger2505/items/d3b9ee9884c75a7819d8
func LaunchVim(filePath string) error {
	// Open text editor
	err := openEditor("vim", filePath)
	if err != nil {
		return errors.New("Can't open vim.")
	}
	return nil
}

func openEditor(program string, filePath string) error {
	c := exec.Command(program, filePath)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
