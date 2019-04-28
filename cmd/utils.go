package cmd

import (
	"fmt"
	"os"
)

func getHome() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return home
}

func onError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
