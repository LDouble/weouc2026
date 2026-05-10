package config

import "strings"

type ModuleConfig struct {
	AllowedScenes map[string]struct{}
}

func DefaultModuleConfig() ModuleConfig {
	return ModuleConfig{
		AllowedScenes: map[string]struct{}{
			"general":  {},
			"market":   {},
			"errand":   {},
			"resource": {},
		},
	}
}

func (c ModuleConfig) NormalizeScene(scene string) string {
	scene = strings.TrimSpace(strings.ToLower(scene))
	if scene == "" {
		return "general"
	}
	if _, ok := c.AllowedScenes[scene]; ok {
		return scene
	}

	return "general"
}
