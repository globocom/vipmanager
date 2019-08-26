package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/globocom/vipmanager/model"
	"github.com/urfave/cli"
)

var app = cli.NewApp()
var napiURLQA = "http://networkapi.qa01.globoi.com"
var napiURLPROD = "http://networkapi.globoi.com"
var napiURL string
var napiUser string
var napiPass string

type Operation int

const (
	ADD Operation = 0
	REM Operation = 1
)

func main() {
	info()
	flags()
	commands()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func info() {
	app.Name = "Vip manager"
	app.Usage = "A Simple CLI for managing Vips"
	app.Author = "Evolução Infra - Team"
	app.Version = "0.1"
}

func flags() {
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "env, e",
			Value: "qa",
			Usage: "NetworkApi env prod/qa",
		},
		cli.StringFlag{
			Name:   "ldap_user",
			Usage:  "User from ladp to auth on the napi (without domain)",
			EnvVar: "ldap_user",
		},
		cli.StringFlag{
			Name:   "ldap_pass",
			Usage:  "Pass from ladp to auth on the napi",
			EnvVar: "ldap_pass",
		},
	}
}

func commands() {
	app.Commands = []cli.Command{
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "Add machines to the Vip pools",
			Action:  func(c *cli.Context) { updateVip(c, ADD) },
			Flags: []cli.Flag{
				cli.StringFlag{Name: "vip, v"},
				cli.StringFlag{Name: "real, r"},
			},
		},
		{
			Name:    "rem",
			Aliases: []string{"r"},
			Usage:   "Removing machines from Vip pools",
			Action:  func(c *cli.Context) { updateVip(c, REM) },
			Flags: []cli.Flag{
				cli.StringFlag{Name: "vip, v"},
				cli.StringFlag{Name: "real, r"},
			},
		},
	}
}

func updateVip(c *cli.Context, op Operation) {
	vipName := c.String("vip")
	ip := c.String("real")
	if c.GlobalString("env") == "prod" {
		napiURL = napiURLPROD
	} else {
		napiURL = napiURLQA
	}

	napiUser = c.GlobalString("ldap_user")
	napiPass = c.GlobalString("ldap_pass")

	if napiUser == "" || napiPass == "" {
		log.Fatalln("$ldap_user and $ldap_pass need to me set")
	}

	fmt.Printf("%s, %s\n", vipName, ip)
	ipID, err := getIPID(ip)
	if err != nil {
		log.Fatalf("Error retrieving ip: %s\n %v", ip, err)
	}

	log.Printf("The id of IP: %s is: %d", ip, ipID)

	vip, err := getVip(vipName)
	if err != nil {
		log.Fatalf("Error retrieving vip: %s\n %s", vipName, err.Error())
	}

	serverPoolIDs := []int{}
	for _, port := range vip.Ports {
		for _, pool := range port.Pools {
			serverPoolIDs = append(serverPoolIDs, pool.ServerPoolID)
		}
	}

	sPools, err := getServerPools(serverPoolIDs)
	if err != nil {
		log.Fatalf("Error retrieving the pools: %v\n %s", serverPoolIDs, err.Error())
	}

	if op == ADD {
		sPools.AddMember(ipID, ip)
	} else {
		sPools.RemMember(ipID)
	}

	err = storeServerPools(serverPoolIDs, sPools)

	if err != nil {
		log.Fatalf("Error updating the server pools: %v for the vip: %s\n %v", serverPoolIDs, vipName, err.Error())
	}
	log.Println("Operation concluded successfully")
}

func storeServerPools(sPoolIds []int, sPools model.ServerPools) error {
	req := model.ServerPoolPutRequest{IDs: sPoolIds}
	resp := &model.ServerPools{}
	sPoolsp := &sPools
	pURL, err := req.Build(napiURL)
	if err != nil {
		return err
	}

	err = callHTTP("PUT", pURL, sPoolsp, resp)
	if err != nil {
		return err
	}

	return nil
}

