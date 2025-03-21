package file

import "os/exec"

func Move(src, dest string) error {
	return exec.Command("mv", src, dest).Run()
}
