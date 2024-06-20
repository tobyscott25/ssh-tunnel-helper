# SSH Tunnel Helper

A CLI tool to manage SSH tunnels for SOCKS proxies and port forwarding using configuration from a YAML file.

## Usage

```shell
go run main.go start

# or, to specify a config file that's not the default ($HOME/.config/ssh-tunnel-helper/config.yaml)
go run main.go start --config config.yaml
```

## Configuration

The configuration file is a YAML file. The file should contain a list of connections. Each connection should have a name, host, port, user, and key. The key should be the path to the private key file.

Example configuration file:

```yaml
connections:
  - name: "Home SOCKS Proxy"
    type: "socks"
    user: "toby"
    server: "myserver.com:22"
    key: "/path/to/your/secret_key"
    port: 1337
  - name: "Office SOCKS Proxy"
    type: "socks"
    user: "toby"
    server: "office.myserver.com:22"
    key: "/path/to/your/office_key"
    port: 1338
  - name: "Home NAS"
    type: "portforward"
    user: "toby"
    server: "myserver.com:22"
    key: "/path/to/your/id_ed25519"
    forwardings:
      - local_port: 1139
        remote_host: "nas"
        remote_port: 139
      - local_port: 1455
        remote_host: "nas"
        remote_port: 445
```
