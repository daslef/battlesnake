package rulesets

import (
	"rules"
	"rules/settings"
	"sort"
)

var standardRulesetStages = []string{
	StageGameOverStandard,
	StageMovementStandard,
	StageStarvationStandard,
	StageFeedSnakesStandard,
	StageEliminationStandard,
}

func MoveSnakesStandard(b *rules.BoardState, settings settings.Settings, moves []SnakeMove) (bool, error) {
	if IsInitialization(b, settings, moves) {
		return false, nil
	}

	// no-op when moves are empty
	if len(moves) == 0 {
		return false, nil
	}

	// Sanity check that all non-eliminated snakes have moves and bodies.
	for i := 0; i < len(b.Snakes); i++ {
		snake := &b.Snakes[i]
		if snake.EliminatedCause != rules.NotEliminated {
			continue
		}

		if len(snake.Body) == 0 {
			return false, rules.ErrorZeroLengthSnake
		}
		moveFound := false
		for _, move := range moves {
			if snake.ID == move.ID {
				moveFound = true
				break
			}
		}
		if !moveFound {
			return false, rules.ErrorNoMoveFound
		}
	}

	for i := 0; i < len(b.Snakes); i++ {
		snake := &b.Snakes[i]
		if snake.EliminatedCause != rules.NotEliminated {
			continue
		}

		for _, move := range moves {
			if move.ID == snake.ID {
				appliedMove := move.Move
				switch move.Move {
				case rules.MoveUp, rules.MoveDown, rules.MoveRight, rules.MoveLeft:
					break
				default:
					appliedMove = getDefaultMove(snake.Body)
				}

				newHead := rules.Point{}
				switch appliedMove {
				// Guaranteed to be one of these options given the clause above
				case rules.MoveUp:
					newHead.X = snake.Body[0].X
					newHead.Y = snake.Body[0].Y + 1
				case rules.MoveDown:
					newHead.X = snake.Body[0].X
					newHead.Y = snake.Body[0].Y - 1
				case rules.MoveLeft:
					newHead.X = snake.Body[0].X - 1
					newHead.Y = snake.Body[0].Y
				case rules.MoveRight:
					newHead.X = snake.Body[0].X + 1
					newHead.Y = snake.Body[0].Y
				}

				// Append new head, pop old tail
				snake.Body = append([]rules.Point{newHead}, snake.Body[:len(snake.Body)-1]...)
			}
		}
	}
	return false, nil
}

func getDefaultMove(snakeBody []rules.Point) string {
	if len(snakeBody) >= 2 {
		// Use neck to determine last move made
		head, neck := snakeBody[0], snakeBody[1]
		// Situations where neck is next to head
		if head.X == neck.X+1 {
			return rules.MoveRight
		} else if head.X == neck.X-1 {
			return rules.MoveLeft
		} else if head.Y == neck.Y+1 {
			return rules.MoveUp
		} else if head.Y == neck.Y-1 {
			return rules.MoveDown
		}
		// Consider the wrapped cases using zero axis to anchor
		if head.X == 0 && neck.X > 0 {
			return rules.MoveRight
		} else if neck.X == 0 && head.X > 0 {
			return rules.MoveLeft
		} else if head.Y == 0 && neck.Y > 0 {
			return rules.MoveUp
		} else if neck.Y == 0 && head.Y > 0 {
			return rules.MoveDown
		}
	}
	return rules.MoveUp
}

func ReduceSnakeHealthStandard(b *rules.BoardState, settings settings.Settings, moves []SnakeMove) (bool, error) {
	if IsInitialization(b, settings, moves) {
		return false, nil
	}
	for i := 0; i < len(b.Snakes); i++ {
		if b.Snakes[i].EliminatedCause == rules.NotEliminated {
			b.Snakes[i].Health = b.Snakes[i].Health - 1
		}
	}
	return false, nil
}

