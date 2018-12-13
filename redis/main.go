package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"redis/argscheck"
	"redis/client"
	"redis/serv"
	"redis/types"
)

const (
	protocolTCP = "tcp"
	separator   = " "
	buff        = 512
)

var (
	gPort       = ":9090"
	gMemoryMode = "..." //expected disk or blank
	gIP         = "127.0.0.1"
)

func main() {
	fmt.Println("Start")
	argscheck.Start(os.Args, gPort, gMemoryMode, gIP)
	fmt.Println("ServFunc")
	ServFunc(gPort)
}

//ServFunc - containing Server and Client
func ServFunc(gPort string) {
	fmt.Printf("Port->%s\n", gPort)
	li, err := net.Listen(protocolTCP, gPort)
	fmt.Println("Listen")
	if err != nil {
		fmt.Println("Error: ", err)
		log.Fatal(err)
	}
	defer li.Close()
	fmt.Println("Entered with :" + gIP + gPort)
	go client.Commands(gPort, gIP, buff, protocolTCP)
	for {
		conn, err := li.Accept()
		fmt.Println("li.Accept")
		if err != nil {
			fmt.Println("Error9: ", err)
			log.Fatal(err)
		}
		defer conn.Close()
		ServConnCh := make(chan types.Server)
		fmt.Println("ServConnHandler start")
		go serv.ServConnHandler(ServConnCh, conn)
		fmt.Println("ServCmndsHandler start")
		go serv.ServCmndsHandler(ServConnCh, gMemoryMode)
		fmt.Println("end ServCmndsHandler")
	}
}
