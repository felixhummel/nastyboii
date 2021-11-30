/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"

	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var Directory string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nastyboii",
	Short: "reminds you of nasty stuff",
	Run: func(cmd *cobra.Command, args []string) {
		err := filepath.Walk(".", func(path_ string, info fs.FileInfo, err error) error {
			if err != nil {
				fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path_, err)
				return err
			}
			if info.IsDir() {
				name := info.Name()
				repolog := log.WithFields(log.Fields{
					"repo": name,
				})
				_, err := os.Stat(path.Join(name, ".git"))
				if err != nil {
					// debug log: "not a git repo: %s", name
					return filepath.SkipDir
				}
				repolog.Info("found")
			}
			return nil
		})
		if err != nil {
			fmt.Printf("error walking the path %q: %v\n", '.', err)
			return
		}
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&Directory, "directory", "", "where to look for git repos")
	viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
	viper.SetDefault("license", "apache")
}

func initConfig() {
	// Find home directory.
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	// Search config in home directory with name ".nastyboii" (without extension).
	viper.AddConfigPath(home)
	viper.SetConfigType("yaml")
	viper.SetConfigName(".nastyboii")

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
