package rulesets

import (
	"rules"
	"rules/settings"
)

type Ruleset interface {
	Name() string
	Settings() settings.Settings
	Execute(prevState *rules.BoardState, moves []SnakeMove) (gameOver bool, nextState *rules.BoardState, err error)
}

type SnakeMove struct {
	ID   string
	Move string
}

type rulesetBuilder struct {
	params   map[string]string  // game customisation parameters
	seed     int64              // used for random events in games
	rand     rules.Rand         // used for random number generation
	solo     bool               // if true, only 1 alive snake is required to keep the game from ending
	settings *settings.Settings // used to set settings directly instead of via string params
}

// NewRulesetBuilder returns an instance of a builder for the Ruleset types.
func NewRulesetBuilder() *rulesetBuilder {
	return &rulesetBuilder{
		params: map[string]string{},
	}
}

// WithParams accepts a map of string parameters for customizing games.
func (rb *rulesetBuilder) WithParams(params map[string]string) *rulesetBuilder {
	for k, v := range params {
		rb.params[k] = v
	}
	return rb
}

// WithSeed sets the seed used for randomisation by certain game modes.
func (rb *rulesetBuilder) WithSeed(seed int64) *rulesetBuilder {
	rb.seed = seed
	return rb
}

// WithRandom overrides the random number generator with a specific instance
// instead of a Rand initialized from the seed.
func (rb *rulesetBuilder) WithRand(rand rules.Rand) *rulesetBuilder {
	rb.rand = rand
	return rb
}

// WithSolo sets whether the ruleset is a solo game.
func (rb *rulesetBuilder) WithSolo(value bool) *rulesetBuilder {
	rb.solo = value
	return rb
}

// WithSettings sets the settings object for the ruleset directly.
func (rb *rulesetBuilder) WithSettings(settings settings.Settings) *rulesetBuilder {
	rb.settings = &settings
	return rb
}

// NamedRuleset constructs a known ruleset by using name to look up a standard pipeline.
func (rb rulesetBuilder) NamedRuleset(name string) Ruleset {
	var stages []string
	if rb.solo {
		stages = append(stages, StageGameOverSoloSnake)
	} else {
		stages = append(stages, StageGameOverStandard)
	}

	if name == rules.GameTypeStandard {
		stages = append(stages, standardRulesetStages[1:]...)
	} else if name == rules.GameTypeSolo {
		stages = soloRulesetStages
	}

	return rb.PipelineRuleset(name, NewPipeline(stages...))
}

// PipelineRuleset constructs a ruleset with the given name and pipeline using the parameters passed to the builder.
// This can be used to create custom rulesets.
func (rb rulesetBuilder) PipelineRuleset(name string, p Pipeline) Ruleset {
	var settingsInstance settings.Settings
	if rb.settings != nil {
		settingsInstance = *rb.settings
	} else {
		settingsInstance = settings.NewSettings(rb.params).WithRand(rb.rand).WithSeed(rb.seed)
	}
	return &pipelineRuleset{
		name:     name,
		pipeline: p,
		settings: settingsInstance,
	}
}

type pipelineRuleset struct {
	pipeline Pipeline
	name     string
	settings settings.Settings
}

// impl Ruleset
func (r pipelineRuleset) Settings() settings.Settings {
	return r.settings
}

// impl Ruleset
func (r pipelineRuleset) Name() string { return r.name }

// impl Ruleset
func (r pipelineRuleset) Execute(bs *rules.BoardState, sm []SnakeMove) (bool, *rules.BoardState, error) {
	return r.pipeline.Execute(bs, r.Settings(), sm)
}

func (r pipelineRuleset) Err() error {
	return r.pipeline.Err()
}
