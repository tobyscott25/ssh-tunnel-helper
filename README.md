# SSH Tunnel Helper

A CLI tool to manage SSH tunnels for SOCKS proxies and port forwarding using configuration from a YAML file.

## Usage

```shell
go run main.go start

# or, to specify a config file that's not the default ($HOME/.config/ssh-tunnel-helper/config.yaml)
go run main.go start --config example-config/config.yaml
```

## Configuration

The configuration file must be in YAML format and match [this schema](config-schema.json). Here is an [example config file](example-config/config.yaml).
