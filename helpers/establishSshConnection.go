package helpers

import (
	"fmt"
	"log"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

func EstablishSshConnection(user string, host string, port int, keyPath string) (*ssh.Client, error) {

	serverAddress := fmt.Sprintf("%s:%d", host, port)

	// Try to connect to the ssh-agent
	sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	if err == nil {
		agentClient := agent.NewClient(sshAgent)
		config := &ssh.ClientConfig{
			User: user,
			Auth: []ssh.AuthMethod{
				ssh.PublicKeysCallback(agentClient.Signers),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Don't use this in production
		}
		client, err := ssh.Dial("tcp", serverAddress, config)
		if err == nil {
			return client, nil
		}
		log.Printf("ssh-agent connection failed: %v, falling back to private key", err)
	} else {
		log.Printf("could not connect to ssh-agent: %v, falling back to private key", err)
	}

	// Fallback to using private key
	signer, err := GetSignerFromPrivateKey(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", serverAddress, config)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to [%s]: %v", serverAddress, err)
	}
	return client, nil
}
