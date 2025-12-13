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

type _admin struct {
	Emails string `yaml:"emails"` // 逗号分隔的邮箱列表
}

type _config struct {
	Env        string              `yaml:"Env"`
	Mode       constant.ServerMode `yaml:"Mode"`
	Redis      _redis              `yaml:"Redis"`
	Mysql      _mysql              `yaml:"Mysql"`
	Scheduler  _scheduler          `yaml:"Scheduler"`
	Logger     _logger             `yaml:"Logger"`
	System     _system             `yaml:"System"`
	ThirdParty _thirdParty         `yaml:"ThirdParty"`
	Admin      _admin              `yaml:"Admin"`
}

var (
	Database  *gorm.DB
	Config    *_config
	Rdb       *redis.Client
	Cache     util.Cache
	Scheduler *gocron.Scheduler
)

func initConfig() error {
	// 第一步：读取主配置文件获取运行环境
	mainConfigPath := filepath.Join(constant.RootDir, constant.ConfigFileName)
	yamlFile, err := os.ReadFile(mainConfigPath)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", constant.ConfigFileName, err)
	}

	// 解析主配置以获取 Env
	mainConfig := &_config{}
	err = yaml.Unmarshal(yamlFile, mainConfig)
	if err != nil {
		return fmt.Errorf("failed to unmarshal %s: %w", constant.ConfigFileName, err)
	}

	// 第二步：根据 Env 读取对应的配置文件
	var envMode constant.EnvMode
	env := mainConfig.Env

	switch env {
	case string(constant.EnvLocalhost):
		envMode = constant.EnvLocalhost
	case string(constant.EnvDocker):
		envMode = constant.EnvDocker
	default:
		return fmt.Errorf("unknown Env: %s, expected '%s' or '%s'",
			env, constant.EnvLocalhost, constant.EnvDocker)
	}

	actualConfigFile := constant.GetConfigFileName(envMode)
	actualConfigPath := filepath.Join(constant.RootDir, actualConfigFile)
	fmt.Printf("Loading config from: %s\n", actualConfigFile)

	actualYamlFile, err := os.ReadFile(actualConfigPath)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", actualConfigFile, err)
	}

	// 第三步：解析实际配置文件到 Config
	if Config == nil {
		Config = &_config{}
	}

	err = yaml.Unmarshal(actualYamlFile, Config)
	if err != nil {
		return fmt.Errorf("failed to unmarshal %s: %w", actualConfigFile, err)
	}

	// 第四步：用主配置的 Admin 覆盖环境配置的 Admin（实现 Admin 统一管理）
	Config.Admin = mainConfig.Admin

	fmt.Printf("Config loaded successfully from: %s (Env: %s, Mode: %s)\n", actualConfigFile, env, Config.Mode)
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
