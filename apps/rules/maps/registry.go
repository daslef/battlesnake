package maps

import (
	"rules"
)

type MapRegistry map[string]GameMap

var globalRegistry = MapRegistry{
	"standard": StandardMap{},
}

// GetMap returns the map associated with the given ID.
func (registry MapRegistry) GetMap(id string) (GameMap, error) {
	if m, ok := registry[id]; ok {
		return m, nil
	}
	return nil, rules.ErrorMapNotFound
}

// GetMap returns the map associated with the given ID from the global registry.
func GetMap(id string) (GameMap, error) {
	return globalRegistry.GetMap(id)
}

func TestMap(id string, m GameMap, callback func()) {
	globalRegistry[id] = m
	callback()
	delete(globalRegistry, id)
}
