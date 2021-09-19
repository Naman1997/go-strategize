package services

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func Ansible_galaxy(requirements string) {
	cmd0 := "ansible-galaxy"
	cmd1 := "collection"
	cmd2 := "install"
	cmd3 := "-r"
	cmd4 := requirements
	cmd := exec.Command(cmd0, cmd1, cmd2, cmd3, cmd4)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("[ERROR] %v", err)
	}

	fmt.Println("[INFO] Finished installing ansible requirements")
}

func Ansible_playbook(playbook string, inventory string, vars string, user string) {
	cmd0 := "ansible-playbook"
	cmd1 := playbook
	cmd2 := "-i"
	cmd3 := inventory
	cmd4 := "--user"
	cmd5 := user

	if len(vars) > 0 {
		cmd6 := "-e"
		cmd7 := "@" + vars
		fmt.Println("[INFO] Executing: ", cmd0, cmd1, cmd2, cmd3, cmd4, cmd5, cmd6, cmd7)
		cmd := exec.Command(cmd0, cmd1, cmd2, cmd3, cmd4, cmd5, cmd6, cmd7)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Fatalf("[ERROR] %v", err)
		}
	} else {
		fmt.Println("[INFO] Executing: ", cmd0, cmd1, cmd2, cmd3, cmd4, cmd5)
		cmd := exec.Command(cmd0, cmd1, cmd2, cmd3, cmd4, cmd5)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Fatalf("[ERROR] %v", err)
		}
	}
}
