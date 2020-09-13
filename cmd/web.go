/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"context"

	"net"

	"net/http"
	_ "net/http/pprof"

	"github.com/czhj/ahfs/modules/log"
	"github.com/czhj/ahfs/modules/setting"
	"github.com/czhj/ahfs/routers"
	"github.com/czhj/ahfs/routers/routes"
	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/crypto/acme/autocert"
)

// webCmd represents the web command
var webCmd = &cobra.Command{
	Use:   "web",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runWeb(cmd, args)
	},
}

func runLetsEncryptWeb(e *gin.Engine, listenAddr, domain, directory string) error {
	m := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(domain),
		Cache:      autocert.DirCache(directory),
	}

	return autotls.RunWithManager(e, m)
}

func runHTTPS(e *gin.Engine, listenAddr, certFile, keyFile string) error {
	return e.RunTLS(listenAddr, certFile, keyFile)
}

func runHTTP(e *gin.Engine, listenAddr string) error {
	return e.Run(listenAddr)
}

func runWeb(cmd *cobra.Command, args []string) error {
	defer log.Sync()

	rootCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	routers.GlobalInit(rootCtx)

	e := routes.NewEngine()
	routes.RegisterRoutes(e)

	setting.SaveSetting()

	if setting.PprofService.Enabled {
		go func() {
			addr := net.JoinHostPort(setting.PprofService.HTTPAddr, setting.PprofService.HTTPPort)
			if err := http.ListenAndServe(addr, nil); err != nil {
				log.Error("Pprof service is shutdown", zap.Error(err))
			}
		}()
	}

	listenAddr := net.JoinHostPort(setting.HTTPAddr, setting.HTTPPort)

	var err error
	switch setting.Protocol {
	case "http":
		err = runHTTP(e, listenAddr)
	case "https":
		if setting.EnableLetsEncrypt {
			err = runLetsEncryptWeb(e, listenAddr, setting.Domain, setting.LetsEncryptDir)
			break
		}
		err = runHTTPS(e, listenAddr, setting.CertFile, setting.KeyFile)
	default:
		log.Fatal("Invalid protocol", zap.String("protocol", setting.Protocol))
	}

	if err != nil {
		log.Fatal("Failed to start server", zap.Error(err))
	}

	log.Info("HTTP Listener is shutdown", zap.String("addr", listenAddr))
	return nil
}

func init() {
	rootCmd.AddCommand(webCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// webCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// webCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
