/*
 * Copyright 2026 Xphyr
 */
package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/xphyr/synology-csi/pkg/dsm/webapi"
)

var cmdTarget = &cobra.Command{
	Use:   "target",
	Short: "target API",
	Long:  `DSM target API`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var cmdTargetList = &cobra.Command{
	Use:   "list",
	Short: "list targets",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		dsms, err := ListDsms(DsmId)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		targetInfos := make(map[string][]webapi.TargetInfo)
		for _, dsm := range dsms {
			if err := dsm.Login(); err != nil {
				fmt.Printf("Failed to login to DSM: [%s]. err: %v\n", dsm.Ip, err)
				os.Exit(1)
			}
			infos, err := dsm.TargetList()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			targetInfos[dsm.Ip] = infos
			dsm.Logout()
		}

		tw := tabwriter.NewWriter(os.Stdout, 8, 0, 2, ' ', 0)
		fmt.Fprintf(tw, "%-16s\t", "Host:")
		fmt.Fprintf(tw, "%-52s\t", "Name:")
		fmt.Fprintf(tw, "%-52s\t", "IQN:")
		fmt.Fprintf(tw, "%-10s\t", "Status:")
		fmt.Fprintf(tw, "%-12s\t", "MaxSessions:")
		fmt.Fprintf(tw, "%-16s\t", "TargetID:")
		fmt.Fprintf(tw, "%-16s\t", "MappedLuns:")
		fmt.Fprintf(tw, "\n")
		for ip, v := range targetInfos {
			for _, info := range v {
				fmt.Fprintf(tw, "%-16s\t", ip)
				fmt.Fprintf(tw, "%-52s\t", info.Name)
				fmt.Fprintf(tw, "%-52s\t", info.Iqn)
				fmt.Fprintf(tw, "%-10s\t", info.Status)
				fmt.Fprintf(tw, "%-12d\t", info.MaxSessions)
				fmt.Fprintf(tw, "%-16d\t", info.TargetId)
				fmt.Fprintf(tw, "%-16s\t", info.MappedLuns)
				fmt.Fprintf(tw, "\n")
				_ = tw.Flush()
			}
		}

		fmt.Printf("Success, TargetList()\n")
	},
}

func init() {
	cmdTarget.AddCommand(cmdTargetList)
}
