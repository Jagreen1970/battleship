package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/Jagreen1970/battleship/internal/app"
	"github.com/Jagreen1970/battleship/internal/game"
	"github.com/Jagreen1970/battleship/internal/storage"
)

type CLI struct {
	db            storage.Storage
	config        *app.Config
	reader        *bufio.Reader
	api           *game.API
	currentGameID string
	input         io.Reader
	output        io.Writer
}

func New(db storage.Storage, cfg *app.Config) *CLI {
	return &CLI{
		db:     db,
		config: cfg,
		reader: bufio.NewReader(os.Stdin),
		api:    game.NewApi(db),
		input:  os.Stdin,
		output: os.Stdout,
	}
}

// SetIO sets custom input and output for testing
func (c *CLI) SetIO(input io.Reader, output io.Writer) {
	c.input = input
	if input != nil {
		c.reader = bufio.NewReader(input)
	}
	c.output = output
}

func (c *CLI) Run() {
	fmt.Fprintln(c.output, "Battleship CLI Mode")
	fmt.Fprintln(c.output, "Available commands:")
	fmt.Fprintln(c.output, "  create-game <player> [name]: Create a new game with optional friendly name")
	fmt.Fprintln(c.output, "  show-games [page] [count]: List all active games (paginated)")
	fmt.Fprintln(c.output, "  show-game <game-id|name>: Show status and boards of a specific game")
	fmt.Fprintln(c.output, "  join-game <game-id|name> <player>: Join an existing game as a player")
	fmt.Fprintln(c.output, "  set-game <game-id|name>: Set the game ID for the current session")
	fmt.Fprintln(c.output, "  place-ship <player> <ship-type> <x> <y> <orientation>: Place a ship")
	fmt.Fprintln(c.output, "  start-game <player>: Start the game with the given player")
	fmt.Fprintln(c.output, "  fire <player> <x> <y>: Fire at coordinates")
	fmt.Fprintln(c.output, "  delete-game <game-id|name|all>: Delete a specific game or all games")
	fmt.Fprintln(c.output, "  exit: Exit CLI mode")
	fmt.Fprintln(c.output, "\nShip types: Battleship, Cruiser, Destroyer, Submarine")
	fmt.Fprintln(c.output, "Orientation: Horizontal, Vertical")

	for {
		fmt.Fprint(c.output, "> ")
		input, _ := c.reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "exit" {
			break
		}

		c.handleCommand(input)
	}
}

func (c *CLI) handleCommand(input string) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return
	}

	cmd := parts[0]
	args := parts[1:]

	switch cmd {
	case "create-game":
		if len(args) < 1 {
			fmt.Fprintln(c.output, "Usage: create-game <player> [name]")
			return
		}

		playerName := args[0]
		var gameName string
		if len(args) > 1 {
			gameName = args[1]
		}

		c.createGame(playerName, gameName)

	case "join-game":
		if len(args) < 2 {
			fmt.Fprintln(c.output, "Usage: join-game <game-id> <player>")
			return
		}
		c.joinGame(args[0], args[1])

	case "set-game":
		if len(args) < 1 {
			fmt.Fprintln(c.output, "Usage: set-game <game-id>")
			return
		}
		c.setGame(args[0])

	case "start-game":
		if c.currentGameID == "" {
			fmt.Fprintln(c.output, "Error: No game ID set. Use set-game first or specify a game ID.")
			return
		}
		if len(args) < 1 {
			fmt.Fprintln(c.output, "Usage: start-game <player>")
			return
		}
		c.startGame(c.currentGameID, args[0])

	case "show-games":
		page := 0
		count := 10
		if len(args) > 0 {
			var err error
			page, err = strconv.Atoi(args[0])
			if err != nil {
				fmt.Fprintln(c.output, "Invalid page number")
				return
			}
		}
		if len(args) > 1 {
			var err error
			count, err = strconv.Atoi(args[1])
			if err != nil {
				fmt.Fprintln(c.output, "Invalid count number")
				return
			}
		}
		c.showGames(page, count)

	case "delete-game":
		if len(args) < 1 {
			fmt.Fprintln(c.output, "Usage: delete-game <game-id|all>")
			return
		}

		if args[0] == "all" {
			c.deleteAllGames()
		} else {
			c.deleteGame(args[0])
		}

	case "show-game":
		if len(args) < 1 && c.currentGameID == "" {
			fmt.Fprintln(c.output, "Usage: show-game <game-id> or set a game ID with set-game")
			return
		}
		gameID := c.currentGameID
		if len(args) > 0 {
			gameID = args[0]
		}
		c.showGame(gameID)

	case "place-ship":
		// First check if we have all the required parameters
		if len(args) < 5 {
			fmt.Fprintln(c.output, "Usage: place-ship <player> <ship-type> <x> <y> <orientation>")
			fmt.Fprintln(c.output, "Orientation must be either 'Horizontal' or 'Vertical'")
			return
		}

		// Then check if we have a game ID
		if c.currentGameID == "" {
			fmt.Fprintln(c.output, "You must set a game ID first with set-game or include it in the command")
			return
		}

		// Parse the coordinates
		x, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Fprintln(c.output, "Invalid x coordinate")
			return
		}

		y, err := strconv.Atoi(args[3])
		if err != nil {
			fmt.Fprintln(c.output, "Invalid y coordinate")
			return
		}

		// Validate orientation
		orientation := args[4]
		if orientation != "Horizontal" && orientation != "Vertical" {
			fmt.Fprintln(c.output, "Invalid orientation. Must be either 'Horizontal' or 'Vertical'")
			return
		}

		c.placeShip(c.currentGameID, args[0], args[1], x, y, orientation)

	case "fire":
		if c.currentGameID == "" && len(args) < 3 {
			fmt.Fprintln(c.output, "Usage: fire <player> <x> <y>")
			fmt.Fprintln(c.output, "You must set a game ID first with set-game or include it in the command")
			return
		}

		x, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Fprintln(c.output, "Invalid x coordinate")
			return
		}

		y, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Fprintln(c.output, "Invalid y coordinate")
			return
		}

		c.fire(c.currentGameID, args[0], x, y)

	default:
		fmt.Fprintf(c.output, "Unknown command: %s\n", cmd)
	}
}

