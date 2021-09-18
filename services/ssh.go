package services

import (
	"fmt"

	"io/ioutil"
	"log"
	"strings"

	"golang.org/x/crypto/ssh"
)

type Connection struct {
	*ssh.Client
}

func ValidateConn(username string, privateKey string, homedir string, addr string, port string) {
	privateKey = HomeFix(privateKey, homedir)
	conn, err := connectInsecure(username, privateKey, homedir, addr, port)
	if err != nil {
		log.Fatalf("[ERROR] [SSH Failure] %v", err)
	}
	output, err := conn.sendCommands("echo '[INFO] Connected to' `hostname`")
	if err != nil {
		log.Fatalf("[ERROR] [SSH Failure] [%s] %v", addr, err)
	}
	fmt.Println(strings.TrimSuffix(string(output), "\n"))
}

func connectInsecure(username string, privateKey string, homedir string, addr string, port string) (*Connection, error) {
	key, err := ioutil.ReadFile(privateKey)
	if err != nil {
		log.Fatalf("[ERROR] %v", err)
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("[ERROR] %v", err)
	}

	path := "~/.ssh/known_hosts"
	path = HomeFix(path, homedir)

	sshConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	host := addr + ":" + port
	conn, err := ssh.Dial("tcp", host, sshConfig)
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
