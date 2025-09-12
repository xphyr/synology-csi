/*
 * Copyright 2025 Mark DeNeve
 * based on findings from here: https://github.com/acmesh-official/acme.sh/issues/4696
 */

package webapi

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type DNSRecord struct {
	ZoneName   string `json:"zone_name"`   // required always the same as domain_name
	DomainName string `json:"domain_name"` // required always the same as zone_name
	RROwner    string `json:"rr_owner"`    // fully qualified domain name
	RRttl      string `json:"rr_ttl"`      // time to live for record
	RRtype     string `json:"rr_type"`     // type of record A,TXT,CNAME
	RRInfo     string `json:"rr_info"`     // value for the record
	FullRecord string `json:"full_record"`
}

type ZoneRecord struct {
	DomainName  string `json:"domain_name"`
	DomainType  string `json:"domain_type"`
	ReadOnly    bool   `json:"is_readonly"`
	ZoneEnabled bool   `json:"zone_enable"`
	ZoneName    string `json:"zone_name"`
	ZoneType    string `json:"zone_type"`
}

func (dsm *DSM) RecordCreate(dnsRecord DNSRecord) {
	params := url.Values{}
	params.Add("api", "SYNO.DNSServer.Zone.Record")
	params.Add("method", "create")
	params.Add("version", "1")
	params.Add("zone_name", dnsRecord.ZoneName)
	params.Add("domain_name", dnsRecord.DomainName)
	params.Add("rr_owner", dnsRecord.RROwner)
	params.Add("rr_ttl", dnsRecord.RRttl)
	params.Add("rr_type", dnsRecord.RRtype)
	params.Add("rr_info", dnsRecord.RRInfo)

	resp, err := dsm.sendRequest("", &struct{}{}, params, "webapi/entry.cgi")
	if err != nil {
		fmt.Println("There was an error.")
		fmt.Printf("The error was: %s, the response was %s\n", err, resp)
		//return nil, errCodeMapping(resp.ErrorCode, err)
	}
	fmt.Println(resp)

}

func (dsm *DSM) RecordDelete(dnsRecord DNSRecord) {
	params := url.Values{}
	params.Add("api", "SYNO.DNSServer.Zone.Record")
	params.Add("method", "delete")
	params.Add("version", "1")

	jsonData, err := json.Marshal(dnsRecord)
	if err != nil {
		fmt.Println("error")
	}
	params.Add("items", string(jsonData))

	resp, err := dsm.sendRequest("", &struct{}{}, params, "webapi/entry.cgi")
	if err != nil {
		fmt.Println("There was an error.")
		fmt.Printf("The error was: %s, the response was %s\n", err, resp)
		//return nil, errCodeMapping(resp.ErrorCode, err)
	}
	fmt.Println(resp)
}

/*
func (dsm *DSM) updateDNSrecord() {
	params := url.Values{}
	params.Add("api", "SYNO.DNSServer.Zone.Record")
	params.Add("method", "update")
	params.Add("version", "1")

	resp, err := dsm.sendRequest("", &info, params, "webapi/entry.cgi")
}

*/

func (dsm *DSM) RecordList(zoneName string, zoneType string) ([]DNSRecord, error) {
	params := url.Values{}
	params.Add("api", "SYNO.DNSServer.Zone.Record")
	params.Add("method", "list")
	params.Add("version", "1")
	params.Add("action", "enum")
	params.Add("sort_by", "rr_owner")
	params.Add("sort_direction", "ASC")
	params.Add("filter_by", "")
	params.Add("zone_name", zoneName)
	params.Add("domain_name", zoneName)
	params.Add("zone_type", zoneType)

	type DNSRecords struct {
		Record []DNSRecord `json:"items"`
	}

	resp, err := dsm.sendRequest("", &DNSRecords{}, params, "webapi/entry.cgi")
	if err != nil {
		fmt.Println("There was an error.")
		fmt.Printf("The error was: %s, the response was %s\n", err, resp)
		return nil, errCodeMapping(resp.ErrorCode, err)
	}
	fmt.Println(resp)

	dNSRecords, ok := resp.Data.(*DNSRecords)
	if !ok {
		return nil, fmt.Errorf("failed to assert response to %T", &DNSRecords{})
	}
	return dNSRecords.Record, nil
}

func (dsm *DSM) ZoneList() ([]ZoneRecord, error) {
	params := url.Values{}
	params.Add("api", "SYNO.DNSServer.Zone")
	params.Add("method", "list")
	params.Add("version", "1")

	type ZoneRecords struct {
		Record []ZoneRecord `json:"items"`
	}

	resp, err := dsm.sendRequest("", &ZoneRecords{}, params, "webapi/entry.cgi")
	if err != nil {
		fmt.Println("There was an error.")
		fmt.Printf("The error was: %s, the response was %s\n", err, resp)
		return nil, errCodeMapping(resp.ErrorCode, err)
	}
	fmt.Println(resp)

	dNSRecords, ok := resp.Data.(*ZoneRecords)
	if !ok {
		return nil, fmt.Errorf("failed to assert response to %T", &ZoneRecords{})
	}
	return dNSRecords.Record, nil
}