func (c *CLI) createGame(playerName, gameName string) {
	// First ensure the player exists
	player, err := c.api.NewPlayer(playerName)
	if err != nil {
		fmt.Fprintf(c.output, "Error creating player: %v\n", err)
		return
	}

	g, err := c.api.NewGame(player.Name, gameName)
	if err != nil {
		fmt.Fprintf(c.output, "Error creating game: %v\n", err)
		return
	}

	if gameName != "" {
		fmt.Fprintf(c.output, "Created new game '%s' with ID: %s\n", gameName, g.ID)
	} else {
		fmt.Fprintf(c.output, "Created new game with ID: %s\n", g.ID)
	}

	c.currentGameID = g.ID
}

// getGameByIDOrName tries to find a game by ID first, then by name if that fails
func (c *CLI) getGameByIDOrName(idOrName string) (*game.Game, error) {
	// Try by ID first
	g, err := c.api.GetGame(idOrName)
	if err == nil {
		return g, nil
	}

	// If not found and not an invalid ID error, return the error
	if !errors.Is(err, game.ErrorNotFound) && !strings.Contains(err.Error(), "invalid game ID") {
		return nil, err
	}

	// Try by name
	return c.api.GetGameByName(idOrName)
}

func (c *CLI) joinGame(gameIDOrName, playerName string) {
	g, err := c.getGameByIDOrName(gameIDOrName)
	if err != nil {
		fmt.Fprintf(c.output, "Error getting game: %v\n", err)
		return
	}

	// First ensure the player exists
	player, err := c.api.NewPlayer(playerName)
	if err != nil {
		fmt.Fprintf(c.output, "Error getting player: %v\n", err)
		return
	}

	err = g.Join(player)
	if err != nil {
		fmt.Fprintf(c.output, "Error joining game: %v\n", err)
		return
	}

	_, err = c.api.UpdateGame(g)
	if err != nil {
		fmt.Fprintf(c.output, "Error updating game: %v\n", err)
		return
	}

	gameName := ""
	if g.Name != "" {
		gameName = fmt.Sprintf("'%s'", g.Name)
	} else {
		gameName = g.ID
	}

	fmt.Fprintf(c.output, "Player %s joined game %s\n", playerName, gameName)
	c.currentGameID = g.ID
}

func (c *CLI) setGame(gameIDOrName string) {
	// Verify the game exists
	g, err := c.getGameByIDOrName(gameIDOrName)
	if err != nil {
		fmt.Fprintf(c.output, "Error: Game %s not found\n", gameIDOrName)
		return
	}

	c.currentGameID = g.ID

	displayName := g.ID
	if g.Name != "" {
		displayName = fmt.Sprintf("%s (ID: %s)", g.Name, g.ID)
	}

	fmt.Fprintf(c.output, "Current game set to: %s\n", displayName)
}

