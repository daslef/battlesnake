package settings_test

import (
	"rules/settings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSettings(t *testing.T) {
	params := map[string]string{
		"invalidSetting": "abcd",
		"intSetting":     "1234",
		"boolSetting":    "true",
	}

	testSettings := settings.NewSettings(params)

	assert.Equal(t, 4567, testSettings.Int("missingIntSetting", 4567))
	assert.Equal(t, 4567, testSettings.Int("invalidSetting", 4567))
	assert.Equal(t, 1234, testSettings.Int("intSetting", 4567))

	assert.Equal(t, false, testSettings.Bool("missingBoolSetting", false))
	assert.Equal(t, true, testSettings.Bool("missingBoolSetting", true))
	assert.Equal(t, false, testSettings.Bool("invalidSetting", true))
	assert.Equal(t, true, testSettings.Bool("boolSetting", true))

	assert.Equal(t, 4567, settings.NewSettingsWithParams("newIntSetting").Int("newIntSetting", 4567))
	assert.Equal(t, 1234, settings.NewSettingsWithParams("newIntSetting", "1234").Int("newIntSetting", 4567))
	assert.Equal(t, 4567, settings.NewSettingsWithParams("x", "y", "newIntSetting").Int("newIntSetting", 4567))
}
