package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.aporeto.io/remod/internal/remod"
)

func init() {
	cmdDevon.Flags().StringSliceP("include", "m", nil, "Set the prefix of the modules to include")
	cmdDevon.Flags().StringSliceP("exclude", "E", nil, "Set the prefix of the modules to exclude")
	cmdDevon.Flags().StringP("prefix", "p", "../", "The prefix to use for the replacements")
	cmdDevon.Flags().String("replace-version", "", "Set the version to use for replacement. It must be set if prefix is not ../ and must not be if different")
}

var cmdDevon = &cobra.Command{
	Use:     "on",
	Aliases: []string{"devon"},
	Short:   "Apply developpment replace directive",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return viper.BindPFlags(cmd.Flags())
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		included := viper.GetStringSlice("include")
		excluded := viper.GetStringSlice("exclude")
		prefix := viper.GetString("prefix")
		version := viper.GetString("replace-version")

		if !strings.HasPrefix(prefix, ".") && version == "" {
			return fmt.Errorf("you must set --replace-version if --prefix is not local")
		}

		if strings.HasPrefix(prefix, ".") && version != "" {
			return fmt.Errorf("you must not set --replace-version if --prefix is local")
		}

		if err := remod.GitConfig(); err != nil {
			return fmt.Errorf("unable to install git config: %s", err)
		}

		idata, err := ioutil.ReadFile("go.mod")
		if err != nil {
			return fmt.Errorf("unable to read go.mod: %s", err)
		}

		modules, err := remod.Extract(idata, included, excluded)
		if err != nil {
			return fmt.Errorf("unable to extract modules: %s", err)
		}

		odata, err := remod.Enable(idata, modules, prefix, version)
		if err != nil {
			return fmt.Errorf("unable to apply dev replacements: %s", err)
		}
		if odata == nil {
			return nil
		}

		if err := ioutil.WriteFile("go.mod", odata, 0655); err != nil {
			return err
		}

		return remod.GitAdd()
	},
}