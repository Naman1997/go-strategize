package services

import (
	"os"
	"os/exec"
	"strings"
)

/*
AnsibleGalaxy executes ansible-galaxy collection install
using the requirements filepath passed
Expects ansible binary to be present
in PATH
*/
func AnsibleGalaxy(requirements string) {
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
		ColorPrint(ERROR, "%v", err)
	}

	ColorPrint(INFO, "Finished installing ansible requirements")
}

/*
AnsiblePlaybook executes ansible-playbook
using the playbook filepath passed
Expects ansible binary to be present
in PATH
*/
func AnsiblePlaybook(playbook string, inventory string, vars string, user string) {
	cmd0 := "ansible-playbook"
	cmd1 := playbook
	cmd2 := "-i"
	cmd3 := inventory
	cmd4 := "--user"
	cmd5 := user

	if len(vars) > 0 {
		cmd6 := "-e"
		cmd7 := "@" + vars
		command := []string{cmd0, cmd1, cmd2, cmd3, cmd4, cmd5, cmd6, cmd7}
		ColorPrint(INFO, "Executing: "+strings.Join(command, " "))
		cmd := exec.Command(cmd0, cmd1, cmd2, cmd3, cmd4, cmd5, cmd6, cmd7)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			ColorPrint(ERROR, "%v", err)
		}
	} else {
		command := []string{cmd0, cmd1, cmd2, cmd3, cmd4, cmd5}
		ColorPrint(INFO, "Executing: "+strings.Join(command, " "))
		cmd := exec.Command(cmd0, cmd1, cmd2, cmd3, cmd4, cmd5)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			ColorPrint(ERROR, "%v", err)
		}
	}
}
