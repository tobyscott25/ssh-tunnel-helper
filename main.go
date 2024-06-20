package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"syscall"

	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
	"gopkg.in/yaml.v2"
)

type PortForwarding struct {
	LocalPort  int    `yaml:"local_port"`
	RemoteHost string `yaml:"remote_host"`
	RemotePort int    `yaml:"remote_port"`
}

type Connection struct {
	Name        string           `yaml:"name"`
	Type        string           `yaml:"type"`
	User        string           `yaml:"user"`
	Server      string           `yaml:"server"`
	Key         string           `yaml:"key"`
	Port        int              `yaml:"port,omitempty"`
	Forwardings []PortForwarding `yaml:"forwardings,omitempty"`
}

type Config struct {
	Connections []Connection `yaml:"connections"`
}

func main() {
	configFile := flag.String("config", "config.yaml", "Path to YAML configuration file")
	flag.Parse()

	configData, err := os.ReadFile(*configFile)
	if err != nil {
		log.Fatalf("unable to read config file: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		log.Fatalf("unable to parse config file: %v", err)
	}

	fmt.Println("Available connections:")
	for i, conn := range config.Connections {
		fmt.Printf("[%d] %s\n", i, conn.Name)
	}

	var choice int
	fmt.Print("Select a connection: ")
	fmt.Scan(&choice)

	if choice < 0 || choice >= len(config.Connections) {
		log.Fatalf("invalid choice")
	}

	selected := config.Connections[choice]
	if selected.Type == "socks" {
		startSOCKSTunnel(selected)
	} else if selected.Type == "portforward" {
		startPortForwarding(selected)
	} else {
		log.Fatalf("unknown connection type: %s", selected.Type)
	}
}

func promptPassphrase() []byte {
	fmt.Print("Enter passphrase: ")
	passphrase, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		log.Fatalf("unable to read passphrase: %v", err)
	}
	return passphrase
}

func parsePrivateKey(keyPath string) (ssh.Signer, error) {
	key, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err == nil {
		return signer, nil
	}

	if _, ok := err.(*ssh.PassphraseMissingError); ok {
		passphrase := promptPassphrase()
		signer, err = ssh.ParsePrivateKeyWithPassphrase(key, passphrase)
		if err != nil {
			return nil, fmt.Errorf("unable to parse private key with passphrase: %v", err)
		}
		return signer, nil
	}

	return nil, fmt.Errorf("unable to parse private key: %v", err)
}

func startSOCKSTunnel(conn Connection) {
	signer, err := parsePrivateKey(conn.Key)
	if err != nil {
		log.Fatalf("%v", err)
	}

	config := &ssh.ClientConfig{
		User: conn.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	sshConn, err := ssh.Dial("tcp", conn.Server, config)
	if err != nil {
		log.Fatalf("unable to connect to [%s]: %v", conn.Server, err)
	}
	defer sshConn.Close()

	localListener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", conn.Port))
	if err != nil {
		log.Fatalf("unable to create local listener: %v", err)
	}
	defer localListener.Close()

	log.Printf("SOCKS proxy listening on 127.0.0.1:%d", conn.Port)

	for {
		localConn, err := localListener.Accept()
		if err != nil {
			log.Printf("failed to accept local connection: %v", err)
			continue
		}

		go handleSOCKSConnection(sshConn, localConn)
	}
}

func startPortForwarding(conn Connection) {
	signer, err := parsePrivateKey(conn.Key)
	if err != nil {
		log.Fatalf("%v", err)
	}

	config := &ssh.ClientConfig{
		User: conn.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	sshConn, err := ssh.Dial("tcp", conn.Server, config)
	if err != nil {
		log.Fatalf("unable to connect to [%s]: %v", conn.Server, err)
	}
	defer sshConn.Close()

	for _, fwd := range conn.Forwardings {
		go func(fwd PortForwarding) {
			localListener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", fwd.LocalPort))
			if err != nil {
				log.Fatalf("unable to create local listener: %v", err)
			}
			defer localListener.Close()

			log.Printf("Port forwarding listening on 127.0.0.1:%d", fwd.LocalPort)

			for {
				localConn, err := localListener.Accept()
				if err != nil {
					log.Printf("failed to accept local connection: %v", err)
					continue
				}

				go handlePortForwarding(sshConn, localConn, fwd)
			}
		}(fwd)
	}

	select {} // Block forever
}

func handleSOCKSConnection(sshConn *ssh.Client, localConn net.Conn) {
	defer localConn.Close()

	remoteConn, err := sshConn.Dial("tcp", "localhost:1080")
	if err != nil {
		log.Printf("failed to establish remote connection: %v", err)
		return
	}
	defer remoteConn.Close()

	go func() {
		_, _ = io.Copy(remoteConn, localConn)
	}()
	_, _ = io.Copy(localConn, remoteConn)
}

func handlePortForwarding(sshConn *ssh.Client, localConn net.Conn, fwd PortForwarding) {
	defer localConn.Close()

	remoteConn, err := sshConn.Dial("tcp", fmt.Sprintf("%s:%d", fwd.RemoteHost, fwd.RemotePort))
	if err != nil {
		log.Printf("failed to establish remote connection to %s:%d: %v", fwd.RemoteHost, fwd.RemotePort, err)
		return
	}
	defer remoteConn.Close()

	go func() {
		_, _ = io.Copy(remoteConn, localConn)
	}()
	_, _ = io.Copy(localConn, remoteConn)
}
