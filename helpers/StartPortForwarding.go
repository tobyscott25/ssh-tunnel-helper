package helpers

import (
	"fmt"
	"log"
	"net"

	"golang.org/x/crypto/ssh"
)

func StartPortForwarding(conn PortForwardConnectionConfig, server SshServerConfig) {
	sshConn, err := EstablishSshConnection(server.User, server.Host, server.Port, server.Key)
	if err != nil {
		log.Fatalf("%v", err)
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
	serverAddress := fmt.Sprintf("%s:%d", fwd.RemoteHost, fwd.RemotePort)
	OpenSshTcpConnection(sshConn, localConn, serverAddress)
}
