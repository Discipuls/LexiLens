package bot

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
)

type config struct {
	Mode string `json:"mode"`
	Bot  struct {
		Token string `json:"token"`
	} `json:"bot"`
	MiniApp struct {
		Url string `json:"url"`
	} `json:"mini-app"`
	Database struct {
		Host         string `json:"host"`
		User         string `json:"user"`
		Password     string `json:"password"`
		DatabaseName string `json:"database-name"`
		Port         string `json:"port"`
	} `json:"database"`
	Seeker struct {
		Host string `json:"host"`
		Port string `json:"port"`
	} `json:"seeker"`
}

func GetToken() (token string, err error) {
	environmentToken := os.Getenv("BOT_TOKEN")
	if environmentToken == "" {
		configFilename := os.Getenv("PATH_TO_CONFIG")
		if configFilename == "" {
			configFilename = "config.json"
		}
		data, err1 := ioutil.ReadFile(configFilename)
		if err1 != nil {
			log.Println("Warning: Couldn't read config file: " + err.Error())
		}
		var conf config
		err = json.Unmarshal(data, &conf)
		if err != nil {
			return "", errors.New("GetToken:No env var and couldn't unmarshal config file: " + err.Error())
		}
		return conf.Bot.Token, nil
	} else {
		return environmentToken, nil
	}
}

func GetConfiguration() (conf config, err error) {
	dir, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	log.Println("Current working directory:", dir)

	configFilename := os.Getenv("PATH_TO_CONFIG")
	if configFilename == "" {
		configFilename = "config.json"
	}

	data, err := ioutil.ReadFile(configFilename)
	if err != nil {
		log.Println("Warning: Couldn't read config file: " + err.Error())
	}

	err = json.Unmarshal(data, &conf)
	if err != nil {
		log.Println("Warning: Couldn't unmarshal config file: " + err.Error())
	}

	environmentToken := os.Getenv("BOT_TOKEN")
	if environmentToken != "" {
		conf.Bot.Token = environmentToken
	}

	environmentMiniAppUrl := os.Getenv("MINI_APP_URL")
	if environmentMiniAppUrl != "" {
		conf.MiniApp.Url = environmentMiniAppUrl
	}

	environmentDBHost := os.Getenv("DB_HOST")
	if environmentDBHost != "" {
		conf.Database.Host = environmentDBHost
	}

	environmentDBUser := os.Getenv("DB_USER")
	if environmentDBUser != "" {
		conf.Database.User = environmentDBUser
	}

	environmentDBPassword := os.Getenv("DB_PASSWORD")
	if environmentDBPassword != "" {
		conf.Database.Password = environmentDBPassword
	}

	environmentDBName := os.Getenv("DB_NAME")
	if environmentDBName != "" {
		conf.Database.DatabaseName = environmentDBName
	}

	environmentDBPort := os.Getenv("DB_PORT")
	if environmentDBPort != "" {
		conf.Database.Port = environmentDBPort
	}

	environmentBotMode := os.Getenv("BOT_MODE")
	if environmentDBPort != "" {
		conf.Mode = environmentBotMode
	}

	if seekerHost := os.Getenv("SEEKER_HOST"); seekerHost != "" {
		conf.Seeker.Host = seekerHost
	}
	if seekerPort := os.Getenv("SEEKER_PORT"); seekerPort != "" {
		conf.Seeker.Port = seekerPort
	}
	return conf, nil
}
