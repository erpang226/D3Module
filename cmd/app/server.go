package app

import (
	"fmt"
	"main/cmd/app/options"
	"main/common/global"
	actor "main/pkg/actor/core"
	adminModule "main/pkg/admin/module"
	"main/pkg/common/dbm"
	"main/pkg/common/dbm/dao/core"
	"main/pkg/common/dbm/dao/module"
	"main/pkg/common/dbm/dao/moduleproperty"
	"main/pkg/config"
	"main/pkg/north"

	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

// NewGatewayCommand create app cmd
func NewGatewayCommand() *cobra.Command {
	appOptions := options.InitGatewayOptions(options.ConfigPath)
	c, err := appOptions.Config()
	if err != nil {
		panic(err)
	}
	options.SetAppConfig(c)
	// init log
	err = c.Log.LogInit()
	if err != nil {
		_ = fmt.Errorf("set log flag error. %v", err)
	}
	defer c.Log.FlushLogs()

	// init db
	core.RegisterCoreDBTable()
	dbm.InitDBConfig(c.DataBase.DriverName, c.DataBase.DbName, c.DataBase.DataSource)
	cmd := &cobra.Command{
		Use:  "app",
		Long: `D3Module app...`,
		Run: func(cmd *cobra.Command, args []string) {
			err := environmentCheck(false)
			if err != nil {
				klog.Error("environment not satisfied.", err)
			}
			registerModules(options.GetAppConfig())
			// start all modules
			actor.Run()
		},
	}

	cmd.Flags().StringVar(&options.ConfigPath, "config", global.DefaultConfigDir+global.DefaultConfigFile,
		fmt.Sprintf("app config file path, the default value is %s", global.DefaultConfigDir+global.DefaultConfigFile))
	return cmd
}

// environmentCheck check the environment before start
// if Check failed,  return errors
func environmentCheck(skipCheck bool) error {
	if skipCheck {
		return nil
	}
	// ...
	return nil
}

// registerModules register all the modules to be started
func registerModules(appConfig *config.AppConfig) {
	registerAdminModule(appConfig)
	registerMiddleModules(appConfig)
	registerNorthModules(appConfig)
	registerSouthModules(appConfig)
}

func registerNorthModules(appConfig *config.AppConfig) {
	// init modules according to db
	rows, queryErr := module.Query("type", "north")
	if queryErr != nil {
		klog.Errorf("Query sqlite failed while syncing sqlite, err: %#v", queryErr)
	}
	if rows == nil {
		klog.Info("Query sqlite nil while syncing sqlite")
		return
	}
	for _, m := range *rows {
		if !m.Enable {
			continue
		}
		property, err := moduleproperty.QueryModuleProperty(m.Id)
		if err != nil {
			klog.Error(err)
			continue
		}
		north.RegisterModuleFromDB(&m, property)
	}
}

func registerMiddleModules(appConfig *config.AppConfig) {
	// init modules according to db
	rows, queryErr := module.Query("type", "middle")
	if queryErr != nil {
		klog.Errorf("Query sqlite failed while syncing sqlite, err: %#v", queryErr)
	}
	if rows == nil {
		klog.Info("Query sqlite nil while syncing sqlite")
		return
	}
	for _, m := range *rows {
		if !m.Enable {
			continue
		}
	}
}

func registerSouthModules(appConfig *config.AppConfig) {
	// init modules according to db
	rows, queryErr := module.Query("type", "south")
	if queryErr != nil {
		klog.Errorf("Query sqlite failed while syncing sqlite, err: %#v", queryErr)
	}
	if rows == nil {
		klog.Info("Query sqlite nil while syncing sqlite")
		return
	}
	for _, m := range *rows {
		if !m.Enable {
			continue
		}
	}

}

func registerAdminModule(appConfig *config.AppConfig) {
	if appConfig.Admin.Enable {
		adminModule.Register(&appConfig.Admin)
	}
}
