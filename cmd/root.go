package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var Directory string
var exitCode int

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nastyboii",
	Short: "reminds you of nasty stuff",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.WarnLevel)
		log.SetFormatter(&log.TextFormatter{
			// double-click a repo and paste :)
			ForceQuote: true,
		})

		wd, _ := os.Getwd()
		log.WithFields(log.Fields{
			"wd": wd,
		}).Debug()
		err := filepath.Walk(".", func(path_ string, info fs.FileInfo, err error) error {
			// log.WithField("path", path_).Debug()
			if err != nil {
				fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path_, err)
				return err
			}
			if !info.IsDir() {
				return nil
			}
			name := info.Name()
			repolog := log.WithFields(log.Fields{
				"repo": path_,
			})
			_, err = os.Stat(path.Join(name, ".git"))
			if err != nil {
				repolog.Debug("not a git repo")
				return nil
			}
			repolog.Info("found")
			out, err := exec.Command("git", "-C", path_, "status", "-s").Output()
			if err != nil {
				log.Error(err)
				return nil
			}
			if len(out) > 0 {
				repolog.Debug(out)
				repolog.Warn("nasty boii!")
				exitCode = 1
			}
			return nil
		})
		if err != nil {
			fmt.Printf("error walking the path %q: %v\n", '.', err)
			return
		}
		os.Exit(exitCode)
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
