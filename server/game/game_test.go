package game_test

import (
	"fmt"
	"testing"

	"github.com/heartles/uttt/server/game"
)

func verboseFail(t *testing.T, format string, args ...interface{}) {
	if testing.Verbose() {
		t.Errorf(format, args...)
	} else {
		t.Fail()
	}
}

func testMove(g game.Game, m game.Move, expected error) func(*testing.T) {
	return func(t *testing.T) {
		err := g.PlayMove(m)
		if err != expected {
			t.Errorf("PlayMove returned %#v, expected %#v", err, expected)
		}

		if expected == nil && err == nil {
			testSquare(g, m.Coordinate, m.PlayerID)(t)
		}
	}
}

func testSquare(g game.Game, c game.Coordinate, expected string) func(*testing.T) {
	return func(t *testing.T) {
		playerID, err := g.SquareOwner(c)
		if err != nil {
			t.Errorf("IsSquareTaken returned error %#v", err)
		}
		if playerID != expected {
			t.Errorf("IsSquareTaken returned %#v, expected %#v", err, expected)
		}
	}
}

func Test_Game(t *testing.T) {
	g, err := game.NewGame("X", "O")
	if err != nil {
		t.Fatal(err)
	}

	// all squares should be empty at the start of the game
	for x := 1; x <= 3; x++ {
		for y := 1; y <= 3; y++ {
			for z := 1; z <= 3; z++ {
				for w := 1; w <= 3; w++ {
					t.Run(
						fmt.Sprintf("NewGameSquareEmpty%v%v%v%v", x, y, z, w),
						testSquare(g, game.NewCoordinate(1, 1, 1, 1), ""),
					)
				}
			}
		}
	}

	// PlayerID must match one of the two that the game was created with
	t.Run("InvalidPlayerFails", testMove(g, game.Move{"A", game.NewCoordinate(1, 1, 1, 1)}, game.ErrInvalidPlayer))

	// Failed moves should not change the game board
	t.Run("SquareEmptyAfterFailedMove", testSquare(g, game.NewCoordinate(1, 1, 1, 1), ""))

	// coordinates must be in interval [1, 3]
	t.Run("InvalidCoordFails", testMove(g, game.Move{"X", game.NewCoordinate(0, 1, 1, 1)}, game.ErrInvalidCoordinate))
	t.Run("InvalidCoordFails", testMove(g, game.Move{"X", game.NewCoordinate(1, 1, 1, 4)}, game.ErrInvalidCoordinate))

	t.Run("ValidMovesPass", func(t *testing.T) {
		t.Run("X1111", testMove(g, game.Move{"X", game.NewCoordinate(1, 1, 1, 1)}, nil))
		t.Run("O1121", testMove(g, game.Move{"O", game.NewCoordinate(1, 1, 2, 1)}, nil))
		t.Run("X2111", testMove(g, game.Move{"X", game.NewCoordinate(2, 1, 1, 1)}, nil))
		t.Run("O1113", testMove(g, game.Move{"O", game.NewCoordinate(1, 1, 1, 3)}, nil))
	})

	// A move can only be played if it is for the player whose turn it is
	t.Run("WrongPlayerFails", testMove(g, game.Move{"O", game.NewCoordinate(1, 1, 3, 1)}, game.ErrWrongTurn))

	// Players can only play in the subgrid corresponding to the last move
	// EXCEPT in special circumstances, see below tests
	t.Run("WrongSubgridFails", testMove(g, game.Move{"X", game.NewCoordinate(3, 3, 3, 1)}, game.ErrWrongSubgrid))

	// TODO: Finish writing tests
	lastTurn := game.NewCoordinate(2, 1, 1, 1)
	g, err = game.LoadGame("X", "O",
		`XX_    O__    ___
		 ___    ___    ___
		 ___    ___    ___

		 ___    ___    ___
		 ___    ___    ___
		 ___    ___    ___

		 ___    ___    ___
		 ___    ___    ___
		 ___    ___    ___`,
		&lastTurn)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("ConnectThreeWinsBlock", func(t *testing.T) {
		testMove(g, game.Move{"X", game.NewCoordinate(1, 1, 3, 1)}, nil)(t)
		winner, err := g.BlockWinner(game.SubCoordinate{1, 1})
		if err != nil {
			t.Errorf("Error: %v", err)
		} else if winner != "X" {
			t.Errorf("Subgrid outcome was %v, expected %v", winner, "X")
		}
	})

	testResetMove := func(m game.Move, expected error) func(*testing.T) {
		return func(t *testing.T) {
			// next turn is X's
			lastMove := game.NewCoordinate(2, 1, 1, 1)
			g, err := game.LoadGame("X", "O",
				`XXX    O__    ___
			 ___    _O_    XX_
			 ___    __O    ___

			 ___    ___    ___
			 ___    ___    OOO
			 ___    ___    ___

			 ___    ___    ___
			 ___    ___    _O_
			 ___    ___    ___
		`, &lastMove)
			if err != nil {
				t.Fatalf("%v", err)
			}

			err = g.PlayMove(m)
			if err != expected {
				t.Errorf("PlayMove returned %v, expected %v", err, expected)
			}
		}
	}

	// TODO: should probably properly name these
	t.Run("TestMoveIntoFullSubgrid", testResetMove(game.Move{"X", game.NewCoordinate(1, 1, 2, 3)}, game.ErrWrongSubgrid))
	t.Run("TestMoveIntoFullSubgrid", testResetMove(game.Move{"X", game.NewCoordinate(1, 1, 1, 1)}, game.ErrWrongSubgrid))
	t.Run("TestMoveIntoFullSubgrid", testResetMove(game.Move{"X", game.NewCoordinate(2, 1, 2, 2)}, game.ErrWrongSubgrid))
	t.Run("TestMoveIntoFullSubgrid", testResetMove(game.Move{"X", game.NewCoordinate(3, 1, 2, 3)}, nil))
	t.Run("TestMoveIntoFullSubgrid", testResetMove(game.Move{"X", game.NewCoordinate(3, 1, 2, 2)}, game.ErrSquarePlayed))
	t.Run("TestMoveIntoFullSubgrid", testResetMove(game.Move{"X", game.NewCoordinate(3, 2, 2, 2)}, game.ErrWrongSubgrid))
	t.Run("TestMoveIntoFullSubgrid", testResetMove(game.Move{"X", game.NewCoordinate(3, 3, 2, 2)}, game.ErrSquarePlayed))
	t.Run("TestMoveIntoFullSubgrid", testResetMove(game.Move{"X", game.NewCoordinate(2, 3, 2, 2)}, nil))
	t.Run("TestMoveIntoFullSubgrid", testResetMove(game.Move{"X", game.NewCoordinate(2, 2, 2, 2)}, nil))
	t.Run("TestMoveIntoFullSubgrid", testResetMove(game.Move{"X", game.NewCoordinate(1, 3, 2, 2)}, nil))
	t.Run("TestMoveIntoFullSubgrid", testResetMove(game.Move{"X", game.NewCoordinate(1, 2, 2, 2)}, nil))
	t.Run("TestMoveIntoFullSubgrid", testResetMove(game.Move{"O", game.NewCoordinate(1, 2, 2, 2)}, game.ErrWrongTurn))

	testFirstMove := func(m game.Move, expected error) {

		g, err := game.LoadGame("X", "O",
			`___    ___    ___
			 ___    ___    ___
			 ___    ___    ___

			 ___    ___    ___
			 ___    ___    ___
			 ___    ___    ___

			 ___    ___    ___
			 ___    ___    ___
			 ___    ___    ___
		`, nil)
		if err != nil {
			t.Fatalf("%v", err)
		}

		err = g.PlayMove(m)
		if err != expected {
			t.Errorf("PlayMove returned %v, expected %v", err, expected)
		}
	}
	testFirstMove(game.Move{"X", game.NewCoordinate(1, 1, 2, 3)}, nil)
	testFirstMove(game.Move{"X", game.NewCoordinate(1, 1, 1, 1)}, nil)
	testFirstMove(game.Move{"X", game.NewCoordinate(2, 1, 2, 2)}, nil)
	testFirstMove(game.Move{"X", game.NewCoordinate(3, 1, 2, 3)}, nil)
	testFirstMove(game.Move{"X", game.NewCoordinate(3, 1, 2, 2)}, nil)
	testFirstMove(game.Move{"X", game.NewCoordinate(3, 2, 2, 2)}, nil)
	testFirstMove(game.Move{"X", game.NewCoordinate(3, 3, 2, 2)}, nil)
	testFirstMove(game.Move{"X", game.NewCoordinate(2, 3, 2, 2)}, nil)
	testFirstMove(game.Move{"X", game.NewCoordinate(2, 2, 2, 2)}, nil)
	testFirstMove(game.Move{"X", game.NewCoordinate(1, 3, 2, 2)}, nil)
	testFirstMove(game.Move{"X", game.NewCoordinate(1, 2, 2, 2)}, nil)
}

