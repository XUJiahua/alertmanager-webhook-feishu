/*
Copyright © 2021 xujiahua <littleguner@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/xujiahua/alertmanager-webhook-feishu/config"
	"github.com/xujiahua/alertmanager-webhook-feishu/feishu"
	"github.com/xujiahua/alertmanager-webhook-feishu/server"
	"os"
	"os/signal"
	"syscall"
)

var port int
var emailEnabled bool
var cfgFile string

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "start webhook server",
	Run: func(cmd *cobra.Command, args []string) {
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})

		if verbose {
			logrus.SetReportCaller(true)
			logrus.SetLevel(logrus.DebugLevel)
		}

		// cfg
		cfg, err := config.Load(cfgFile)
		handleErr(err)

		// email helper
		var emailHelper *feishu.EmailHelper
		if emailEnabled {
			var err error
			emailHelper, err = feishu.NewEmailHelper(cfg.App)
			handleErr(err)
		}

		// bots
		bots := make(map[string]feishu.IBot)
		for group, botCfg := range cfg.Bots {
			bot, err := feishu.New(botCfg, emailHelper)
			handleErr(err)
			bots[group] = bot
		}

		// start server
		s := server.New(bots)
		go func() {
			address := fmt.Sprintf("0.0.0.0:%d", port)
			handleErr(s.Start(address))
		}()

		signalChan := make(chan os.Signal, 1)
		signal.Notify(
			signalChan,
			syscall.SIGHUP,  // kill -SIGHUP XXXX
			syscall.SIGINT,  // kill -SIGINT XXXX or Ctrl+c
			syscall.SIGQUIT, // kill -SIGQUIT XXXX
		)
		<-signalChan
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().IntVarP(&port, "port", "p", 8000, "server port")
	serverCmd.Flags().BoolVarP(&emailEnabled, "email", "e", false, "if email supported, need feishu appid/secret for enabling")
	serverCmd.Flags().StringVarP(&cfgFile, "config", "c", "", "config file for bot webhook")
}
