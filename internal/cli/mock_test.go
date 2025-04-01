package cli

import (
	"fmt"
	"testing"

	"github.com/Jagreen1970/battleship/internal/game"
	"github.com/stretchr/testify/require"
)

// mockStorage implements the storage.Storage interface for testing
type mockStorage struct {
	t              *testing.T
	players        map[string]*game.Player
	games          map[string]*game.Game
	mockCreateGame func(*game.Game) (*game.Game, error)
}

// newMockStorage creates a new mock storage for testing
func newMockStorage(t *testing.T) *mockStorage {
	return &mockStorage{
		t:       t,
		players: make(map[string]*game.Player),
		games:   make(map[string]*game.Game),
		mockCreateGame: func(g *game.Game) (*game.Game, error) {
			g.ID = "mock-game-id"
			return g, nil
		},
	}
}

// Connect implements the storage.Storage interface
func (m *mockStorage) Connect() error {
	return nil
}

// Disconnect implements the storage.Storage interface
func (m *mockStorage) Disconnect() error {
	return nil
}

// Ping implements the storage.Storage interface
func (m *mockStorage) Ping() error {
	return nil
}

// Close implements the storage.Storage interface
func (m *mockStorage) Close() error {
	return nil
}

// CreatePlayer implements the storage.Storage interface
func (m *mockStorage) CreatePlayer(playerName string) (*game.Player, error) {
	player := &game.Player{
		Name: playerName,
		ID:   "player-" + playerName,
	}
	m.players[playerName] = player
	return player, nil
}

// FindPlayerByName implements the storage.Storage interface
func (m *mockStorage) FindPlayerByName(username string) (*game.Player, error) {
	player, ok := m.players[username]
	if !ok {
		return nil, fmt.Errorf("player not found: %w", game.ErrorNotFound)
	}
	return player, nil
}

// QueryGames implements the storage.Storage interface
func (m *mockStorage) QueryGames(page int, count int) ([]*game.Game, error) {
	var games []*game.Game
	for _, g := range m.games {
		games = append(games, g)
	}
	
	// Apply pagination
	start := page * count
	end := start + count
	
	// Handle edge cases
	if start >= len(games) {
		return []*game.Game{}, nil
	}
	if end > len(games) {
		end = len(games)
	}
	
	return games[start:end], nil
}

// CreateGame implements the storage.Storage interface
func (m *mockStorage) CreateGame(g *game.Game) (*game.Game, error) {
	if m.mockCreateGame != nil {
		game, err := m.mockCreateGame(g)
		if err == nil && game != nil && game.ID != "" {
			m.games[game.ID] = game  // Make sure to store the game in the map
		}
		return game, err
	}
	g.ID = "mock-game-id"
	m.games[g.ID] = g
	return g, nil
}

// FindGameByID implements the storage.Storage interface
func (m *mockStorage) FindGameByID(id string) (*game.Game, error) {
	g, ok := m.games[id]
	if !ok {
		return nil, fmt.Errorf("game not found: %w", game.ErrorNotFound)
	}
	return g, nil
}

// FindGameByName implements the storage.Storage interface
func (m *mockStorage) FindGameByName(name string) (*game.Game, error) {
	if name == "" {
		return nil, fmt.Errorf("game name cannot be empty: %w", game.ErrorInvalidInput)
	}
	
	for _, g := range m.games {
		if g.Name == name {
			return g, nil
		}
	}
	
	return nil, fmt.Errorf("game with name '%s' not found: %w", name, game.ErrorNotFound)
}

// UpdateGame implements the storage.Storage interface
func (m *mockStorage) UpdateGame(g *game.Game) (*game.Game, error) {
	require.NotEmpty(m.t, g.ID, "Game ID should not be empty in UpdateGame")
	m.games[g.ID] = g
	return g, nil
}

// DeleteGame implements the storage.Storage interface
func (m *mockStorage) DeleteGame(id string) error {
	_, exists := m.games[id]
	if !exists {
		return game.ErrorNotFound
	}
	
	delete(m.games, id)
	return nil
}

// DeleteAllGames implements the storage.Storage interface
func (m *mockStorage) DeleteAllGames() (int, error) {
	count := len(m.games)
	m.games = make(map[string]*game.Game)
	return count, nil
}