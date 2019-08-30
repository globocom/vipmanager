package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
)

//Ipv4Request Ipv4 request representation from the Network API
type Ipv4Request struct {
	IP            string             `json:"-"`
	ExtendsSearch []Ipv4ExtendSearch `json:"extends_search"`
}

//Ipv4ExtendSearch Extended search representations from the Network API
type Ipv4ExtendSearch struct {
	Oct1 int `json:"oct1"`
	Oct2 int `json:"oct2"`
	Oct3 int `json:"oct3"`
	Oct4 int `json:"oct4"`
}

//Ipv4Resp Just a simple representation of the Ipv4 Response from the Network API
type Ipv4Resp struct {
	Ips []Ipv4ResIP `json:"ips"`
}

//Ipv4ResIP Ids from Ips returned in the query
type Ipv4ResIP struct {
	ID int `json:"id"`
}

//Build a ipv4 query url
func (ipReq *Ipv4Request) Build() (string, error) {
	i := net.ParseIP(ipReq.IP)
	if i == nil {
		return "", errors.New("Error parsing IP")
	}
	octs := strings.Split(ipReq.IP, ".")
	oct1, _ := strconv.Atoi(octs[0])
	oct2, _ := strconv.Atoi(octs[1])
	oct3, _ := strconv.Atoi(octs[2])
	oct4, _ := strconv.Atoi(octs[3])

	ipReq.ExtendsSearch = []Ipv4ExtendSearch{
		Ipv4ExtendSearch{
			Oct1: oct1,
			Oct2: oct2,
			Oct3: oct3,
			Oct4: oct4,
		},
	}
	bs, err := json.Marshal(ipReq)
	if err != nil {
		return "", err
	}

	params := url.Values{}
	params.Add("search", string(bs))
	url := url.URL{}
	url.Path += "ipv4"
	url.RawQuery = params.Encode()
	return url.String(), nil
}

//GET one Ip from Napi. If more the one ip or none was returned
//return an error
func (ipReq *Ipv4Request) GET() (int, error) {
	urlSuffix, err := ipReq.Build()

	var res Ipv4Resp
	err = HTTP.Call("GET", urlSuffix, nil, &res)
	if err != nil {
		return 0, err
	}

	if len(res.Ips) == 1 {
		printIP(res.Ips[0].ID, ipReq.IP)
		return res.Ips[0].ID, nil
	} else if len(res.Ips) > 1 {
		return 0, errors.New("More than one Ip was returned")
	}

	return 0, errors.New("Ip not found")
}

func printIP(ID int, IP string) {
	fmt.Println()
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.TabIndent)
	fmt.Fprintf(w, "ID\tIPV4\t\n")
	fmt.Fprintf(w, "%d\t%s\t\n\n", ID, IP)
	w.Flush()
}
