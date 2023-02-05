package domain

type Targets struct {
	Urls []struct {
		Link             string   `yaml:"link"`
		Selectors        []string `yaml:"selectors"`
		ExcludeSelectors []string `yaml:"excludeSelectors"`
	} `yaml:"urls"`
	Themes []string `yaml:"themes"`
}
