package config

import log "github.com/sirupsen/logrus"

func InitLogger() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
}
