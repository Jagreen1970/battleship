func TestNewBoard(t *testing.T) {
    // Test data
    playerName := "player1"
    opponentName := "player2"

    // Create new board
    board := NewBoard(playerName, opponentName)

    // Assert board is not nil
    assert.NotNil(t, board)

    // Check initial pins available
    assert.Equal(t, 30, board.PinsAvailable)

    // Check maps initialization
    assert.Len(t, board.Maps, 2)
    assert.NotNil(t, board.Maps[0])
    assert.NotNil(t, board.Maps[1])

    // Check map titles
    assert.Equal(t, playerName, board.Maps[0].Title)
    assert.Equal(t, opponentName, board.Maps[1].Title)

    // Check that fleet is initially nil
    assert.Nil(t, board.Fleet)

    // Check that all fields are initialized to FieldStateEmpty
    for x := 0; x < 10; x++ {
        for y := 0; y < 10; y++ {
            assert.Equal(t, FieldStateEmpty, board.Maps[0].FieldState(x, y), 
                "Player map field at (%d,%d) should be empty", x, y)
            assert.Equal(t, FieldStateEmpty, board.Maps[1].FieldState(x, y), 
                "Opponent map field at (%d,%d) should be empty", x, y)
        }
    }
}