package config

type Config struct {
	DBUser     string
	DBName     string
	DBPassword string
	DBHost     string
	DBPort     string
}

func LoadConfig() Config {
	return Config{
		DBUser:     "shahzod",
		DBName:     "food_ordering_bot",
		DBPassword: "1",
		DBHost:     "localhost",
		DBPort:     "5432",
	}
}
