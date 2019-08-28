package main

import (
	"log"
	"os"

	"github.com/globocom/vipmanager/http"
	"github.com/globocom/vipmanager/model"
	"github.com/urfave/cli"
)

var app = cli.NewApp()

//Operation Kind of pool operation
type Operation int

const (
	//ADD a real to a pool
	ADD Operation = 0

	//REM a real from a pool
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

func updateVip(c *cli.Context, op Operation) {
	vipName := c.String("vip")
	ip := c.String("real")
	env := c.GlobalString("env")
	napiUser := c.GlobalString("ldap_user")
	napiPass := c.GlobalString("ldap_pass")

	if napiUser == "" || napiPass == "" {
		log.Fatalln("$ldap_user and $ldap_pass need to me set")
	}

	model.HTTP = http.New(env, napiUser, napiPass)

	ipReq := model.Ipv4Request{IP: ip}
	ipID, err := ipReq.GET()
	if err != nil {
		log.Fatalf("Error retrieving ip: %s\n %v", ip, err)
	}
	log.Printf("The id of IP: %s is: %d", ip, ipID)

	vipReq := model.VipRequest{}
	vip, err := vipReq.GET(vipName)
	if err != nil {
		log.Fatalf("Error retrieving vip: %s\n %s", vipName, err.Error())
	}

	spReq := model.ServerPoolReq{}
	err = spReq.GET(vip)
	if err != nil {
		log.Fatalf("Error retrieving the pools from : %s\n %s", vipName, err.Error())
	}

	if op == ADD {
		err = spReq.AddMember(ipID, ip)
	} else {
		err = spReq.RemMember(ipID)
	}

	if err != nil {
		log.Fatalf("Error updating the server pools: %v for the vip: %s\n %v", spReq.Ids, vipName, err.Error())
	}

	log.Println("Operation concluded successfully")
}
