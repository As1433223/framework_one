package mysql

import (
	"fmt"
	"github.com/As1433223/framework_one/config"
	"gopkg.in/yaml.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type MysqlConf struct {
	Username string `yaml:"Username"`
	Password string `yaml:"Password"`
	Host     string `yaml:"Host"`
	Port     int    `yaml:"Port"`
	Database string `yaml:"Database"`
}

func InitMysql(servername string, MysqlFunc func(Db *gorm.DB) error) error {
	conf, err := config.GetConfig(servername)
	if err != nil {
		return err
	}
	var (
		data        map[interface{}]interface{}
		MysqlConfig MysqlConf
	)
	yaml.Unmarshal([]byte(conf), &data)
	MysqlConfig.Username = data["MysqlConf"].(map[interface{}]interface{})["Username"].(string)
	MysqlConfig.Password = data["MysqlConf"].(map[interface{}]interface{})["Password"].(string)
	MysqlConfig.Host = data["MysqlConf"].(map[interface{}]interface{})["Host"].(string)
	MysqlConfig.Port = data["MysqlConf"].(map[interface{}]interface{})["Port"].(int)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		MysqlConfig.Username, MysqlConfig.Password, MysqlConfig.Host, MysqlConfig.Port, MysqlConfig.Database,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return err
	}
	clone, _ := db.DB()
	defer clone.Close()
	err = MysqlFunc(db)
	if err != nil {
		return err
	}
	return nil
}
