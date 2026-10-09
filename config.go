package main

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func LoadDotEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println(err)
	}
}

func GetOptionalString(key string, def func() string) string {
	val := os.Getenv(key)
	if val == "" {
		return def()
	}
	return val
}

func GetOptionalInt(key string, def func() int) int {
	text := os.Getenv(key)
	if text == "" {
		return def()
	}

	val, err := strconv.Atoi(text)
	if err != nil {
		return def()
	}

	return val
}