func TestCalculateVictor(t *testing.T) {
	var g game.Game

	testBlock := func(x, y int, expected string) {
		winner, err := g.BlockWinner(game.SubCoordinate{x, y})
		if err != nil {
			t.Errorf("%v", err)
			return
		}

		if winner != expected {
			t.Errorf("block {%v, %v} had winner %#v, expected %#v", x, y, winner, expected)
		}
	}

	lastTurn := game.NewCoordinate(2, 1, 1, 1)
	g, err := game.LoadGame("X", "O",
		`XX_    OX_    X_O
		 X__    _X_    _O_
		 ___    _X_    O_X

		 _O_    ___    X__
		 OOO    _X_    _O_
		 _OX    ___    __X

		 ___    ___    XOX
		 XXX    ___    OOX
		 ___    OOO    XXO`,
		&lastTurn)
	if err != nil {
		t.Fatal(err)
	}

	testBlock(1, 1, "")
	testBlock(2, 1, "X")
	testBlock(3, 1, "O")
	testBlock(1, 2, "O")
	testBlock(2, 2, "")
	testBlock(3, 2, "")
	testBlock(1, 3, "X")
	testBlock(2, 3, "O")
	testBlock(3, 3, game.StalematePlayer)
	if g.GameWinner() != "" {
		t.Errorf("game in progress but winner declared: %v", g.GameWinner())
	}

	lastTurn = game.NewCoordinate(2, 1, 1, 1)
	g, err = game.LoadGame("X", "O",
		`XXX    OX_    X_O
		 X__    _X_    _X_
		 ___    _X_    O_X

		 _O_    ___    X__
		 OOO    _X_    _O_
		 _OX    ___    __X

		 ___    ___    XOX
		 XXX    ___    OOX
		 ___    OOO    XXO`,
		&lastTurn)
	if err != nil {
		t.Fatal(err)
	}

	testBlock(1, 1, "X")
	testBlock(2, 1, "X")
	testBlock(3, 1, "X")
	testBlock(1, 2, "O")
	testBlock(2, 2, "")
	testBlock(3, 2, "")
	testBlock(1, 3, "X")
	testBlock(2, 3, "O")
	testBlock(3, 3, game.StalematePlayer)
	if g.GameWinner() != "X" {
		t.Errorf("incorrect game winner: %#v expected %v", g.GameWinner(), "X")
	}

	lastTurn = game.NewCoordinate(2, 1, 1, 1)
	g, err = game.LoadGame("X", "O",
		`OXO    OX_    X_O
		 XXO   _X_    _X_
		 __O    _X_    O_X

		 XO_    __X    X_O
		 OXO    _XX    _O_
		 _OX    __X    O_X

		 OOO    ___    XOO
		 XXO    ___    OXX
		 ___    OOO    XXO`,
		&lastTurn)
	if err != nil {
		t.Fatal(err)
	}

	testBlock(1, 1, "O")
	testBlock(2, 1, "X")
	testBlock(3, 1, "X")
	testBlock(1, 2, "X")
	testBlock(2, 2, "X")
	testBlock(3, 2, "O")
	testBlock(1, 3, "O")
	testBlock(2, 3, "O")
	testBlock(3, 3, game.StalematePlayer)
	if g.GameWinner() != game.StalematePlayer {
		t.Errorf("incorrect game winner: %#v expected %v", g.GameWinner(), game.StalematePlayer)
	}

	lastTurn = game.NewCoordinate(2, 1, 1, 1)
	g, err = game.LoadGame("X", "O",
		`OOO    OX_    X_O
		 X__    _X_    _X_
		 ___    _X_    O_X

		 _O_    ___    X__
		 OOO    _X_    _O_
		 _OX    ___    __X

		 OOO    ___    XOX
		 XXO    ___    OOX
		 ___    OOO    XXO`,
		&lastTurn)
	if err != nil {
		t.Fatal(err)
	}

	testBlock(1, 1, "O")
	testBlock(2, 1, "X")
	testBlock(3, 1, "X")
	testBlock(1, 2, "O")
	testBlock(2, 2, "")
	testBlock(3, 2, "")
	testBlock(1, 3, "O")
	testBlock(2, 3, "O")
	testBlock(3, 3, game.StalematePlayer)
	if g.GameWinner() != "O" {
		t.Errorf("incorrect game winner: %#v expected %v", g.GameWinner(), "O")
	}
}
