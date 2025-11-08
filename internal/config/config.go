package config

type Config struct {
	Address string
	DBPath  string
}

func Load() Config {
	return Config{
		Address: ":8080",
		DBPath:  "data.db",
	}
}
