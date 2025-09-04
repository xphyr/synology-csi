package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/xphyr/synology-csi/pkg/dsm/webapi"
)

var cmdHost = &cobra.Command{
	Use:   "host",
	Short: "host API",
	Long:  `DSM host API`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var cmdHostList = &cobra.Command{
	Use:   "list",
	Short: "list hosts",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		dsms, err := ListDsms(DsmId)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		hostInfos := make(map[string][]webapi.HostInfo)
		for _, dsm := range dsms {
			if err := dsm.Login(); err != nil {
				fmt.Printf("Failed to login to DSM: [%s]. err: %v\n", dsm.Ip, err)
				os.Exit(1)
			}
			infos, err := dsm.HostList()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println(infos)
			hostInfos[dsm.Ip] = infos
			dsm.Logout()
		}

		tw := tabwriter.NewWriter(os.Stdout, 8, 0, 2, ' ', 0)
		fmt.Fprintf(tw, "%-16s\t", "DSM:")
		fmt.Fprintf(tw, "%-16s\t", "HostName:")
		fmt.Fprintf(tw, "%-36s\t", "Uuid:")
		fmt.Fprintf(tw, "%-12s\t", "Description:")
		fmt.Fprintf(tw, "%-58s\t", "IQN:")
		fmt.Fprintf(tw, "%-8s\t", "HostID:")
		fmt.Fprintf(tw, "%-8s\t", "Protocol:")
		fmt.Fprintf(tw, "\n")
		for ip, v := range hostInfos {
			for _, info := range v {
				fmt.Fprintf(tw, "%-16s\t", ip)
				fmt.Fprintf(tw, "%-16s\t", info.Name)
				fmt.Fprintf(tw, "%-36s\t", info.Uuid)
				fmt.Fprintf(tw, "%-12s\t", info.Description)
				fmt.Fprintf(tw, "%-58s\t", info.InitiatorIDs)
				fmt.Fprintf(tw, "%-8d\t", info.HostID)
				fmt.Fprintf(tw, "%-8s\t", info.Protocol)
				fmt.Fprintf(tw, "\n")
				_ = tw.Flush()
			}
		}

		fmt.Printf("Success, HostList()\n")
	},
}

func init() {
	cmdHost.AddCommand(cmdHostList)
}
