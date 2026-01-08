/*
 * Copyright 2025 Mark DeNeve
 * based on findings from here: https://github.com/acmesh-official/acme.sh/issues/4696
 */

package webapi

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type DNSRecord struct {
	ZoneName   string `json:"zone_name"`   // required always the same as domain_name
	DomainName string `json:"domain_name"` // required always the same as zone_name
	Record     string `json:"rr_owner"`    // fully qualified domain name
	Type       string `json:"rr_type"`     // type of record A,TXT,CNAME
	TTL        string `json:"rr_ttl"`      // time to live for record
	Value      string `json:"rr_info"`     // value for the record
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

// RecordCreate creates a new DNS record
// dnsRecord.Record MUST end with a dot (.) or it will be prepended to the zone name by the Synology
func (dsm *DSM) RecordCreate(dnsRecord DNSRecord) error {
	params := url.Values{}
	params.Add("api", "SYNO.DNSServer.Zone.Record")
	params.Add("method", "create")
	params.Add("version", "1")
	params.Add("zone_name", strconv.Quote(dnsRecord.ZoneName))
	params.Add("domain_name", strconv.Quote(dnsRecord.DomainName))
	params.Add("rr_owner", strconv.Quote(dnsRecord.Record))
	params.Add("rr_ttl", strconv.Quote(dnsRecord.TTL))
	params.Add("rr_type", strconv.Quote(dnsRecord.Type))
	params.Add("rr_info", strconv.Quote(dnsRecord.Value))

	resp, err := dsm.sendRequest("", &struct{}{}, params, "webapi/entry.cgi")
	if err != nil {
		log.Error("There was an error.")
		log.Debugf("The Params were: %s", params.Encode())
		log.Errorf("The error was: %s, the response was %v\n", err, resp)
		return err
	}
	log.Infof("Record %s created successfully.", dnsRecord.Record)
	log.Debug(params.Encode())
	log.Debug(resp)

	return nil
}

// RecordDelete deletes a DNS record
// dnsRecord.Record MUST end with a dot (.) if it is fully qualified
func (dsm *DSM) RecordDelete(dnsRecord DNSRecord) error {
	// The Synology API for DNS record deletion expects the full record to be passed to it
	// the value of dnsRecord.fullRecord changes, but I am not sure why
	// in order to make sure that we find the right version, we will need to do a search of all existing records, and then compare
	// to the one we want to delete EXCLUDING the fullRecord field
	// we can then delete that record if we find a match

	existingRecords, err := dsm.RecordFind(dnsRecord, "master")
	log.Infof("Found existing record: %v\n", existingRecords)

	if err != nil {
		return err
	}

	deletedRecords := 0

	for _, record := range existingRecords {
		if compareRecords(record, dnsRecord) {
			params := url.Values{}
			params.Add("api", "SYNO.DNSServer.Zone.Record")
			params.Add("method", "delete")
			params.Add("version", "1")

			jsonData, err := json.Marshal(record)
			if err != nil {
				log.Error("Error marshaling record to JSON")
				return err
			}

			// The API expects an array of records to delete, even if it's just one.
			jsonData = []byte("[" + string(jsonData) + "]")
			params.Add("items", string(jsonData))

			resp, err := dsm.sendRequest("", &struct{}{}, params, "webapi/entry.cgi")
			if err != nil {
				log.Error("There was an error.")
				log.Debugf("The Params were: %s\n", params.Encode())
				log.Errorf("The error was: %s, the response was %v\n", err, resp)
				return err
			}
			deletedRecords++
			log.Infof("Record %s deleted successfully.", record.Record)
			log.Debug(params.Encode())
			log.Debug(resp)
		}
	}
	if deletedRecords == 0 {
		return fmt.Errorf("no matching records found to delete")
	}
	return nil
}

func compareRecords(r1, r2 DNSRecord) bool {
	return r1.ZoneName == r2.ZoneName &&
		r1.DomainName == r2.DomainName &&
		r1.Record == r2.Record &&
		r1.Type == r2.Type &&
		r1.TTL == r2.TTL &&
		r1.Value == r2.Value
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
// RecordList gets all records for the specified zone names and zone type (master/slave)
func (dsm *DSM) RecordList(zoneName []string, zoneType string) ([]DNSRecord, error) {

	type DNSRecords struct {
		Record []DNSRecord `json:"items"`
	}

	var allDNSRecords []DNSRecord

	for _, zn := range zoneName {
		params := url.Values{}
		params.Add("api", "SYNO.DNSServer.Zone.Record")
		params.Add("method", "list")
		params.Add("version", "1")
		params.Add("action", "enum")
		params.Add("sort_by", "rr_owner")
		params.Add("sort_direction", "ASC")
		params.Add("filter_by", "")
		params.Add("zone_name", zn)
		params.Add("domain_name", zn)
		params.Add("zone_type", zoneType)

		resp, err := dsm.sendRequest("", &DNSRecords{}, params, "webapi/entry.cgi")
		if err != nil {
			log.Error("There was an error.")
			log.Debugf("The Params were: %s", params.Encode())
			log.Errorf("The error was: %s, the response was %v\n", err, resp)
			return nil, errCodeMapping(resp.ErrorCode, err)
		}

		dNSRecords, ok := resp.Data.(*DNSRecords)
		if !ok {
			return nil, fmt.Errorf("failed to assert response to %T", &DNSRecords{})
		}
		allDNSRecords = append(allDNSRecords, dNSRecords.Record...)
	}

	return allDNSRecords, nil
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
		log.Error("There was an error.")
		log.Errorf("The error was: %s, the response was %v\n", err, resp)
		return nil, errCodeMapping(resp.ErrorCode, err)
	}
	log.Debug(resp)

	dNSRecords, ok := resp.Data.(*ZoneRecords)
	if !ok {
		return nil, fmt.Errorf("failed to assert response to %T", &ZoneRecords{})
	}
	return dNSRecords.Record, nil
}

func (dsm *DSM) RecordFind(dnsRecord DNSRecord, zoneType string) ([]DNSRecord, error) {

	type DNSRecords struct {
		Record []DNSRecord `json:"items"`
	}

	var allDNSRecords []DNSRecord

	params := url.Values{}
	params.Add("api", "SYNO.DNSServer.Zone.Record")
	params.Add("method", "list")
	params.Add("version", "1")
	params.Add("offset", "0")
	params.Add("limit", "50")
	params.Add("action", "find")
	params.Add("filterString", dnsRecord.Record)
	params.Add("sort_by", "rr_owner")
	params.Add("sort_direction", "ASC")
	params.Add("filter_by", dnsRecord.Record)
	params.Add("zone_name", dnsRecord.ZoneName)
	params.Add("domain_name", dnsRecord.DomainName)
	params.Add("zone_type", zoneType)

	resp, err := dsm.sendRequest("", &DNSRecords{}, params, "webapi/entry.cgi")
	if err != nil {
		log.Error("There was an error.")
		log.Debugf("The Params were: %s", params.Encode())
		log.Errorf("The error was: %s, the response was %v\n", err, resp)
		return nil, errCodeMapping(resp.ErrorCode, err)
	}

	dNSRecords, ok := resp.Data.(*DNSRecords)
	if !ok {
		return nil, fmt.Errorf("failed to assert response to %T", &DNSRecords{})
	}
	allDNSRecords = append(allDNSRecords, dNSRecords.Record...)

	// OK this is so dumb but the Synology API doesn't return the full record in the find or list methods
	// so we have to take the response, and fill in the rest of the missing details
	for i := range allDNSRecords {
		allDNSRecords[i].ZoneName = dnsRecord.ZoneName
		allDNSRecords[i].DomainName = dnsRecord.DomainName
	}
	return allDNSRecords, nil

}
