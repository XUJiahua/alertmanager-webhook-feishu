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
var splitByStatus bool
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
			logrus.Infof("bot %s created", group)
		}

		// start server
		s := server.New(bots, splitByStatus)
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
	serverCmd.Flags().BoolVarP(&splitByStatus, "split", "", false, "if enabled, sending firing and resolved alerts in two notifications")
	serverCmd.Flags().StringVarP(&cfgFile, "config", "c", "", "config file for bot webhook")
}
