package cmd

import (
	"github.com/sirupsen/logrus"
	"os"
)

func handleErr(err error) {
	if err != nil {
		logrus.Error(err)
		os.Exit(2)
	}
}
