package utils

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type config struct {
	ServerAddr string
	RedisAddr  string
	DB         db
}

type db struct {
	Host string
	Port string
	SSL  string
	Name string
	User string
	Pass string
}

var Config config

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	serverAddr, serverAddrExist := os.LookupEnv("SERVER_ADDR")
	if !serverAddrExist {
		log.Fatal("Env SERVER_ADDR does not exist")
	}

	redisAddr, redisAddrExist := os.LookupEnv("REDIS_ADDR")
	if !redisAddrExist {
		log.Fatal("Env REDIS_ADDR does not exist")
	}

	dbHost, dbHostExist := os.LookupEnv("DB_HOST")
	if !dbHostExist {
		log.Fatal("Env DB_HOST does not exist")
	}

	dbPort, dbPortExist := os.LookupEnv("DB_PORT")
	if !dbPortExist {
		log.Fatal("Env DB_PORT does not exist")
	}

	dbSSL, dbSSLExist := os.LookupEnv("DB_SSL_MODE")
	if !dbSSLExist {
		log.Fatal("Env DB_SSL_MODE does not exist")
	}

	dbName, dbNameExist := os.LookupEnv("DB_NAME")
	if !dbNameExist {
		log.Fatal("Env DB_NAME does not exist")
	}

	dbUser, dbUserExist := os.LookupEnv("DB_USER")
	if !dbUserExist {
		log.Fatal("Env DB_USER does not exist")
	}

	dbPass, dbPassExist := os.LookupEnv("DB_PASS")
	if !dbPassExist {
		log.Fatal("Env DB_PASS does not exist")
	}

	Config.ServerAddr = serverAddr
	Config.RedisAddr = redisAddr

	Config.DB.Host = dbHost
	Config.DB.Port = dbPort
	Config.DB.SSL = dbSSL
	Config.DB.Name = dbName
	Config.DB.User = dbUser
	Config.DB.Pass = dbPass
}
