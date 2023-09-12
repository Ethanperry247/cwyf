package main

import "os"

func GetApplicationPort() string {
	return os.Getenv("APP_PORT")
}

func GetHasherKey() string {
	return os.Getenv("HASH_KEY")
}

func GetDatabasePath() string {
	return os.Getenv("DB_PATH")
}