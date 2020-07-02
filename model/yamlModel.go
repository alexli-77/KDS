package model

type PolicyYaml struct {
	ApiVersion    string `yaml:"apiVersion"`
	Kind     string `yaml:"kind"`
	Strategies      Strategies `yaml:"strategies"`
}

type Strategies struct {
	RemoveDuplicates    RemoveDuplicates `yaml:"RemoveDuplicates"`
	RemovePodsViolatingInterPodAntiAffinity     RemovePodsViolatingInterPodAntiAffinity `yaml:"RemovePodsViolatingInterPodAntiAffinity"`
	LowNodeUtilization      LowNodeUtilization `yaml:"LowNodeUtilization"`
}

type RemoveDuplicates struct {
	Enabled    bool `yaml:"enabled"`
}

type RemovePodsViolatingInterPodAntiAffinity struct {
	Enabled    bool `yaml:"enabled"`
}

type LowNodeUtilization struct {
	Enabled    bool `yaml:"enabled"`
	Params     LowNodeUtilizationParams `yaml:"params"`
}

type IsEnabled struct {
	Enabled    bool `yaml:"enabled"`
}

type LowNodeUtilizationParams struct {
	NodeResourceUtilizationThresholds    NodeResourceUtilizationThresholds `yaml:"nodeResourceUtilizationThresholds"`

}

type NodeResourceUtilizationThresholds struct {
	Thresholds    Thresholds `yaml:"thresholds"`
	TargetThresholds TargetThresholds `yaml:"targetThresholds"`
}

type Thresholds struct {
	Cpu      float64 `yaml:"cpu"`
	Memory   float64 `yaml:"memory"`
	Pods     float64 `yaml:"pods"`
}

type TargetThresholds struct {
	Cpu      float64 `yaml:"cpu"`
	Memory   float64 `yaml:"memory"`
	Pods     float64 `yaml:"pods"`
}

type Config struct {
	PolicyFilepath  string `yaml:"policyfilepath"`
	TemplateFilepath  string `yaml:"templatefilepath"`
}
