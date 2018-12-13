package serv

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"redis/save"
	"redis/types"
	"strings"
)

//ServConnHandler - scanning input connection and send to ServCmndsHandler
func ServConnHandler(ServConnCh chan types.Server, conn net.Conn) {
	scnnr := bufio.NewScanner(conn)
	for scnnr.Scan() {
		line := scnnr.Text()
		fmt.Printf("scanner ln: %s\n", line)
		inptFlds := strings.Fields(line)
		fmt.Println("inptFlds", line)
		rslt := make(chan string)
		ServConnCh <- types.Server{
			HandFlds: inptFlds,
			Rslt:     rslt,
		}
		fmt.Println("ServConnCh", ServConnCh)
		fmt.Println("13:conn <-rslt")
		io.WriteString(conn, <-rslt)
	}
}

//ServCmndsHandler - containing GET, SET, DEL commands.
func ServCmndsHandler(ServConnCh chan types.Server, gMemoryMode string) {
	fmt.Println("memData start")
	var memData = make(map[string]string)
	for cmnd := range ServConnCh {
		if len(cmnd.HandFlds) < 2 {
			cmnd.Rslt <- "At least 2 arguments for start\n"
			continue
		}
		switch cmnd.HandFlds[0] {
		//GET <key>
		case "GET":
			fmt.Println("case GET")
			if len(cmnd.HandFlds) != 2 {
				cmnd.Rslt <- "Get what?"
			}
			key := cmnd.HandFlds[1]
			value := memData[key]
			if len(memData) == 0 {
				fmt.Println("case GET == 0", memData)
				cmnd.Rslt <- "Map is empty"
			} else {
				fmt.Println("case GET", memData)
				cmnd.Rslt <- value
			}
			//SET <key>
		case "SET":
			fmt.Println("case SET")
			if len(cmnd.HandFlds) != 3 {
				fmt.Println("case SET error Missing value")
				cmnd.Rslt <- "Missing value\n"
			}
			key := cmnd.HandFlds[1]
			value := cmnd.HandFlds[2]
			memData[key] = value
			fmt.Println("case SET", memData)
			if gMemoryMode == "disk" {
				memDisk, err := json.Marshal(memData)
				if err != nil {
					fmt.Println("JSON", string(memDisk), err)
				}
				info := string(memDisk)
				save.SaveOnDisk(info)
				cmnd.Rslt <- "JSON: KEY - VALUE SET\n"
			} else {
				cmnd.Rslt <- "KEY - VALUE SET\n"
			}
			//DEL <KEY>
		case "DEL":
			fmt.Println("case DEL")
			key := cmnd.HandFlds[1]
			delete(memData, key)
			fmt.Println("case DET", memData)
			cmnd.Rslt <- "KEY - VALUE DELETED\n"
		default:
			fmt.Println("case noCommand default")
			cmnd.Rslt <- "I don't know this command :" + cmnd.HandFlds[0] + "\n"
		}
	}
}
