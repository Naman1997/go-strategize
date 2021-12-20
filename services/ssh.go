package services

import (
	"os"
	"os/exec"
)

/*
ValidateConn validates all SSH connections. Does the following:
> Connect to all VMs with known_hosts for callback
> Tries to execute : 'echo 'Connected to' `hostname`' on VMs
*/
func ValidateConn(username string, privateKey string, homedir string, addr string, port string, strict bool) {
	privateKey = HomeFix(privateKey, homedir)
	copySSHKey(username, addr, port, privateKey+".pub", strict)
	sendCommands(username, addr, port)
}

func copySSHKey(user string, addr string, port string, key string, strict bool) {
	cmd0 := "ssh-copy-id"
	cmd1 := "-i"
	cmd2 := key
	cmd3 := "-p"
	cmd4 := port
	host := user + "@" + addr

	if !strict {
		cmd5 := "-o"
		cmd6 := "StrictHostKeyChecking=no"
		cmd := exec.Command(cmd0, cmd1, cmd2, cmd3, cmd4, cmd5, cmd6, host)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			ColorPrint(ERROR, "[ssh-copy-id] %v", err)
		}
	} else {
		cmd := exec.Command(cmd0, cmd1, cmd2, cmd3, cmd4, host)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			ColorPrint(ERROR, "[ssh-copy-id] %v", err)
		}
	}

	ColorPrint(INFO, "[INFO] Finished executing ssh-copy-id "+addr)
}

func sendCommands(user string, addr string, port string) {
	cmd0 := "ssh"
	cmd1 := user + "@" + addr
	cmd2 := "-p"
	cmd3 := port
	cmd4 := "echo"
	cmd5 := "'Connected to' $(uname -n)"

	cmd := exec.Command(cmd0, cmd1, cmd2, cmd3, cmd4, cmd5)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		ColorPrint(ERROR, "[SSH Failure] %v", err)
	}

	ColorPrint(INFO, "[INFO] Finished executing ssh-copy-id "+addr)
}
