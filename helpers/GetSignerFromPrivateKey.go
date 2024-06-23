package helpers

import (
	"fmt"
	"log"
	"os"
	"syscall"

	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

func GetSignerFromPrivateKey(keyPath string) (ssh.Signer, error) {
	key, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err == nil {
		return signer, nil
	}

	if _, ok := err.(*ssh.PassphraseMissingError); ok {

		fmt.Print("Enter passphrase: ")
		passphrase, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		if err != nil {
			log.Fatalf("unable to read passphrase: %v", err)
		}

		signer, err = ssh.ParsePrivateKeyWithPassphrase(key, passphrase)
		if err != nil {
			return nil, fmt.Errorf("unable to parse private key with passphrase: %v", err)
		}
		return signer, nil
	}

	return nil, fmt.Errorf("unable to parse private key: %v", err)
}
