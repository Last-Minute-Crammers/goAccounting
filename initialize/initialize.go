package initialize

import (
	"context"
	"fmt"
	"goAccounting/global/constant"
	"goAccounting/util"
	"os"
	"path/filepath"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/go-redis/redis/v8"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

type _config struct {
	Mode       constant.ServerMode `yaml:"Mode"`
	Redis      _redis              `yaml:"Redis"`
	Mysql      _mysql              `yaml:"Mysql"`
	Scheduler  _scheduler          `yaml:"Scheduler"`
	Logger     _logger             `yaml:"Logger"`
	System     _system             `yaml:"System"`
	ThirdParty _thirdParty         `yaml:"ThirdParty"`
}

var (
	Database  *gorm.DB
	Config    *_config
	Rdb       *redis.Client
	Cache     util.Cache
	Scheduler *gocron.Scheduler
)

func initConfig() error {
	configFileName := "config.yaml"
	fmt.Printf("RootDir is %s\n", constant.RootDir)
	configPath := filepath.Join(constant.RootDir, configFileName)
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	if Config == nil {
		Config = &_config{}
	}

	err = yaml.Unmarshal(yamlFile, Config)
	if err != nil {
		return err
	}
	return nil
}

// core
func init() {
	fmt.Println("Starting initialization ...")
	var err error
	Config = &_config{
		Redis:      _redis{},
		Mysql:      _mysql{},
		Scheduler:  _scheduler{},
		System:     _system{},
		ThirdParty: _thirdParty{},
		Logger:     _logger{},
	}
	if err = initConfig(); err != nil {
		fmt.Println("Failed to initialize config:", err)
		panic(err)
	}
	group, _ := errgroup.WithContext(context.Background())
	group.Go(Config.Mysql.initializeMysql)
	group.Go(Config.Redis.initializeRedis)
	group.Go(Config.Scheduler.initScheduler)
	if err = group.Wait(); err != nil {
		panic(err)
	}
	fmt.Println("Config loaded successfully")
}

func reconnection[T any](connect func() (T, error), retryTimes int) (result T, err error) {
	result, err = connect()
	if err != nil && retryTimes > 0 {
		time.Sleep(time.Second * 3)
		result, err = reconnection[T](connect, retryTimes-1)
	}
	return result, err
}
