package core

const (
	appName        = "GoSteamRestarter"
	configFileName = "config.json"
)

type Config struct {
	SteamPath string `json:"steamPath"`
	SteamArgs string `json:"steamArgs"`
}

func DefaultConfig() Config {
	return Config{}
}
