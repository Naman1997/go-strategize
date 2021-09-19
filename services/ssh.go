package services

import (
	"fmt"
	"os"
	"os/exec"

	"io/ioutil"
	"log"
	"strings"

	"golang.org/x/crypto/ssh"
	kh "golang.org/x/crypto/ssh/knownhosts"
)

type Connection struct {
	*ssh.Client
}

func ValidateConn(username string, privateKey string, homedir string, addr string, port string, strict bool) {
	privateKey = HomeFix(privateKey, homedir)
	conn, err := connectInsecure(username, privateKey, homedir, addr, port, strict)
	if err != nil {
		log.Fatalf("[ERROR] [SSH Failure] %v", err)
	}
	output, err := conn.sendCommands("echo '[INFO] Connected to' `hostname`")
	if err != nil {
		log.Fatalf("[ERROR] [SSH Failure] [%s] %v", addr, err)
	}
	defer conn.Close()
	fmt.Println(strings.TrimSuffix(string(output), "\n"))
}

func connectInsecure(username string, privateKey string, homedir string, addr string, port string, strict bool) (*Connection, error) {
	key, err := ioutil.ReadFile(privateKey)
	if err != nil {
		log.Fatalf("[ERROR] %v", err)
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("[ERROR] %v", err)
	}

	//Add fingerprint to known_hosts and copy SSH key
	copySSHKey(username, addr, port, privateKey+".pub", strict)

	path := "~/.ssh/known_hosts"
	path = HomeFix(path, homedir)
	hostKeyCallback, err := kh.New(path)
	if err != nil {
		log.Fatalf("[ERROR] %v", err)
	}

	sshConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: hostKeyCallback,
	}

	conn, err := ssh.Dial("tcp", addr+":"+port, sshConfig)
	if err != nil {
		return nil, err
	}

	return &Connection{conn}, nil

}

func (conn *Connection) sendCommands(cmds ...string) ([]byte, error) {
	session, err := conn.NewSession()
	if err != nil {
		log.Fatalf("[ERROR] %v", err)
	}
	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	err = session.RequestPty("xterm", 80, 40, modes)
	if err != nil {
		return []byte{}, err
	}

	cmd := strings.Join(cmds, "; ")
	output, err := session.Output(cmd)
	if err != nil {
		return output, fmt.Errorf("[ERROR] Failed to execute command '%s' on server: %v", cmd, err)
	}

	return output, err
}

func copySSHKey(user string, addr string, port string, key string, strict bool) {
	cmd0 := "ssh-copy-id"
	cmd1 := user + "@" + addr
	cmd2 := "-i"
	cmd3 := key
	cmd4 := "-p"
	cmd5 := port

	if !strict {
		cmd6 := "-o"
		cmd7 := "StrictHostKeyChecking=no"
		cmd := exec.Command(cmd0, cmd1, cmd2, cmd3, cmd4, cmd5, cmd6, cmd7)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Fatalf("[ERROR] [ssh-copy-id] %v", err)
		}
	} else {
		cmd := exec.Command(cmd0, cmd1, cmd2, cmd3, cmd4, cmd5)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Fatalf("[ERROR] [ssh-copy-id] %v", err)
		}
	}

	fmt.Println("[INFO] Finished executing ssh-copy-id", addr)
}
