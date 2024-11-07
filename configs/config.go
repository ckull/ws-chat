package configs

import (
	"log"
	"os"
	"ws-chat/utils"

	"github.com/joho/godotenv"
)

type (
	Config struct {
		*Server
		*Db
		*Redis
		*Jwt
	}

	Server struct {
		Url  string
		Port int64
	}

	Db struct {
		URI string
	}

	Jwt struct {
		AccessTokenSecret    string
		RefreshTokenSecret   string
		ApiSecret            string
		AccessTokenDuration  int64
		RefreshTokenDuration int64
		ApiDuration          int64
	}

	Redis struct {
		Address  string
		Password string
		DB       int
	}
)

func LoadConfig(path string) *Config {
	if err := godotenv.Load(path); err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		Server: &Server{
			Url:  os.Getenv("SERVER_URL"),
			Port: utils.ParseStringToInt64(os.Getenv("SERVER_PORT")),
		},
		Db: &Db{
			URI: os.Getenv("DB_URI"),
		},
		// Jwt: &Jwt{
		// 	AccessTokenSecret:    os.Getenv("ACCESS_TOKEN_SECRET"),
		// 	RefreshTokenSecret:   os.Getenv("REFRESH_TOKEN_SECRET"),
		// 	ApiSecret:            os.Getenv("API_SECRET"),
		// 	AccessTokenDuration:  utils.ParseStringToInt64(os.Getenv("ACCESS_TOKEN_DURATION")),
		// 	RefreshTokenDuration: utils.ParseStringToInt64(os.Getenv("REFRESH_TOKEN_DURATION")),
		// 	ApiDuration:          utils.ParseStringToInt64(os.Getenv("ACCESS_TOKEN_DURATION")),
		// },
		// Redis: &Redis{
		// 	Address:  os.Getenv("REDIS_ADDRESS"),
		// 	Password: os.Getenv("REDIS_PASSWORD"),
		// 	DB:       utils.ParseStringToInt(os.Getenv("REDIS_DB")),
		// },
	}

}
