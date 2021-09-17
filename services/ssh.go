package services

import (
	"fmt"

	"io/ioutil"
	"log"
	"os/user"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh"
)

type Connection struct {
	*ssh.Client
}

func ConnectInsecure(username string, privateKey string) (*Connection, error) {
	key, err := ioutil.ReadFile(privateKey)
	if err != nil {
		log.Fatalf("Unable to read private key: %v", err)
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}

	usr, _ := user.Current()
	dir := usr.HomeDir

	path := ".ssh/known_hosts"
	path = filepath.Join(dir, path)
	fmt.Println(path)

	sshConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", "192.168.0.105:22", sshConfig)
	if err != nil {
		return nil, err
	}

	return &Connection{conn}, nil

}

func (conn *Connection) SendCommands(cmds ...string) ([]byte, error) {
	session, err := conn.NewSession()
	if err != nil {
		log.Fatal(err)
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
		return output, fmt.Errorf("failed to execute command '%s' on server: %v", cmd, err)
	}

	return output, err
}
