package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/xphyr/synology-csi/pkg/dsm/webapi"
)

var cmdDNS = &cobra.Command{
	Use:   "dns",
	Short: "dns API",
	Long:  `DSM dns API`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var cmdDNSList = &cobra.Command{
	Use:   "record",
	Short: "list record",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		dsms, err := ListDsms(DsmId)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		var dnsRecords []webapi.DNSRecord

		for _, dsm := range dsms {
			if err := dsm.Login(); err != nil {
				fmt.Printf("Failed to login to DSM: [%s]. err: %v\n", dsm.Ip, err)
				os.Exit(1)
			}
			records, err := dsm.RecordList([]string{"xphyrlab.net", "home.xphyr.net"}, "master")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			//fmt.Println(dnsRecords)
			dnsRecords = records
			dsm.Logout()
		}

		tw := tabwriter.NewWriter(os.Stdout, 8, 0, 2, ' ', 0)
		fmt.Fprintf(tw, "%-36s\t", "Record:")
		fmt.Fprintf(tw, "%-12s\t", "TTL:")
		fmt.Fprintf(tw, "%-6s\t", "Type")
		fmt.Fprintf(tw, "%-36s\t", "Value")
		fmt.Fprintf(tw, "%-36s\t", "FullRecord:")
		fmt.Fprintf(tw, "%-16s\t", "ZoneName:")
		fmt.Fprintf(tw, "%-16s\t", "DomainName:")
		fmt.Fprintf(tw, "\n")
		for _, dnsRecord := range dnsRecords {
			fmt.Fprintf(tw, "%-36s\t", dnsRecord.Record)
			fmt.Fprintf(tw, "%-12s\t", dnsRecord.TTL)
			fmt.Fprintf(tw, "%-6s\t", dnsRecord.Type)
			fmt.Fprintf(tw, "%-36s\t", dnsRecord.Value)
			fmt.Fprintf(tw, "%-36s\t", strings.Replace(dnsRecord.FullRecord, "\t", ",", -1))
			fmt.Fprintf(tw, "%-16s\t", dnsRecord.ZoneName)
			fmt.Fprintf(tw, "%-16s\t", dnsRecord.DomainName)
			fmt.Fprintf(tw, "\n")
			_ = tw.Flush()

		}

		fmt.Printf("Success, DNSrecordsList()\n")
	},
}

var cmdZoneList = &cobra.Command{
	Use:   "zone",
	Short: "list zone",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		dsms, err := ListDsms(DsmId)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		var zoneRecords []webapi.ZoneRecord

		for _, dsm := range dsms {
			if err := dsm.Login(); err != nil {
				fmt.Printf("Failed to login to DSM: [%s]. err: %v\n", dsm.Ip, err)
				os.Exit(1)
			}
			records, err := dsm.ZoneList()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			//fmt.Println(dnsRecords)
			zoneRecords = records
			dsm.Logout()
		}

		tw := tabwriter.NewWriter(os.Stdout, 8, 0, 2, ' ', 0)
		fmt.Fprintf(tw, "%-36s\t", "Domain Name:")
		fmt.Fprintf(tw, "%-14s\t", "Domain Type:")
		fmt.Fprintf(tw, "%-36s\t", "Zone Name:")
		fmt.Fprintf(tw, "%-10s\t", "Zone Type:")
		fmt.Fprintf(tw, "%-10s\t", "Read Only:")
		fmt.Fprintf(tw, "%-13s\t", "Zone Enabled:")
		fmt.Fprintf(tw, "\n")
		for _, zoneRecord := range zoneRecords {
			fmt.Fprintf(tw, "%-36s\t", zoneRecord.DomainName)
			fmt.Fprintf(tw, "%-14s\t", zoneRecord.DomainType)
			fmt.Fprintf(tw, "%-36s\t", zoneRecord.ZoneName)
			fmt.Fprintf(tw, "%-10s\t", zoneRecord.ZoneType)
			fmt.Fprintf(tw, "%-10t\t", zoneRecord.ReadOnly)
			fmt.Fprintf(tw, "%-13t\t", zoneRecord.ZoneEnabled)
			fmt.Fprintf(tw, "\n")
			_ = tw.Flush()

		}

		fmt.Printf("Success, DNSrecordsList()\n")
	},
}

var cmdRecordCreate = &cobra.Command{
	Use:   "record-create",
	Short: "create record",
	Args:  cobra.MinimumNArgs(5),
	Run: func(cmd *cobra.Command, args []string) {
		dsms, err := ListDsms(DsmId)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		var dnsRecord webapi.DNSRecord

		dnsRecord.ZoneName = args[0]
		dnsRecord.DomainName = args[0]
		dnsRecord.Record = args[1]
		dnsRecord.TTL = args[2]
		dnsRecord.Type = args[3]
		dnsRecord.Value = args[4]

		for _, dsm := range dsms {
			if err := dsm.Login(); err != nil {
				fmt.Printf("Failed to login to DSM: [%s]. err: %v\n", dsm.Ip, err)
				os.Exit(1)
			}
			err := dsm.RecordCreate(dnsRecord)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			dsm.Logout()
		}

		fmt.Printf("Success, RecordCreate()\n")
	},
}

var cmdRecordDelete = &cobra.Command{
	Use:   "record-delete",
	Short: "delete record",
	Args:  cobra.MinimumNArgs(5),
	Run: func(cmd *cobra.Command, args []string) {
		dsms, err := ListDsms(DsmId)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		var dnsRecord webapi.DNSRecord

		dnsRecord.ZoneName = args[0]
		dnsRecord.DomainName = args[0]
		dnsRecord.Record = args[1]
		dnsRecord.TTL = args[2]
		dnsRecord.Type = args[3]
		dnsRecord.Value = args[4]

		for _, dsm := range dsms {
			if err := dsm.Login(); err != nil {
				fmt.Printf("Failed to login to DSM: [%s]. err: %v\n", dsm.Ip, err)
				os.Exit(1)
			}
			err := dsm.RecordDelete(dnsRecord)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			dsm.Logout()
		}

		fmt.Printf("Success, RecordDelete()\n")
	},
}

var cmdRecordFind = &cobra.Command{
	Use:   "record-find",
	Short: "find record",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		dsms, err := ListDsms(DsmId)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		var dnsRecord webapi.DNSRecord

		dnsRecord.ZoneName = args[0]
		dnsRecord.DomainName = args[0]
		dnsRecord.Record = args[1]

		for _, dsm := range dsms {
			if err := dsm.Login(); err != nil {
				fmt.Printf("Failed to login to DSM: [%s]. err: %v\n", dsm.Ip, err)
				os.Exit(1)
			}
			record, err := dsm.RecordFind(dnsRecord, "master")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			dsm.Logout()
			for _, r := range record {
				fmt.Printf("%v\n", r)
			}
		}

		fmt.Printf("Success, RecordFind()\n")
	},
}

func init() {
	cmdDNS.AddCommand(cmdDNSList)
	cmdDNS.AddCommand(cmdZoneList)
	cmdDNS.AddCommand(cmdRecordCreate)
	cmdDNS.AddCommand(cmdRecordDelete)
	cmdDNS.AddCommand(cmdRecordFind)
}
