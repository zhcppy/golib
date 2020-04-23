package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zhcppy/golib/common/util"
	"github.com/zhcppy/golib/logger"
	"github.com/zhcppy/golib/server"
	"github.com/zhcppy/golib/version"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if util.RunEnv() == util.EnvDev {
		if err != nil {
			logger.L.Errorf("read config file error: %s", err.Error())
		} else {
			logger.L.Debugf("config file found and successfully parsed, %s\n", viper.ConfigFileUsed())
		}
	}
}

func main() {
	app := &cobra.Command{
		Short:             "This is demo",
		PersistentPreRunE: util.BindFlagsToViper,
		RunE: func(*cobra.Command, []string) error {
			var cfg = &server.Config{}
			if err := viper.Unmarshal(cfg); err != nil {
				return err
			}
			logger.L.Debug(cfg.String())

			if viper.GetBool("init") {
				if err := server.InitDB(cfg.DbConfig); err != nil {
					return err
				}
			}

			appServer, err := server.NewServer(cfg)
			if err != nil {
				return err
			}
			return appServer.Run(":8899")
		},
	}
	app.Flags().Bool("init", false, "create mysql db")
	app.AddCommand(version.NewCmd())
	util.SilenceMsg(app)
	if err := app.Execute(); err != nil {
		fmt.Printf("\033[1;31m%s\033[0m", fmt.Sprintf("Failed to command execute: %s\n", err.Error()))
	}
}
