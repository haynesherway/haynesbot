package pokemon

import (
    "strings"
)

type Moves struct {
    Fast MoveList `json:"quickMoves"`
    Charge MoveList `json:"cinematicMoves"`
}

type MoveList []*PokemonMove

type PokemonMove struct {
    ID string `json:"id"`
    Name string `json:"name"`
}

func (moveList MoveList) Print() string {
    moves := []string{}
    for _, m := range moveList {
        //Remove "Fast"
        move := strings.TrimSpace(strings.Replace(m.Name, "Fast", "", 1))
        moves = append(moves, move)
    }
    return strings.Join(moves, ", ")
}

