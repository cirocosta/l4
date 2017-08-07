package lib

type Server struct {
	Address string `yaml:"address"`
}

type Config struct {
	Port    int      `yaml:"port"`
	Servers []Server `yaml:"servers"`
}
