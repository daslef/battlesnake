package rulesets

import (
	"fmt"
	"rules"
	"rules/settings"
)

const (
	StageSpawnFoodStandard   = "spawn_food.standard"
	StageGameOverStandard    = "game_over.standard"
	StageStarvationStandard  = "starvation.standard"
	StageFeedSnakesStandard  = "feed_snakes.standard"
	StageMovementStandard    = "movement.standard"
	StageEliminationStandard = "elimination.standard"

	StageGameOverSoloSnake = "game_over.solo_snake"
)

// globalRegistry is a global, default mapping of stage names to stage functions.
// It can be extended by plugins through the use of registration functions.
// Plugins that wish to extend the available game stages should call RegisterPipelineStageError
// to add additional stages.
var globalRegistry = StageRegistry{
	StageGameOverSoloSnake:   GameOverSolo,
	StageGameOverStandard:    GameOverStandard,
	StageStarvationStandard:  ReduceSnakeHealthStandard,
	StageFeedSnakesStandard:  FeedSnakesStandard,
	StageEliminationStandard: EliminateSnakesStandard,
	StageMovementStandard:    MoveSnakesStandard,
}

// Pipeline is an ordered sequences of game stages which are executed to produce the
// next game state.
//
// If a stage produces an error or an ended game state, the pipeline is halted at that stage.
type Pipeline interface {
	Execute(*rules.BoardState, settings.Settings, []SnakeMove) (bool, *rules.BoardState, error)
	Err() error
}

// StageFunc represents a single stage of an ordered pipeline and applies custom logic to the board state each turn.
// It is expected to modify the boardState directly.
// The return values are a boolean (to indicate whether the game has ended as a result of the stage)
// and an error if any errors occurred during the stage.
//
// Errors should be treated as meaning the stage failed and the board state is now invalid.
type StageFunc func(*rules.BoardState, settings.Settings, []SnakeMove) (bool, error)

// IsInitialization checks whether the current state means the game is initialising (turn zero).
func IsInitialization(b *rules.BoardState, settings settings.Settings, moves []SnakeMove) bool {
	// We can safely assume that the game state is in the initialisation phase when
	// the turn hasn't advanced and the moves are empty
	return b.Turn <= 0 && len(moves) == 0
}

// StageRegistry is a mapping of stage names to stage functions
type StageRegistry map[string]StageFunc

// RegisterPipelineStage adds a stage to the registry.
// If a stage has already been mapped it will be overwritten by the newly
// registered function.
func (sr StageRegistry) RegisterPipelineStage(s string, fn StageFunc) {
	sr[s] = fn
}

// RegisterPipelineStageError adds a stage to the registry.
// If a stage has already been mapped an error will be returned.
func (sr StageRegistry) RegisterPipelineStageError(s string, fn StageFunc) error {
	if _, ok := sr[s]; ok {
		return rules.RulesetError(fmt.Sprintf("stage '%s' has already been registered", s))
	}

	sr.RegisterPipelineStage(s, fn)
	return nil
}

// RegisterPipelineStage adds a stage to the global stage registry.
// It will panic if the a stage has already been registered with the same name.
func RegisterPipelineStage(s string, fn StageFunc) {
	err := globalRegistry.RegisterPipelineStageError(s, fn)
	if err != nil {
		panic(err)
	}
}

type pipeline struct {
	stages []StageFunc
	err    error
}

// NewPipeline constructs an instance of Pipeline using the global registry.
func NewPipeline(stageNames ...string) Pipeline {
	if len(globalRegistry) == 0 {
		return &pipeline{err: rules.ErrorEmptyRegistry}
	}

	if len(stageNames) == 0 {
		return &pipeline{err: rules.ErrorNoStages}
	}

	p := pipeline{}
	for _, s := range stageNames {
		fn, ok := globalRegistry[s]
		if !ok {
			return pipeline{err: rules.ErrorStageNotFound}
		}

		p.stages = append(p.stages, fn)
	}

	return &p
}

// impl
func (p pipeline) Err() error {
	return p.err
}

// impl
func (p pipeline) Execute(state *rules.BoardState, settings settings.Settings, moves []SnakeMove) (bool, *rules.BoardState, error) {
	// If the pipeline is in an error state, Execute must return that error
	// because the pipeline is invalid and cannot execute.
	//
	// This is done for API use convenience to satisfy the common pattern
	// of wanting to write NewPipeline().Execute(...).

	if p.err != nil {
		return false, nil, p.err
	}

	var ended bool
	var err error
	state = state.Clone()
	for _, fn := range p.stages {
		// execute current stage
		ended, err = fn(state, settings, moves)

		// stop if we hit any errors or if the game is ended
		if err != nil || ended {
			return ended, state, err
		}
	}

	return ended, state, err
}
