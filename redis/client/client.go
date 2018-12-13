package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

//Commands - read stdin, input connection, write to connection
func Commands(gPort string, gIP string, buff int, protocolTCP string) (Dconn net.Conn) {
	var (
		SizeConn  uint16 = 1
		buffer           = make([]byte, buff)
		IPHandler        = []string{gIP + gPort}
	)
	fmt.Println("Dial")
	Dconn, err := net.Dial(protocolTCP, gPort)
	if nil != err {
		log.Fatal(err)
	}
	for {
		fmt.Println("NewReader")
		mssg, err := bufio.NewReader(os.Stdin).ReadString('\n')
		fmt.Println("NewReader", mssg)
		if nil != err {
			fmt.Printf("Error - %s", err)
		}
		Splitmssg := strings.Fields(mssg)
		fmt.Println("Splitmssg", Splitmssg)

		switch Splitmssg[0] {
		case "--exit", "-q":
			fmt.Println("exit Commands. See ya.")
			os.Exit(0)

		case "--connect":
			if len(Splitmssg) == 2 {
				fmt.Println("case connectTo")
				var toggle = false
				for _, addr := range IPHandler {
					if addr == Splitmssg[1] {
						toggle = true
						break
					}
				}
				if !toggle {
					D2conn, err := net.Dial(protocolTCP, Splitmssg[1])
					defer D2conn.Close()
					fmt.Println("21:!toggle case connectTo")
					if err == nil {
						IPHandler = append(IPHandler, Splitmssg[1])
						SizeConn++
						fmt.Println("22:conn2.Write! !toggle")
						D2conn.Write([]byte(fmt.Sprintf("User '%s:%s join'\n", gIP, gPort)))
					}
				}
			}
		default:
			//write to socket
			fmt.Println("Dconn.Write")
			Dconn.Write([]byte(mssg))

			//read from socket
			fmt.Println("Dconn.Write")
			var content string
			read, err := Dconn.Read(buffer)
			fmt.Println("Dconn.Read", read)
			if err != nil || read == 0 {
				fmt.Println("read == 0", content)
				log.Println(string(buffer[:read]))
			}
			fmt.Println("content", content)
			log.Println(string(buffer[:read]))
		}
	}
}
