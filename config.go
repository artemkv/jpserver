package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadDotEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println(err)
	}
}

func GetOptionalString(key string, def string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Printf("Could not find the value for the key '%s'. Using default value '%s'", key, def)
		return def
	}
	return val
}