func getServerPools(serverPoolIds []int) (model.ServerPools, error) {
	req := model.ServerPools{Ids: serverPoolIds}

	pURL, err := req.Build(napiURL)
	if err != nil {
		return model.ServerPools{}, err
	}

	var sPools model.ServerPools
	err = callHTTP("GET", pURL, nil, &sPools)
	if err != nil {
		return model.ServerPools{}, err
	}
	if len(sPools.ServerPools) == 0 {
		return model.ServerPools{}, errors.New("No server pools were found")
	}

	if len(sPools.ServerPools) != len(serverPoolIds) {
		return model.ServerPools{}, fmt.Errorf(
			"the number of ServerPoolIds: %d is differente from the returned in the response: %d",
			len(serverPoolIds),
			len(sPools.ServerPools),
		)
	}

	log.Printf("%v", sPools.ServerPools[0].Members[0].IPV6)

	return sPools, nil
}

func getVip(name string) (model.Vip, error) {
	req := model.VipRequest{
		ExtendsSearch: []map[string]string{
			{"name": name},
		},
	}
	pURL, err := req.Build(napiURL)
	if err != nil {
		return model.Vip{}, err
	}

	var vips model.Vips
	err = callHTTP("GET", pURL, nil, &vips)
	if err != nil {
		return model.Vip{}, err
	}
	if len(vips.Vips) == 0 {
		return model.Vip{}, errors.New("Vip not found")

	}

	return vips.Vips[0], nil
}

func callHTTP(method string, url string, body interface{}, resp interface{}) error {
	var mBody []byte
	var err error
	if body != nil {
		mBody, err = json.Marshal(body)

		fmt.Println(string(mBody))
		if err != nil {
			return err
		}
	}

	rq, err := http.NewRequest(method, url, bytes.NewBuffer(mBody))
	if err != nil {
		return err
	}
	rq.Header.Set("Content-Type", "application/json")
	rq.SetBasicAuth(napiUser, napiPass)

	client := &http.Client{}
	rs, err := client.Do(rq)
	if rs.StatusCode != 200 && rs.StatusCode != 201 {
		s, _ := ioutil.ReadAll(rs.Body)
		log.Printf("Status code error: %d \n %s\n", rs.StatusCode, string(s))
		return fmt.Errorf("Status code: %d", rs.StatusCode)
	}
	if err != nil {
		return err
	}

	dec := json.NewDecoder(rs.Body)
	return dec.Decode(resp)
}

func buildIPReq(ip string) ([]byte, error) {
	i := net.ParseIP(ip)
	if i == nil {
		return nil, errors.New("Error parsing IP")
	}
	octs := strings.Split(ip, ".")
	oct1, _ := strconv.Atoi(octs[0])
	oct2, _ := strconv.Atoi(octs[1])
	oct3, _ := strconv.Atoi(octs[2])
	oct4, _ := strconv.Atoi(octs[3])
	req := model.Ipv4Request{
		ExtendsSearch: []model.Ipv4ExtendSearch{
			model.Ipv4ExtendSearch{
				Oct1: oct1,
				Oct2: oct2,
				Oct3: oct3,
				Oct4: oct4,
			},
		},
	}
	return json.Marshal(req)
}

func getIPID(ip string) (int, error) {
	j, err := buildIPReq(ip)
	if err != nil {
		return 0, err
	}
	params := url.Values{}
	params.Add("search", string(j))
	url, _ := url.Parse(napiURL)
	url.Path += "api/v3/ipv4"
	url.RawQuery = params.Encode()

	var res model.Ipv4Resp
	err = callHTTP("GET", url.String(), nil, &res)
	if err != nil {
		return 0, err
	}

	if len(res.Ips) > 0 {
		return res.Ips[0].Id, nil
	}
	return 0, errors.New("Ip not found")
}
