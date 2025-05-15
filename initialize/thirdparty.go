package initialize

type _thirdParty struct {
	Ai _ai `yaml:"Ai"`
}
type _ai struct {
	Host          string  `yaml:"Host"`
	Port          string  `yaml:"Port"`
	MinSimilarity float32 `yaml:"MinSimilarity"`
}

func (ai _ai) GetPortalSite() string {
	return "http://" + ai.Host + ":" + ai.Port
}
func (ai _ai) IsOpen() bool {
	return ai.Host != "" && ai.Port != ""
}
