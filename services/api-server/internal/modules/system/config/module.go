package config

type ModuleConfig struct {
	DependencyNames []string
}

func DefaultModuleConfig() ModuleConfig {
	return ModuleConfig{
		DependencyNames: []string{
			"postgres",
			"redis",
			"object_storage",
		},
	}
}
