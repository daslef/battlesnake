package rulesets

import (
	"errors"
	"rules"
	"rules/settings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPipeline(t *testing.T) {
	r := StageRegistry{}

	// test empty registry error
	p := NewPipelineFromRegistry(r)
	require.Equal(t, rules.ErrorEmptyRegistry, p.Err())
	_, _, err := p.Execute(nil, settings.Settings{}, nil)
	require.Equal(t, rules.ErrorEmptyRegistry, err)

	// test empty stages names error
	r.RegisterPipelineStage("astage", mockStageFun(false, nil))
	p = NewPipelineFromRegistry(r)
	require.Equal(t, rules.ErrorNoStages, p.Err())
	_, _, err = p.Execute(rules.NewBoardState(0, 0), settings.Settings{}, nil)
	require.Equal(t, rules.ErrorNoStages, err)

	// test that an unregistered stage name errors
	p = NewPipelineFromRegistry(r, "doesntexist")
	_, _, err = p.Execute(rules.NewBoardState(0, 0), settings.Settings{}, nil)
	require.Equal(t, rules.ErrorStageNotFound, p.Err())
	require.Equal(t, rules.ErrorStageNotFound, err)

	// simplest case - one stage
	ended, next, err := NewPipelineFromRegistry(r, "astage").Execute(rules.NewBoardState(0, 0), settings.Settings{}, nil)
	require.NoError(t, err)
	require.NoError(t, err)
	require.NotNil(t, next)
	require.False(t, ended)

	// test that the pipeline short-circuits for a stage that errors
	r.RegisterPipelineStage("errors", mockStageFun(false, errors.New("")))
	ended, next, err = NewPipelineFromRegistry(r, "errors", "astage").Execute(rules.NewBoardState(0, 0), settings.Settings{}, nil)
	require.Error(t, err)
	require.NotNil(t, next)
	require.False(t, ended)

	// test that the pipeline short-circuits for a stage that ends
	r.RegisterPipelineStage("ends", mockStageFun(true, nil))
	ended, next, err = NewPipelineFromRegistry(r, "ends", "astage").Execute(rules.NewBoardState(0, 0), settings.Settings{}, nil)
	require.NoError(t, err)
	require.NotNil(t, next)
	require.True(t, ended)

	// test that the pipeline runs normally for multiple stages
	ended, next, err = NewPipelineFromRegistry(r, "astage", "ends").Execute(rules.NewBoardState(0, 0), settings.Settings{}, nil)
	require.NoError(t, err)
	require.NotNil(t, next)
	require.True(t, ended)
}

func TestStageRegistry(t *testing.T) {
	sr := StageRegistry{}

	// register a stage without error
	require.NoError(t, sr.RegisterPipelineStageError("test", mockStageFun(false, nil)))
	require.Contains(t, sr, "test")

	// error on duplicate
	var e rules.RulesetError
	err := sr.RegisterPipelineStageError("test", mockStageFun(false, nil))
	require.Error(t, err)
	require.True(t, errors.As(err, &e), "error should be a RulesetError")
	require.Equal(t, "stage 'test' has already been registered", err.Error())

	// register another stage with no error
	require.NoError(t, sr.RegisterPipelineStageError("other", mockStageFun(false, nil)))
	require.Contains(t, sr, "other")

	// register stage
	sr.RegisterPipelineStage("last", mockStageFun(false, nil))
	require.Contains(t, sr, "last")

	// register existing stage (should just be okay and not panic or anything)
	sr.RegisterPipelineStage("test", mockStageFun(false, nil))
}

func mockStageFun(ended bool, err error) StageFunc {
	return func(b *rules.BoardState, settings settings.Settings, moves []SnakeMove) (bool, error) {
		return ended, err
	}
}
