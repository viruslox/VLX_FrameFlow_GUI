package config

import (
	"log"
	"os"
)

type Config struct {
	AuthUser string
	AuthPass string
}

func LoadConfig() *Config {
	user := os.Getenv("FF_GUI_USER")
	pass := os.Getenv("FF_GUI_PASS")

	if user == "" || pass == "" {
		log.Fatal("ERROR: Insecure configuration. Environment variables FF_GUI_USER and FF_GUI_PASS must be set for authentication.")
	}

	return &Config{
		AuthUser: user,
		AuthPass: pass,
	}
}
