/*
 * Copyright 2021 Synology Inc.
 */

package main

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/xphyr/synology-csi/pkg/driver"
	"github.com/xphyr/synology-csi/pkg/dsm/common"
	"github.com/xphyr/synology-csi/pkg/dsm/service"
	"github.com/xphyr/synology-csi/pkg/logger"
	"github.com/xphyr/synology-csi/pkg/utils/hostexec"
)

type Config struct {
	NodeID         string
	Endpoint       string
	ClientInfoPath string
	LogLevel       string
	WebapiDebug    bool
	MultipathForUC bool
	ChrootDir      string
	IscsiadmPath   string
	MultipathPath  string
	MultipathdPath string
}

func driverStart(cfg *Config) error {
	log.Infof("CSI Options = {%s, %s, %s}", cfg.NodeID, cfg.Endpoint, cfg.ClientInfoPath)

	dsmService := service.NewDsmService()

	// 1. Login DSMs by given ClientInfo
	info, err := common.LoadConfig(cfg.ClientInfoPath)
	if err != nil {
		log.Errorf("Failed to read config: %v", err)
		return err
	}

	for _, client := range info.Clients {
		err := dsmService.AddDsm(client)
		if err != nil {
			log.Errorf("Failed to add DSM: %s, error: %v", client.Host, err)
		}
	}
	defer dsmService.RemoveAllDsms()

	// 2. Create command executor
	cmdMap := map[string]string{
		"iscsiadm":   cfg.IscsiadmPath,
		"multipath":  cfg.MultipathPath,
		"multipathd": cfg.MultipathdPath,
	}
	cmdExecutor, err := hostexec.New(cmdMap, cfg.ChrootDir)
	if err != nil {
		log.Errorf("Failed to create command executor: %v", err)
		return err
	}
	tools := driver.NewTools(cmdExecutor)

	// 3. Create and Run the Driver
	drv, err := driver.NewControllerAndNodeDriver(cfg.NodeID, cfg.Endpoint, dsmService, tools)
	if err != nil {
		log.Errorf("Failed to create driver: %v", err)
		return err
	}
	drv.Activate()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	// Block until a signal is received.
	<-c
	log.Infof("Shutting down.")
	return nil
}

func main() {
	cfg := &Config{}
	rootCmd := &cobra.Command{
		Use:          "synology-csi-driver",
		Short:        "Synology CSI Driver",
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, _ []string) error {
			if cfg.WebapiDebug {
				logger.WebapiDebug = true
				cfg.LogLevel = "debug"
			}
			logger.Init(cfg.LogLevel)

			if !cfg.MultipathForUC {
				driver.MultipathEnabled = false
			}

			if err := driverStart(cfg); err != nil {
				log.Errorf("Failed to start driver: %v", err)
				return err
			}
			return nil
		},
	}

	rootCmd.FParseErrWhitelist.UnknownFlags = true
	addFlags(rootCmd, cfg)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}

func addFlags(cmd *cobra.Command, cfg *Config) {
	cmd.PersistentFlags().StringVar(&cfg.NodeID, "nodeid", "CSINode", "Node ID")
	cmd.PersistentFlags().StringVarP(&cfg.Endpoint, "endpoint", "e", "unix:///var/lib/kubelet/plugins/"+driver.DriverName+"/csi.sock", "CSI endpoint")
	cmd.PersistentFlags().StringVarP(&cfg.ClientInfoPath, "client-info", "f", "/etc/synology/client-info.yml", "Path of Synology config yaml file")
	cmd.PersistentFlags().StringVar(&cfg.LogLevel, "log-level", "info", "Log level (debug, info, warn, error, fatal)")
	cmd.PersistentFlags().BoolVarP(&cfg.WebapiDebug, "debug", "d", false, "Enable webapi debugging logs")
	cmd.PersistentFlags().BoolVar(&cfg.MultipathForUC, "multipath", true, "Set to 'false' to disable multipath for UC")
	cmd.PersistentFlags().StringVar(&cfg.ChrootDir, "chroot-dir", "/host", "Host directory to chroot into (empty disables chroot)")
	cmd.PersistentFlags().StringVar(&cfg.IscsiadmPath, "iscsiadm-path", "", "Full path of iscsiadm executable")
	cmd.PersistentFlags().StringVar(&cfg.MultipathPath, "multipath-path", "", "Full path of multipath executable")
	cmd.PersistentFlags().StringVar(&cfg.MultipathdPath, "multipathd-path", "", "Full path of multipathd executable")

	cmd.MarkFlagRequired("endpoint")
	cmd.MarkFlagRequired("client-info")
	cmd.Flags().SortFlags = false
	cmd.PersistentFlags().SortFlags = false
}
