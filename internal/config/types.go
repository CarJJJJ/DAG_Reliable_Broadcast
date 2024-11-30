package config

type ServerConfig struct {
	ID   string `json:"id"`
	Host string `json:"host"`
	Port string `json:"port"`
}

type Config struct {
	T       int            `json:"T"`
	N       int            `json:"N"`
	Servers []ServerConfig `json:"servers"`
}
