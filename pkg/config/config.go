package config

import (
	"os"

	"sigs.k8s.io/yaml"
)

type AppConfig struct {
	Name     string             `yaml:"name"  json:"name"`
	Log      LogConfig          `yaml:"log" json:"log"`
	DataBase DataBase           `yaml:"dataBase"`
	Admin    AdminServerConfig  `yaml:"admin" json:"admin"`
	Mqtt     []MqttModuleConfig `yaml:"mqtt" json:"mqtt"`
	Kafka    KafkaClientConfig  `yaml:"kafka" json:"kafka"`
}

type DataBase struct {
	// DriverName indicates database driver name
	// default "sqlite3"
	DriverName string `yaml:"driverName" json:"driverName"`
	// AliasName indicates alias name
	// default "default"
	DbName string `yaml:"dbName" json:"dbName"`
	// DataSource indicates the data source path
	// default "/var/lib/gateway/app.db"
	DataSource string `yaml:"dataSource" json:"dataSource"`
}

type AdminServerConfig struct {
	Enable bool   `yaml:"enable" json:"enable"`
	Name   string `yaml:"name" json:"name"`
	Port   string `yaml:"port" json:"port"`
}

type MqttModuleConfig struct {
	Enable          bool   `yaml:"enable" json:"enable"`
	Name            string `yaml:"name" json:"name"`
	DataUploadTopic string `yaml:"dataUploadTopic" json:"dataUploadTopic"`
	Host            string `yaml:"host" json:"host"`
	Port            string `yaml:"port" json:"port"`
	ClientId        string `yaml:"clientId" json:"clientId"`
	User            string `yaml:"user" json:"user"`
	Passwd          string `yaml:"passwd" json:"passwd"`
}

type KafkaClientConfig struct {
	Enable        bool   `yaml:"enable" json:"enable"`
	Name          string `yaml:"name" json:"name"`
	Host          string `yaml:"host" json:"host"`
	Port          string `yaml:"port" json:"port"`
	ClientId      string `yaml:"clientId" json:"clientId"`
	User          string `yaml:"user" json:"user"`
	Passwd        string `yaml:"passwd" json:"passwd"`
	AutoReconnect bool   `yaml:"autoReconnect" json:"autoReconnect"`
}

// NewDefaultAppConfig new default gateway config
func NewDefaultAppConfig() *AppConfig {
	config := &AppConfig{}
	return config
}

// Parse parse config from file
func (c *AppConfig) Parse(filename string) error {
	if !FileIsExist(filename) {
		filename = "./config.yaml"
	}
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, c)
	if err != nil {
		return err
	}
	return nil
}

func FileIsExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}
