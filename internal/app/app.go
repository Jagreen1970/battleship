package app

import "time"

func Name() string {
	return "battleship"
}

func DatabaseDriver() string {
	return "mongo"
}

func DatabaseURL() string {
	return "localhost:27017"
}

func DatabaseTimeout() time.Duration {
	return 10 * time.Second
}

func DatabaseName() string {
	return "battleship"
}

func DatabaseUser() string {
	return "root"
}

func DatabasePassword() string {
	return "battleship"
}