func (c *CLI) startGame(gameID, playerName string) {
	g, err := c.api.GetGame(gameID)
	if err != nil {
		fmt.Fprintf(c.output, "Error getting game: %v\n", err)
		return
	}

	err = g.Start(playerName)
	if err != nil {
		fmt.Fprintf(c.output, "Error starting game: %v\n", err)
		return
	}

	_, err = c.api.UpdateGame(g)
	if err != nil {
		fmt.Fprintf(c.output, "Error updating game: %v\n", err)
		return
	}

	fmt.Fprintf(c.output, "Game %s started! Player to move: %s\n", gameID, g.PlayerToMove)
}

func (c *CLI) placeShip(gameID, playerName, shipTypeStr string, x, y int, orientationStr string) {
	g, err := c.api.GetGame(gameID)
	if err != nil {
		fmt.Fprintf(c.output, "Error getting game: %v\n", err)
		return
	}

	shipType := game.ShipType(shipTypeStr)
	orientation := game.ShipOrientation(orientationStr)

	err = g.PlaceShip(playerName, shipType, x, y, orientation)
	if err != nil {
		fmt.Fprintf(c.output, "Error placing ship: %v\n", err)
		return
	}

	// Check if we can automatically start the game
	if g.Status == game.StatusSetup {
		// Try to start the game - if it's not ready yet (not all pins placed), it will return an error
		err = g.Start(playerName)
		if err == nil {
			fmt.Fprintf(c.output, "All ships placed. Game automatically started! Player to move: %s\n", g.PlayerToMove)
		}
	}

	_, err = c.api.UpdateGame(g)
	if err != nil {
		fmt.Fprintf(c.output, "Error updating game: %v\n", err)
		return
	}

	fmt.Fprintf(c.output, "Placed %s at (%d,%d) %s for player %s\n", shipType, x, y, orientation, playerName)
}

func (c *CLI) fire(gameID, playerName string, x, y int) {
	g, err := c.api.GetGame(gameID)
	if err != nil {
		fmt.Fprintf(c.output, "Error getting game: %v\n", err)
		return
	}

	move := game.Move{
		Player: playerName,
		X:      x,
		Y:      y,
	}

	err = g.MakeMove(move)
	if err != nil {
		fmt.Fprintf(c.output, "Error making move: %v\n", err)
		return
	}

	// Check hit or miss
	lastMove := g.History[len(g.History)-1]
	hitStatus := "miss"
	if lastMove.Hit {
		hitStatus = "hit"
	}

	_, err = c.api.UpdateGame(g)
	if err != nil {
		fmt.Fprintf(c.output, "Error updating game: %v\n", err)
		return
	}

	fmt.Fprintf(c.output, "Player %s fired at (%d,%d) - %s\n", playerName, x, y, hitStatus)

	// Check if game is over
	switch g.Status {
	case game.StatusWon:
		fmt.Fprintf(c.output, "Game over! %s won!\n", playerName)
	case game.StatusLost:
		fmt.Fprintf(c.output, "Game over! %s lost!\n", playerName)
	default:
		fmt.Fprintf(c.output, "Next player to move: %s\n", g.PlayerToMove)
	}
}

func (c *CLI) deleteGame(gameIDOrName string) {
	// Check if the game exists first
	g, err := c.getGameByIDOrName(gameIDOrName)
	if err != nil {
		fmt.Fprintf(c.output, "Error: Game %s not found\n", gameIDOrName)
		return
	}

	// Display the name for feedback
	displayName := g.ID
	if g.Name != "" {
		displayName = fmt.Sprintf("'%s' (ID: %s)", g.Name, g.ID)
	}

	// Delete the game
	err = c.api.DeleteGame(g.ID)
	if err != nil {
		fmt.Fprintf(c.output, "Error deleting game: %v\n", err)
		return
	}

	fmt.Fprintf(c.output, "Game %s successfully deleted\n", displayName)

	// If we deleted the current game, clear the current game ID
	if g.ID == c.currentGameID {
		c.currentGameID = ""
		fmt.Fprintln(c.output, "Current game unset")
	}
}

func (c *CLI) deleteAllGames() {
	// Prompt for confirmation
	fmt.Fprintln(c.output, "WARNING: This will delete ALL games. This action cannot be undone.")
	fmt.Fprintln(c.output, "Type 'confirm' to proceed:")

	input, _ := c.reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input != "confirm" {
		fmt.Fprintln(c.output, "Operation cancelled")
		return
	}

	// Delete all games
	count, err := c.api.DeleteAllGames()
	if err != nil {
		fmt.Fprintf(c.output, "Error deleting all games: %v\n", err)
		return
	}

	fmt.Fprintf(c.output, "Successfully deleted %d games\n", count)

	// Clear the current game ID
	if c.currentGameID != "" {
		c.currentGameID = ""
		fmt.Fprintln(c.output, "Current game unset")
	}
}

