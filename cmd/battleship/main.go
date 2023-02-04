package main

import (
	"github.com/Jagreen1970/battleship/internal/game"
	"log"

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

	gameApi := game.NewApi(db)
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
