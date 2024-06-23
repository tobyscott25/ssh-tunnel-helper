package helpers

import (
	"fmt"
	"log"
	"net"
	"os"
	"syscall"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/term"
)

func EstablishSSHConnection(user string, host string, port int, keyPath string) (*ssh.Client, error) {
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
		client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), config)
		if err == nil {
			return client, nil
		}
		log.Printf("ssh-agent connection failed: %v, falling back to private key", err)
	} else {
		log.Printf("could not connect to ssh-agent: %v, falling back to private key", err)
	}

	// Fallback to using private key
	signer, err := CustomParsePrivateKey(keyPath)
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
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), config)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to [%s]: %v", fmt.Sprintf("%s:%d", host, port), err)
	}
	return client, nil
}

func PromptPassphrase() []byte {
	fmt.Print("Enter passphrase: ")
	passphrase, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		log.Fatalf("unable to read passphrase: %v", err)
	}
	return passphrase
}

func CustomParsePrivateKey(keyPath string) (ssh.Signer, error) {
	key, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err == nil {
		return signer, nil
	}

	if _, ok := err.(*ssh.PassphraseMissingError); ok {
		passphrase := PromptPassphrase()
		signer, err = ssh.ParsePrivateKeyWithPassphrase(key, passphrase)
		if err != nil {
			return nil, fmt.Errorf("unable to parse private key with passphrase: %v", err)
		}
		return signer, nil
	}

	return nil, fmt.Errorf("unable to parse private key: %v", err)
}
