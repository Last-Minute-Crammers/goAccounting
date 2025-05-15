package global

import (
	"fmt"
	"goAccounting/initialize"
)

var (
	GlobalDb  = initialize.Database
	GlobalRdb = initialize.Rdb
	Config    = initialize.Config
	Cache     = initialize.Cache
)

func init() {
	fmt.Println("[global]: starting init")
	GlobalDb = initialize.Database
	GlobalRdb = initialize.Rdb
	Config = initialize.Config
	Cache = initialize.Cache

	if GlobalDb == nil || GlobalRdb == nil || Config == nil || Cache == nil {
		panic("global: failed to get resources from package initialize")
	}
	fmt.Println("[global]: init success")
}
