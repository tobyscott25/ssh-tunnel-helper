package helpers

import (
	"io"
	"log"
	"net"

	"golang.org/x/crypto/ssh"
)

func OpenSshTcpConnection(sshConn *ssh.Client, localConn net.Conn, address string) {
	conn, err := sshConn.Dial("tcp", address)

	if err != nil {
		log.Printf("failed to establish remote connection to %s: %v", address, err)
		return
	}
	defer conn.Close()

	go func() {
		_, _ = io.Copy(conn, localConn)
	}()
	_, _ = io.Copy(localConn, conn)
}
