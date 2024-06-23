package helpers

type SshServerConfig struct {
	Name string `mapstructure:"name"`
	User string `mapstructure:"user"`
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Key  string `mapstructure:"key"`
}

type SocksConnectionConfig struct {
	Name            string `mapstructure:"name"`
	SshServerConfig string `mapstructure:"server"`
	Port            int    `mapstructure:"port"`
}

type PortForwardConnectionConfig struct {
	Name            string `mapstructure:"name"`
	SshServerConfig string `mapstructure:"server"`
	Forwardings     []PortForwarding
}

type PortForwarding struct {
	LocalPort  int    `mapstructure:"local_port"`
	RemoteHost string `mapstructure:"remote_host"`
	RemotePort int    `mapstructure:"remote_port"`
}

type Configuration struct {
	Servers                []SshServerConfig             `mapstructure:"servers"`
	SocksConnections       []SocksConnectionConfig       `mapstructure:"socks_connections"`
	PortForwardConnections []PortForwardConnectionConfig `mapstructure:"portforward_connections"`
}