func (c *CLI) showGames(page, count int) {
	games, err := c.api.Games(page, count)
	if err != nil {
		fmt.Fprintf(c.output, "Error retrieving games: %v\n", err)
		return
	}

	if len(games) == 0 {
		fmt.Fprintf(c.output, "No games found on page %d\n", page)
		return
	}

	fmt.Fprintf(c.output, "=== Games (Page %d, Count %d) ===\n", page, count)
	fmt.Fprintln(c.output, "ID                       | Name                 | Status    | Player 1      | Player 2      | Moves")
	fmt.Fprintln(c.output, "--------------------------|----------------------|-----------|---------------|---------------|------")

	for _, g := range games {
		// Format status string
		status := "Setup"
		switch g.Status {
		case game.StatusPlaying:
			status = "Playing"
		case game.StatusWon:
			status = "Won"
		case game.StatusLost:
			status = "Lost"
		}

		// Get player names or placeholders
		player1 := "-"
		if g.Player1 != nil {
			player1 = g.Player1.Name
		}

		player2 := "-"
		if g.Player2 != nil {
			player2 = g.Player2.Name
		}

		// Format the name field
		name := "-"
		if g.Name != "" {
			name = g.Name
		}

		// Format the output in columns
		fmt.Fprintf(c.output, "%-24s | %-20s | %-9s | %-13s | %-13s | %d\n",
			g.ID, name, status, player1, player2, len(g.History))
	}

	// Add pagination help
	if len(games) == count {
		fmt.Fprintf(c.output, "\nFor next page: show-games %d %d\n", page+1, count)
	}
	if page > 0 {
		fmt.Fprintf(c.output, "For previous page: show-games %d %d\n", page-1, count)
	}
}

func (c *CLI) showGame(gameIDOrName string) {
	g, err := c.getGameByIDOrName(gameIDOrName)
	if err != nil {
		fmt.Fprintf(c.output, "Error getting game: %v\n", err)
		return
	}

	fmt.Fprintf(c.output, "Game ID: %s\n", g.ID)
	if g.Name != "" {
		fmt.Fprintf(c.output, "Game Name: %s\n", g.Name)
	}

	fmt.Fprintf(c.output, "Status: ")
	switch g.Status {
	case game.StatusSetup:
		fmt.Fprintln(c.output, "Setup Phase")
	case game.StatusPlaying:
		fmt.Fprintln(c.output, "Playing")
	case game.StatusWon:
		fmt.Fprintln(c.output, "Game Won")
	case game.StatusLost:
		fmt.Fprintln(c.output, "Game Lost")
	}

	fmt.Fprintf(c.output, "Player 1: %s\n", g.Player1.Name)
	fmt.Fprintf(c.output, "Player 2: %s\n", g.Player2.Name)
	fmt.Fprintf(c.output, "Player to move: %s\n\n", g.PlayerToMove)

	// Show board for each player
	for playerName, board := range g.Boards {
		fmt.Fprintf(c.output, "=== %s's Board ===\n", playerName)

		// Ships map
		fmt.Fprintln(c.output, "Own ships:")
		fmt.Fprintln(c.output, "  0 1 2 3 4 5 6 7 8 9")
		for y := 0; y < 10; y++ {
			fmt.Fprintf(c.output, "%d ", y)
			for x := 0; x < 10; x++ {
				// Now that we've fixed the board.Map FieldState function to swap coordinates,
				// we can call it with the correct (x,y) parameters
				fmt.Fprintf(c.output, "%c ", board.Maps[0].FieldState(x, y))
			}
			fmt.Fprintln(c.output)
		}

		// Shots map
		fmt.Fprintln(c.output, "\nShots fired at opponent:")
		fmt.Fprintln(c.output, "  0 1 2 3 4 5 6 7 8 9")
		for y := 0; y < 10; y++ {
			fmt.Fprintf(c.output, "%d ", y)
			for x := 0; x < 10; x++ {
				// Use the same coordinate system as for the ships
				fmt.Fprintf(c.output, "%c ", board.Maps[1].FieldState(x, y))
			}
			fmt.Fprintln(c.output)
		}
		fmt.Fprintln(c.output)
	}

	// Show game history
	fmt.Fprintln(c.output, "=== Game History ===")
	if len(g.History) == 0 {
		fmt.Fprintln(c.output, "No moves yet")
	} else {
		for i, move := range g.History {
			hitStatus := "miss"
			if move.Hit {
				hitStatus = "hit"
			}
			fmt.Fprintf(c.output, "%d. %s fired at (%d,%d) - %s\n", i+1, move.Player, move.X, move.Y, hitStatus)
		}
	}
}

