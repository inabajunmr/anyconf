package editor

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// https://qiita.com/lighttiger2505/items/d3b9ee9884c75a7819d8
func LaunchEditor(filePath string, command string) error {
	// Open text editor
	err := openEditor(command, filePath)
	if err != nil {
		return errors.New(fmt.Sprintf("Can't open %v.", command))
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
