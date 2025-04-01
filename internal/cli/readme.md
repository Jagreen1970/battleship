# Battleship CLI Mode

The CLI mode provides a command-line interface for testing the Battleship game.

## Running the CLI

```bash
# Start CLI mode with default settings (will likely fail without authentication)
./battleship --cli

# Start with MongoDB authentication
./battleship --cli --dbuser root --dbpass battleship

# Or set environment variables
export DB_USER=root
export DB_PASSWORD=battleship
./battleship --cli
```

## Available Commands

The following commands are available within the CLI:

- `create-game <player> [name]`: Create a new game session with optional friendly name
- `show-games [page] [count]`: List all active games (paginated)
- `show-game <game-id|name>`: Show game status and boards of a specific game
- `join-game <game-id|name> <player>`: Join an existing game as a player
- `set-game <game-id|name>`: Set the game ID for the current session (all future actions will be performed on this game)
- `place-ship <player> <ship-type> <x> <y> <orientation>`: Place a ship. The game will automatically start when all ships are placed.
- `start-game <player>`: Start the game with the given player
- `fire <player> <x> <y>`: Fire at coordinates
- `delete-game <game-id|name|all>`: Delete a specific game or all games
- `exit`: Exit CLI mode

## Game Elements

- **Ship types**: Battleship (5 tiles), Cruiser (4 tiles), Destroyer (3 tiles), Submarine (2 tiles)
- **Orientation**: Horizontal, Vertical

## MongoDB Authentication

The CLI requires MongoDB authentication. You must provide valid credentials either through environment variables or command-line flags:

- Through environment variables: Set `DB_USER` and `DB_PASSWORD`
- Through command-line flags: Use `--dbuser` and `--dbpass`

### Default MongoDB Authentication

If you're running the application with Docker Compose as described in the README, the following credentials are configured:

```
Username: root
Password: battleship
```

### Authentication Error

If you see the following error:

```
Error creating player: error fetching player: (Unauthorized) Command find requires authentication
```

It means you need to provide MongoDB credentials as described above.

## Workflow Example

```
> create-game player1 MyFirstGame
Created new game 'MyFirstGame' with ID: 6440a28f9d5c7b9a0f1a2b3c

> show-games
=== Games (Page 0, Count 10) ===
ID                       | Name                 | Status    | Player 1      | Player 2      | Moves
--------------------------|----------------------|-----------|---------------|---------------|------
6440a28f9d5c7b9a0f1a2b3c | MyFirstGame          | Setup     | player1       | -             | 0

> join-game MyFirstGame player2
Player player2 joined game 'MyFirstGame'

> set-game MyFirstGame
Current game set to: MyFirstGame (ID: 6440a28f9d5c7b9a0f1a2b3c)

> place-ship player1 Battleship 0 0 Horizontal
Placed Battleship at (0,0) Horizontal for player player1

> place-ship player2 Battleship 0 0 Horizontal
Placed Battleship at (0,0) Horizontal for player player2

> show-game
Game ID: 6440a28f9d5c7b9a0f1a2b3c
Game Name: MyFirstGame
Status: Setup Phase
Player 1: player1
Player 2: player2
Player to move: 

> fire player1 0 0
Player player1 fired at (0,0) - hit
Next player to move: player2

> delete-game MyFirstGame
Game 'MyFirstGame' (ID: 6440a28f9d5c7b9a0f1a2b3c) successfully deleted
Current game unset
```

## Using Named Games

With the optional name parameter for `create-game`, you can now refer to games by their friendly name instead of their ID.

```
> create-game player1 Practice
Created new game 'Practice' with ID: 6440a28f9d5c7b9a0f1a2b3c

> set-game Practice
Current game set to: Practice (ID: 6440a28f9d5c7b9a0f1a2b3c)

> join-game Practice player2
Player player2 joined game 'Practice'
```

## Command Structure and Validation

The CLI includes validation for all commands to ensure they have the required parameters:

- Commands always show usage instructions when called incorrectly
- Proper validation prevents errors and crashes
- The `place-ship` command requires all 5 parameters (player, ship-type, x, y, orientation)
- Orientation must be either "Horizontal" or "Vertical" (case-sensitive)

