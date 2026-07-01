package darwin

type Platform struct{}

func New() Platform {
	return Platform{}
}

func (Platform) DefaultSteamPath() string {
	return "/Applications/Steam.app"
}

func (p Platform) FindSteamPath() (string, error) {
	return p.DefaultSteamPath(), nil
}

func (Platform) KillSteam() error {
	return nil
}

func (Platform) StartSteam(path, args string) error {
	return nil
}

func (Platform) FlushDNS() error {
	return nil
}
