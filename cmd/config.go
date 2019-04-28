package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

const (
	configName = "cobra"
	configType = "json"
	fileName   = configName + "." + configType
)

type config struct{}

var (
	// RootConfig represents the global config
	RootConfig = &config{}
)

func (c config) AWSKey() string {
	return c.promptIfRequired("AWSKey")
}

func (c config) AWSSecretKey() string {
	return c.promptIfRequired("AWSSecretKey")
}

func (c config) Reset() {
	err := os.Remove(c.fullpath())
	onError(err)
}

func (c config) initConfig() {
	home := getHome()
	viper.AddConfigPath(home)
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)

	if err := viper.ReadInConfig(); err != nil {
		viper.WriteConfigAs(c.fullpath())
	}
}

func (c config) home() string {
	home, err := os.UserHomeDir()
	onError(err)
	return home
}

func (c config) fullpath() string {
	return getHome() + "/" + fileName
}

func (c config) promptIfRequired(key string) string {
	if c.requiresKey(key) {
		promptWith := fmt.Sprintf("Enter %s: ", key)
		input := c.prompt(promptWith)

		viper.Set(key, input)
		err := viper.WriteConfig()
		onError(err)
	}

	return viper.GetString(key)
}

func (c config) requiresKey(key string) bool {
	return viper.GetString(key) == ""
}

func (config) prompt(input string) string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(input)
	scanner.Scan()
	text := scanner.Text()
	err := scanner.Err()
	onError(err)

	return text
}
