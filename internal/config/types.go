package config

type ServerConfig struct {
	ID   string `json:"id"`
	Host string `json:"host"`
	Port string `json:"port"`
}

type Config struct {
	Servers []ServerConfig `json:"servers"`
}
