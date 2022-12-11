package config

import (
	"flag"
	"fmt"
	"log"
	"os"

	conflib "github.com/JeremyLoy/config"
)

type config struct {
	ImagePathDisk          string `config:"IMAGES_PATH_DISK"`
	Addr                   string `config:"APP_IP"`
	Port                   string `config:"APP_PORT"`
	MongoURL               string `config:"MONGO_URL"`
	DataBaseName           string `config:"DB_NAME"`
	VKAPIVersion           string `config:"VKAPIVersion"`
	SecretKeyForHash       string `config:"SECRET_KEY_HASH"`
	JWTLifeTime            string `config:"JWT_LIFE_TIME"`
	ImagesPathURI          string `config:"IMAGES_PATH_URI"`
	MaxUploadedFileSize    int64  `config:"MAX_UPLOADED_FILE_SIZE"`
	DefaultAnimePosterPath string `config:"DEFAULT_ANIME_POSTER_PATH"`
	DefaultUserAvatarPath  string `config:"DEFAULT_USER_AVATAR_PATH"`
}

var Config config

func init() {
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
