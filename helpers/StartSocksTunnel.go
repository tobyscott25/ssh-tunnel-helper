package helpers

import (
	"fmt"
	"log"
	"net"

	"golang.org/x/crypto/ssh"
)

func StartSocksTunnel(conn SocksConnectionConfig, server SshServerConfig) {
	sshConn, err := EstablishSshConnection(server.User, server.Host, server.Port, server.Key)
	if err != nil {
		log.Fatalf("%v", err)
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

		go handleSocksConnection(sshConn, localConn)
	}
}

func handleSocksConnection(sshConn *ssh.Client, localConn net.Conn) {
	defer localConn.Close()
	OpenSshTcpConnection(sshConn, localConn, "localhost:1080")
}
