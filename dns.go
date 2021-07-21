package namecheap

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/beevik/etree"
)

func GetHosts(SLD string, TLD string, conf Config) ([]Host, error) {

	url := "https://api.namecheap.com/xml.response?apiuser=" + conf.User + "&apikey=" + conf.ApiKey + "&username=" + conf.User + "&Command=namecheap.domains.dns.getHosts&ClientIp=" + conf.IP + "&SLD=" + SLD + "&TLD=" + TLD

	result, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}

	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(body); err != nil {
		return nil, err

	} else {
		hosts := doc.FindElements("//ApiResponse/CommandResponse/DomainDNSGetHostsResult/host")
		resultHosts := []Host{}
		tmpHost := Host{}

		for _, host := range hosts {
			ttl, _ := strconv.Atoi(host.SelectAttrValue("TTL", "1800"))
			active, _ := strconv.ParseBool(host.SelectAttrValue("IsActive", "true"))
			ddns, _ := strconv.ParseBool(host.SelectAttrValue("IsDDNSEnabled", "true"))

			tmpHost = Host{
				XMLName:            xml.Name{},
				HostId:             host.SelectAttrValue("HostId", ""),
				Name:               host.SelectAttrValue("Name", ""),
				Type:               host.SelectAttrValue("Type", ""),
				Address:            host.SelectAttrValue("Address", ""),
				MXPref:             host.SelectAttrValue("MXPref", "10"),
				TTL:                ttl,
				AssociatedAppTitle: host.SelectAttrValue("AssociatedAppTitle", ""),
				FriendlyName:       host.SelectAttrValue("FriendlyName", ""),
				IsActive:           active,
				IsDDNSEnabled:      ddns,
			}

			resultHosts = append(resultHosts, tmpHost)

		}
		return resultHosts, nil
	}
}

func SetHosts(SLD string, TLD string, conf Config, hosts []Host) error {
	if len(hosts) <= 0 {
		return nil
	} else {
		updateUrl := "https://api.namecheap.com/xml.response?apiuser=" + conf.User + "&apikey=" + conf.ApiKey + "&username=" + conf.User + "&Command=namecheap.domains.dns.setHosts&ClientIp=" + conf.IP + "&SLD=" + SLD + "&TLD=" + TLD

		formData := url.Values{}

		count := 1
		countStr := ""

		// check every host
		for _, host := range hosts {

			countStr = strconv.Itoa(count)

			formData.Add("HostName"+countStr, host.Name)
			formData.Add("RecordType"+countStr, host.Type)
			formData.Add("Address"+countStr, host.Address)
			formData.Add("MXPref"+countStr, host.MXPref)
			formData.Add("TTL"+countStr, strconv.Itoa(host.TTL))
			formData.Add("EmailType"+countStr, "MX")
			formData.Add("IsActive"+countStr, strconv.FormatBool(host.IsActive))
			formData.Add("IsDDNSEnabled"+countStr, strconv.FormatBool(host.IsDDNSEnabled))

			count++
		}

		postResp, _ := http.PostForm(updateUrl, formData)
		//We Read the response body on the line below.
		body, _ := ioutil.ReadAll(postResp.Body)
		fmt.Println(string(body))
	}

	return nil
}

func AddHost(SLD string, TLD string, conf Config, newHost Host) error {

	hosts, err := GetHosts(SLD, TLD, conf)

	if err != nil {
		return err
	}

	found := false
	for i, h := range hosts {
		if newHost.Name == h.Name && newHost.Type == h.Type {
			found = true
			hosts[i] = newHost
			break
		}
	}

	if found == false {
		hosts = append(hosts, newHost)
	}

	return SetHosts(SLD, TLD, conf, hosts)
}

func DelHost(SLD string, TLD string, conf Config, name string, recordType string) error {
	hosts, err := GetHosts(SLD, TLD, conf)

	if err != nil {
		return err
	}

	for i, h := range hosts {
		if name == h.Name && recordType == h.Type {

			hosts[i] = hosts[len(hosts)-1]
			hosts = hosts[:len(hosts)-1]
			break
		}
	}

	return SetHosts(SLD, TLD, conf, hosts)
}
