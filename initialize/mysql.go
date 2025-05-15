package initialize

import (
	"goAccounting/global/constant"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type _mysql struct {
	Path     string `yaml:"Path"`
	Port     string `yaml:"Port"`
	Config   string `yaml:"Config"`
	DbName   string `yaml:"DbName"`
	Username string `yaml:"Username"`
	Password string `yaml:"Password"`
}

func (m *_mysql) getDSN() string {
	return m.Username + ":" + m.Password + "@tcp(" + m.Path + ":" + m.Port + ")/" + m.DbName + "?" + m.Config
}
func (m *_mysql) gormConfigInit() *gorm.Config {
	config := &gorm.Config{
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		TranslateError: true,
	}
	switch Config.Mode {
	case constant.Debug:
		config.Logger = logger.Default.LogMode(logger.Info)
	case constant.Production:
		config.Logger = logger.Default.LogMode(logger.Silent)
	default:
		panic("error Mode")
	}
	return config
}

func (m *_mysql) initializeMysql() error {
	var err error
	mysqlConfig := mysql.Config{
		DSN:                       m.getDSN(),
		DefaultStringSize:         191,
		SkipInitializeWithVersion: false,
	}
	var db *gorm.DB
	db, err = reconnection[*gorm.DB](
		func() (*gorm.DB, error) {
			return gorm.Open(mysql.New(mysqlConfig), m.gormConfigInit())
		}, 10,
	)
	if err != nil {
		return err
	}
	sqlDb, _ := db.DB()
	sqlDb.SetMaxIdleConns(50)
	sqlDb.SetMaxOpenConns(50)
	sqlDb.SetConnMaxLifetime(3 * time.Minute)
	db.InstanceSet("gorm:table_options", "ENGINE=InnoDB")
	Database = db
	return nil
}
