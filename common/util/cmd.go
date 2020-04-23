package util

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func BindFlagsToViper(cmd *cobra.Command, _ []string) error {
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}
	return nil
}

func SilenceMsg(cmd *cobra.Command) {
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
	for _, c := range cmd.Commands() {
		c.SilenceUsage = true
		c.SilenceErrors = true
	}
}