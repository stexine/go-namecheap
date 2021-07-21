package namecheap

import "encoding/xml"

type Config struct {
	User   string `json:"user"`
	ApiKey string `json:"api_key"`
	IP     string `json:"ip"`
}

type Host struct {
	XMLName            xml.Name `xml:"host"`
	HostId             string   `xml:"HostId"`
	Name               string   `xml:"Name"`
	Type               string   `xml:"Type"`
	Address            string   `xml:"Address"`
	MXPref             string   `xml:"MXPref"`
	TTL                int      `xml:"TTL"`
	AssociatedAppTitle string   `xml:"AssociatedAppTitle"`
	FriendlyName       string   `xml:"FriendlyName"`
	IsActive           bool     `xml:"IsActive"`
	IsDDNSEnabled      bool     `xml:"IsDDNSEnabled"`
}

type Hosts struct {
	XMLName       xml.Name `xml:"hosts"`
	Domain        string   `xml:"Domain,attr"`
	EmailType     string   `xml:"EmailType,attr"`
	IsUsingOurDNS bool     `xml:"IsUsingOurDNS,attr"`
	Hosts         []Host   `xml:"Host"`
}
