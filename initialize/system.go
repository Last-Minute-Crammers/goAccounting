package initialize

type _system struct {
	// addr means running port
	Addr         int    `yaml:"Addr"`
	RouterPrefix string `yaml:"RouterPrefix"`

	JwtKey        string `yaml:"JwtKey"`
	ClientSignKey string `yaml:"ClientSignKey"`
}
