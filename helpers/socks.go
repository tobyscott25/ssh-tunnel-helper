package helpers

import (
	"fmt"
	"io"
	"log"
	"net"

	"golang.org/x/crypto/ssh"
)

func StartSOCKSTunnel(conn Connection) {
	sshConn, err := EstablishSSHConnection(conn.User, conn.Server, conn.Key)
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

		go handleSOCKSConnection(sshConn, localConn)
	}
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
