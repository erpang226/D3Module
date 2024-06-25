package options

import (
	"main/common/global"
	"main/pkg/config"
	"path"
)

type GatewayOptions struct {
	ConfigFile string
}

var ConfigPath string

var appOptions *GatewayOptions

var appConfig *config.AppConfig

func GetGatewayOptions() *GatewayOptions {
	return appOptions
}
func GetAppConfig() *config.AppConfig {
	return appConfig
}
func SetAppConfig(c *config.AppConfig) {
	appConfig = c
}

func InitGatewayOptions(configFile string) *GatewayOptions {
	if len(configFile) == 0 {
		configFile = path.Join(global.DefaultConfigDir, "app.yaml")
	}
	appOptions = &GatewayOptions{
		ConfigFile: configFile,
	}
	return appOptions
}

func (o *GatewayOptions) Flags() {
	//fs := fss.FlagSet("global")
	//fs.StringVar(&o.ConfigFile, "config", o.ConfigFile, "The path to the configuration file. Flags override values in this file.")
	//return
}

func (o *GatewayOptions) Validate() []error {
	var errs []error
	return errs
}

func (o *GatewayOptions) Config() (*config.AppConfig, error) {
	appConfig = config.NewDefaultAppConfig()
	if err := appConfig.Parse(o.ConfigFile); err != nil {
		return nil, err
	}

	return appConfig, nil
}
