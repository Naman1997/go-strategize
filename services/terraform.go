package services

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func Terraform_init(path string, dir string) {
	cmd0 := "terraform"
	cmd1 := "-chdir=" + dir + "/" + path
	cmd2 := "init"
	cmd := exec.Command(cmd0, cmd1, cmd2)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("[ERROR] %v", err)
	}

	fmt.Println("INFO: Finished executing init stage", cmd.Stdin)
}

func Terraform_apply(path string, dir string) {
	cmd0 := "terraform"
	cmd1 := "-chdir=" + dir + "/" + path
	cmd2 := "apply"
	cmd3 := "-auto-approve"
	cmd := exec.Command(cmd0, cmd1, cmd2, cmd3)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("[ERROR] %v", err)
	}

	fmt.Println("INFO: Finished executing apply stage", cmd.Stdin)
}
