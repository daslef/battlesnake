package rulesets

import (
	"rules"
	"rules/settings"
)

var soloRulesetStages = []string{
	StageGameOverSoloSnake,
	StageMovementStandard,
	StageStarvationStandard,
	StageFeedSnakesStandard,
	StageEliminationStandard,
}

func GameOverSolo(b *rules.BoardState, settings settings.Settings, moves []SnakeMove) (bool, error) {
	for i := 0; i < len(b.Snakes); i++ {
		if b.Snakes[i].EliminatedCause == rules.NotEliminated {
			return false, nil
		}
	}
	return true, nil
}
