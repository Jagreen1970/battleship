package main

import (
	"log"

	"github.com/Jagreen1970/battleship/internal/battleship"
	"github.com/Jagreen1970/battleship/internal/database"
	"github.com/Jagreen1970/battleship/internal/http/endpoints"
	"github.com/Jagreen1970/battleship/internal/http/server"
)

func main() {
	db, err := database.New()

	if err != nil {
		log.Fatal(err)
	}
	defer func(db database.Database) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	err = db.Connect()
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	gameApi := battleship.NewApi(db)
	c := endpoints.NewController(gameApi)

	// Setup http server
	s := server.New()

	// Register endpoints
	c.Register(s)

	// Start http server
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
