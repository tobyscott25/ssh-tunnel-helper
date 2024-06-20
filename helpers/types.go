package helpers

type PortForwarding struct {
	LocalPort  int    `mapstructure:"local_port"`
	RemoteHost string `mapstructure:"remote_host"`
	RemotePort int    `mapstructure:"remote_port"`
}

type Connection struct {
	Name        string           `mapstructure:"name"`
	Type        string           `mapstructure:"type"`
	User        string           `mapstructure:"user"`
	Server      string           `mapstructure:"server"`
	Key         string           `mapstructure:"key"`
	Port        int              `mapstructure:"port,omitempty"`
	Forwardings []PortForwarding `mapstructure:"forwardings,omitempty"`
}
