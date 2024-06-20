package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v2"
)

type Connection struct {
	Name   string `yaml:"name"`
	User   string `yaml:"user"`
	Server string `yaml:"server"`
	Key    string `yaml:"key"`
	Port   int    `yaml:"port"`
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
	startSSHTunnel(selected)
}

func startSSHTunnel(conn Connection) {
	// Read private key
	key, err := os.ReadFile(conn.Key)
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}

	// Create the Signer for this private key
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}

	// Setup SSH client configuration
	config := &ssh.ClientConfig{
		User: conn.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Connect to SSH server
	sshConn, err := ssh.Dial("tcp", conn.Server, config)
	if err != nil {
		log.Fatalf("unable to connect to [%s]: %v", conn.Server, err)
	}
	defer sshConn.Close()

	// Create local listener
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

		go handleConnection(sshConn, localConn)
	}
}

func handleConnection(sshConn *ssh.Client, localConn net.Conn) {
	defer localConn.Close()

	// Create SOCKS tunnel
	remoteConn, err := sshConn.Dial("tcp", "localhost:1080")
	if err != nil {
		log.Printf("failed to establish remote connection: %v", err)
		return
	}
	defer remoteConn.Close()

	// Copy data between local and remote connection
	go func() {
		_, _ = io.Copy(remoteConn, localConn)
	}()
	_, _ = io.Copy(localConn, remoteConn)
}
