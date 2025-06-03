package initialize

import (
	"goAccounting/global/config"
)

var Config = config.Config{
	System: config.System{
		Addr: 8080,
		DbType: "sqlite",
		Mode: "debug",
	},
}
