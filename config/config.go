package config

import (
	"flag"
	"fmt"
	"log"
	"os"

	conflib "github.com/JeremyLoy/config"
)

type config struct {
	PathToPhoto      string `config:"PHOTO_PATH"`
	Addr             string `config:"ADDR"`
	Port             string `config:"PORT"`
	MongoURL         string `config:"MONGO_URL"`
	DataBaseName     string `config:"DB_NAME"`
	VKAPIVersion     string `config:"VKAPIVersion"`
	SecretKeyForHash string `config:"SECRET_KEY_HASH"`
	JWTLifeTime      string `config:"JWT_LIFE_TIME"`
	ImagesPath       string `config:"IMAGES_PATH"`
}

var Config config

func init() {
	Config.MongoURL = "mongodb://localhost:27017"
	Config.DataBaseName = "Wikime_test"

	configPath := flag.String("configPath", "./config/.env", "Path to config file.")
	flag.StringVar(&Config.Addr, "addr", Config.Addr, "")
	flag.StringVar(&Config.Port, "port", Config.Port, "Port")

	flag.Parse()
	log.Println(os.Args, *configPath)
	err := conflib.From(*configPath).FromEnv().To(&Config)
	if err != nil {
		log.Fatal("Config read error:", err)
	}
	log.Printf("%+v\n", Config)
	flag.Parse()
	fmt.Println(Config.JWTLifeTime)
}
