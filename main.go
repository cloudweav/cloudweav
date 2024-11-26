//go:generate go run pkg/codegen/cleanup/main.go
//go:generate /bin/rm -rf pkg/generated
//go:generate go run pkg/codegen/main.go
//go:generate /bin/bash scripts/generate-manifest
//go:generate /bin/bash scripts/generate-openapi

package main

import (
	"fmt"
	_ "net/http/pprof"

	"github.com/rancher/wrangler/v3/pkg/signals"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/cloudweav/cloudweav/pkg/cmd"
	"github.com/cloudweav/cloudweav/pkg/config"
	"github.com/cloudweav/cloudweav/pkg/server"
)

func main() {
	var options config.Options

	flags := []cli.Flag{
		cli.IntFlag{
			Name:        "threadiness",
			EnvVar:      "THREADINESS",
			Usage:       "Specify controller threads",
			Value:       10,
			Destination: &options.Threadiness,
		},
		cli.IntFlag{
			Name:        "http-port",
			EnvVar:      "CLOUDWEAV_SERVER_HTTP_PORT",
			Usage:       "HTTP listen port",
			Value:       8080,
			Destination: &options.HTTPListenPort,
		},
		cli.IntFlag{
			Name:        "https-port",
			EnvVar:      "CLOUDWEAV_SERVER_HTTPS_PORT",
			Usage:       "HTTPS listen port",
			Value:       8443,
			Destination: &options.HTTPSListenPort,
		},
		cli.StringFlag{
			Name:        "namespace",
			EnvVar:      "NAMESPACE",
			Destination: &options.Namespace,
			Usage:       "The default namespace to store management resources",
			Required:    true,
		},
		cli.BoolFlag{
			Name:        "hci-mode",
			EnvVar:      "HCI_MODE",
			Usage:       "Enable HCI mode. Additional controllers are registered in HCI mode",
			Destination: &options.HCIMode,
		},
		cli.BoolFlag{
			Name:        "rancher-embedded",
			EnvVar:      "RANCHER_EMBEDDED",
			Usage:       "Specify whether the Cloudweav is running with embedded Rancher mode, default to false",
			Destination: &options.RancherEmbedded,
		},
		cli.StringFlag{
			Name:        "rancher-server-url",
			EnvVar:      "RANCHER_SERVER_URL",
			Usage:       "Specify the URL to connect to the Rancher server",
			Destination: &options.RancherURL,
			Hidden:      true,
		},
	}

	app := cmd.NewApp("Cloudweav API Server", "", flags, func(commonOptions *config.CommonOptions) error {
		return run(commonOptions, options)
	})
	app.Run()
}

func run(commonOptions *config.CommonOptions, options config.Options) error {
	logrus.Info("Starting controller")
	ctx := signals.SetupSignalContext()

	kubeConfig, err := server.GetConfig(commonOptions.KubeConfig)
	if err != nil {
		return fmt.Errorf("failed to find kubeconfig: %v", err)
	}

	harv, err := server.New(ctx, kubeConfig, options)
	if err != nil {
		return fmt.Errorf("failed to create cloudweav server: %v", err)
	}
	return harv.ListenAndServe(nil, options)
}
