package main

import(
	"encoding/json"
	"io/ioutil"
)

type Configuration struct{
	TelegramBotToken	string
	Webhook				Webhook
	MySQL				MySQL
}

type Webhook struct{
	Port		int
	Private		string
	Public		string
}

type MySQL struct{
	Username	string
	Password	string
	Database	string
	Host		string
	Port		string
}

func (config Configuration) generateDSN() (string){
	return config.MySQL.Username + ":" + config.MySQL.Password + "@tcp(" + config.MySQL.Host + ":" + config.MySQL.Port + ")/" + config.MySQL.Database
}

func load_config() (Configuration, error){
	var config Configuration
	raw, err := ioutil.ReadFile("config.json")
	if err != nil {
		return config, err
	}

	err1 := json.Unmarshal(raw, &config)

	return config, err1
}