func EliminateSnakesStandard(b *rules.BoardState, settings settings.Settings, moves []SnakeMove) (bool, error) {
	if IsInitialization(b, settings, moves) {
		return false, nil
	}
	// First order snake indices by length.
	// In multi-collision scenarios we want to always attribute elimination to the longest snake.
	snakeIndicesByLength := make([]int, len(b.Snakes))
	for i := 0; i < len(b.Snakes); i++ {
		snakeIndicesByLength[i] = i
	}
	sort.Slice(snakeIndicesByLength, func(i int, j int) bool {
		lenI := len(b.Snakes[snakeIndicesByLength[i]].Body)
		lenJ := len(b.Snakes[snakeIndicesByLength[j]].Body)
		return lenI > lenJ
	})

	// First, iterate over all non-eliminated snakes and eliminate the ones
	// that are out of health or have moved out of bounds.
	for i := 0; i < len(b.Snakes); i++ {
		snake := &b.Snakes[i]
		if snake.EliminatedCause != rules.NotEliminated {
			continue
		}
		if len(snake.Body) <= 0 {
			return false, rules.ErrorZeroLengthSnake
		}

		if snakeIsOutOfHealth(snake) {
			rules.EliminateSnake(snake, rules.EliminatedByOutOfHealth, "", b.Turn+1)
			continue
		}

		if snakeIsOutOfBounds(snake, b.Width, b.Height) {
			rules.EliminateSnake(snake, rules.EliminatedByOutOfBounds, "", b.Turn+1)
			continue
		}
	}

	// Next, look for any collisions. Note we apply collision eliminations
	// after this check so that snakes can collide with each other and be properly eliminated.
	type CollisionElimination struct {
		ID    string
		Cause string
		By    string
	}
	collisionEliminations := []CollisionElimination{}
	for i := 0; i < len(b.Snakes); i++ {
		snake := &b.Snakes[i]
		if snake.EliminatedCause != rules.NotEliminated {
			continue
		}
		if len(snake.Body) <= 0 {
			return false, rules.ErrorZeroLengthSnake
		}

		// Check for self-collisions first
		if snakeHasBodyCollided(snake, snake) {
			collisionEliminations = append(collisionEliminations, CollisionElimination{
				ID:    snake.ID,
				Cause: rules.EliminatedBySelfCollision,
				By:    snake.ID,
			})
			continue
		}

		// Check for body collisions with other snakes second
		hasBodyCollided := false
		for _, otherIndex := range snakeIndicesByLength {
			other := &b.Snakes[otherIndex]
			if other.EliminatedCause != rules.NotEliminated {
				continue
			}
			if snake.ID != other.ID && snakeHasBodyCollided(snake, other) {
				collisionEliminations = append(collisionEliminations, CollisionElimination{
					ID:    snake.ID,
					Cause: rules.EliminatedByCollision,
					By:    other.ID,
				})
				hasBodyCollided = true
				break
			}
		}
		if hasBodyCollided {
			continue
		}

		// Check for head-to-heads last
		hasHeadCollided := false
		for _, otherIndex := range snakeIndicesByLength {
			other := &b.Snakes[otherIndex]
			if other.EliminatedCause != rules.NotEliminated {
				continue
			}
			if snake.ID != other.ID && snakeHasLostHeadToHead(snake, other) {
				collisionEliminations = append(collisionEliminations, CollisionElimination{
					ID:    snake.ID,
					Cause: rules.EliminatedByHeadToHeadCollision,
					By:    other.ID,
				})
				hasHeadCollided = true
				break
			}
		}
		if hasHeadCollided {
			continue
		}
	}

	// Apply collision eliminations
	for _, elimination := range collisionEliminations {
		for i := 0; i < len(b.Snakes); i++ {
			snake := &b.Snakes[i]
			if snake.ID == elimination.ID {
				rules.EliminateSnake(snake, elimination.Cause, elimination.By, b.Turn+1)
				break
			}
		}
	}

	return false, nil
}

func snakeIsOutOfHealth(s *rules.Snake) bool {
	return s.Health <= 0
}

func snakeIsOutOfBounds(s *rules.Snake, boardWidth int, boardHeight int) bool {
	for _, point := range s.Body {
		if (point.X < 0) || (point.X >= boardWidth) {
			return true
		}
		if (point.Y < 0) || (point.Y >= boardHeight) {
			return true
		}
	}
	return false
}

func snakeHasBodyCollided(s *rules.Snake, other *rules.Snake) bool {
	head := s.Body[0]
	for i, body := range other.Body {
		if i == 0 {
			continue
		} else if head.X == body.X && head.Y == body.Y {
			return true
		}
	}
	return false
}

func snakeHasLostHeadToHead(s *rules.Snake, other *rules.Snake) bool {
	if s.Body[0].X == other.Body[0].X && s.Body[0].Y == other.Body[0].Y {
		return len(s.Body) <= len(other.Body)
	}
	return false
}

func FeedSnakesStandard(b *rules.BoardState, settings settings.Settings, moves []SnakeMove) (bool, error) {
	newFood := []rules.Point{}
	for _, food := range b.Food {
		foodHasBeenEaten := false
		for i := 0; i < len(b.Snakes); i++ {
			snake := &b.Snakes[i]

			// Ignore eliminated and zero-length snakes, they can't eat.
			if snake.EliminatedCause != rules.NotEliminated || len(snake.Body) == 0 {
				continue
			}

			if snake.Body[0].X == food.X && snake.Body[0].Y == food.Y {
				feedSnake(snake)
				foodHasBeenEaten = true
			}
		}
		// Persist food to next BoardState if not eaten
		if !foodHasBeenEaten {
			newFood = append(newFood, food)
		}
	}

	b.Food = newFood
	return false, nil
}

func feedSnake(snake *rules.Snake) {
	growSnake(snake)
	snake.Health = rules.SnakeMaxHealth
}

func growSnake(snake *rules.Snake) {
	if len(snake.Body) > 0 {
		snake.Body = append(snake.Body, snake.Body[len(snake.Body)-1])
	}
}

func GameOverStandard(b *rules.BoardState, settings settings.Settings, moves []SnakeMove) (bool, error) {
	numSnakesRemaining := 0
	for i := 0; i < len(b.Snakes); i++ {
		if b.Snakes[i].EliminatedCause == rules.NotEliminated {
			numSnakesRemaining++
		}
	}
	return numSnakesRemaining <= 1, nil
}
