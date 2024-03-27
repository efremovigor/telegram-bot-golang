package config

import (
	"telegram-bot-golang/env"
)

const configParentPath = "./build/"

func GetCertPath() string {
	return configParentPath + env.GetEnvVariable("CERT_FILE")
}

func GetCertKeyPath() string {
	return configParentPath + env.GetEnvVariable("CERT_KEY")
}

func GetAppPort() string {
	return configParentPath + env.GetEnvVariable("APP_PORT")
}
