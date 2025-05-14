package initialize

import (
	"goAccounting/global/constant"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type _config struct {
	Mode       constant.ServerMode `yaml:"Mode"`
	Redis      _redis              `yaml:"Redis"`
	Mysql      _mysql              `yaml:"Mysql"`
	Scheduler  _scheduler          `yaml:Scheduler"`
	Logger     _logger             `yaml:"Logger"`
	System     _system             `yaml:"System"`
	ThirdParty _thirdParty         `yaml:"ThirdParty"`
}

var (
	Database *gorm.DB
	Config   *_config
	Rdb      *redis.Client
)

func reconnection[T any](connect func() (T, error), retryTimes int) (result T, err error) {
	result, err = connect()
	if err != nil && retryTimes > 0 {
		time.Sleep(time.Second * 3)
		result, err = reconnection[T](connect, retryTimes-1)
	}
	return result, err
}
