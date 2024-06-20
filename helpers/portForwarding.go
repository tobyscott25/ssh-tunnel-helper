package helpers

import (
	"fmt"
	"io"
	"log"
	"net"

	"golang.org/x/crypto/ssh"
)

func StartPortForwarding(conn Connection) {
	signer, err := CustomParsePrivateKey(conn.Key)
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
