package configuration

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"os"
	"transact-api/constants"
)

// Configuration interface for application configuration
type Configuration interface {
	GetAppName() string
	GetPort() string
	GetReadTimeout() int32
	GetWriteTimeout() int32
	GetMaxIdleTimeOut() int32
	GetServerTimeOut() int
	GetDBHost() string
	GetDBPort() int
	GetDBName() string
	GetDBUser() string
	GetDBPassword() string
	GetAppMode() string
	GetAppHost() string
}

type configuration struct {
	AppConfig
	DBConfig `json:"db_config"`
}

func (c configuration) GetDBName() string {
	return c.DBName
}

func (c configuration) GetAppHost() string {
	return c.AppHost
}

func (c configuration) GetAppMode() string {
	return c.AppMode
}

func (c configuration) GetAppName() string {
	return c.AppName
}

func (c configuration) GetPort() string {
	return c.AppConfig.Port
}

func (c configuration) GetReadTimeout() int32 {
	return c.ReadTimeout
}

func (c configuration) GetWriteTimeout() int32 {
	return c.WriteTimeout
}

func (c configuration) GetMaxIdleTimeOut() int32 {
	return c.MaxIdleTimeOut
}

func (c configuration) GetServerTimeOut() int {
	return c.ServerTimeOut
}

func (c configuration) GetDBHost() string {
	return c.Host
}

func (c configuration) GetDBPort() int {
	return c.DBConfig.Port
}

func (c configuration) GetDBUser() string {
	return c.DBConfig.User
}

func (c configuration) GetDBPassword() string {
	return os.Getenv(constants.DbPassword)
}

// AppConfig holds the application configuration
type AppConfig struct {
	AppName        string `json:"app_name"`
	AppMode        string `json:"app_mode"`
	AppHost        string `json:"app_host"`
	Port           string `json:"port"`
	ReadTimeout    int32  `json:"read_timeout"`
	WriteTimeout   int32  `json:"write_timeout"`
	MaxIdleTimeOut int32  `json:"max_idle_time_out"`
	ServerTimeOut  int    `json:"server_timeout"`
}

// DBConfig holds the database configuration
type DBConfig struct {
	Host   string `json:"host"`
	Port   int    `json:"port"`
	User   string `json:"user"`
	DBName string `json:"dbname"`
}

func Parse(file string) (Configuration, error) {
	var Cfg configuration
	data, err := os.ReadFile(file)
	if err != nil {
		log.Errorf("config: ioutil.ReadFile failed: %s", err.Error())
		return nil, err
	}

	if err = json.Unmarshal(data, &Cfg); err != nil {
		log.Errorf("config: json.unmarshal failed: %s", err.Error())
		return nil, err
	}

	return Cfg, nil
}
