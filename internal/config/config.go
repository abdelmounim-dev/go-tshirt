package config

type Config struct {
	DBPath        string
	ServerAddress string
}

func Load() Config {
	return Config{
		DBPath:        "data.db",
		ServerAddress: ":8080",
	}
}
