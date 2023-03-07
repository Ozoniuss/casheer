package config

type Config struct {
	Server   Server
	Database Database
}

type Database struct {
	Host     string
	Port     int32
	User     string
	Name     string
	Password string
}

type Server struct {
	Address string
	Port    int32
}

func newConfig() Config {
	c := Config{}
	return c
}
