package config

type ModuleConfig struct {
	DependencyNames []string
}

func DefaultModuleConfig() ModuleConfig {
	return ModuleConfig{
		DependencyNames: []string{
			"mysql",
			"mongo",
			"redis",
			"object_storage",
		},
	}
}
