package collection

type Users []struct {
	Username string   `yaml:"Username"`
	Password string   `yaml:"Password"`
	KeyFile  string   `yaml:"KeyFile"`
	Hosts    []string `yaml:"Hosts"`
}

type YamlConfig struct {
	RefreshTime int   `yaml:"RefreshTime"`
	Users       Users `yaml:"Users"`
}
