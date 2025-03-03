package maps

import (
	"rules"
	"rules/settings"
)

type StandardMap struct{}

func (m StandardMap) ID() string {
	return "standard"
}

func (m StandardMap) Meta() Metadata {
	return Metadata{
		Name:       "Standard",
		MinPlayers: 1,
		MaxPlayers: 16,
		BoardSizes: OddSizes(rules.BoardSizeSmall, rules.BoardSizeXXLarge),
	}
}

func (m StandardMap) SetupBoard(initialBoardState *rules.BoardState, settings settings.Settings, editor Editor) error {
	rand := settings.GetRand(0)

	if len(initialBoardState.Snakes) > int(m.Meta().MaxPlayers) {
		return rules.ErrorTooManySnakes
	}

	snakeIDs := make([]string, 0, len(initialBoardState.Snakes))
	for _, snake := range initialBoardState.Snakes {
		snakeIDs = append(snakeIDs, snake.ID)
	}

	tempBoardState, err := rules.CreateDefaultBoardState(rand, initialBoardState.Width, initialBoardState.Height, snakeIDs)
	if err != nil {
		return err
	}

	// Copy food from temp board state
	for _, food := range tempBoardState.Food {
		editor.AddFood(food)
	}

	// Copy snakes from temp board state
	for _, snake := range tempBoardState.Snakes {
		editor.PlaceSnake(snake.ID, snake.Body, snake.Health)
	}

	return nil
}

func (m StandardMap) PreUpdateBoard(lastBoardState *rules.BoardState, settings settings.Settings, editor Editor) error {
	return nil
}

func (m StandardMap) PostUpdateBoard(lastBoardState *rules.BoardState, settings settings.Settings, editor Editor) error {
	rand := settings.GetRand(lastBoardState.Turn)

	foodNeeded := checkFoodNeedingPlacement(rand, settings)
	if foodNeeded > 0 {
		placeFoodRandomly(rand, lastBoardState, editor, foodNeeded)
	}

	return nil
}

func checkFoodNeedingPlacement(rand rules.Rand, settings settings.Settings) int {
	foodSpawnChance := settings.Int(rules.ParamFoodSpawnChance, 0)

	if foodSpawnChance > 0 && (100-rand.Intn(100)) < foodSpawnChance {
		return 1
	}

	return 0
}

func placeFoodRandomly(rand rules.Rand, b *rules.BoardState, editor Editor, n int) {
	unoccupiedPoints := rules.GetUnoccupiedPoints(b, false)
	placeFoodRandomlyAtPositions(rand, editor, n, unoccupiedPoints)
}

func placeFoodRandomlyAtPositions(rand rules.Rand, editor Editor, n int, positions []rules.Point) {
	if len(positions) < n {
		n = len(positions)
	}

	rand.Shuffle(len(positions), func(i int, j int) {
		positions[i], positions[j] = positions[j], positions[i]
	})

	for i := 0; i < n; i++ {
		editor.AddFood(positions[i])
	}
}
