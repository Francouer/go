package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

const (
	PRTCL_TCP = "tcp"
	DNS       = "8.8.4.4:8888"

	EXIT       = "--exit"
	SHORT_EXIT = "-q"
	HOST       = ":connect"

	SEPARATE = " "
	BUFF     = 1024
)

var (
	g_connected   bool   = true
	g_port        string = ":9090"
	g_modememory  string = "..."
	g_def_ip_port string = "127.0.0.1"
	g_con_to_ip          = []string{g_def_ip_port + g_port}
	g_sz_conn     uint16 = 0
	g_path        string = "Info.txt"
)

type ServerFlds struct {
	HandFlds []string
	Rslt     chan string
}

func main() {
	fmt.Println("1:main start")
	inptArgsCheck(os.Args)
	fmt.Println("2:inptArgsCheck end/ServFunc start")
	ServFunc()
	fmt.Println("3:end of main")
}

//Функция проверяет входные аргументы и определяет порт
func inptArgsCheck(args []string) {
	var (
		host bool = false
		mode bool = false
		port bool = false
	)
	fmt.Println("4: 'for' args start")
	for _, value := range args[1:] {
		if port {
			fmt.Println("4.1: 'port' args start")
			g_port = ":" + value
			port = false
			continue
		}
		if mode {
			fmt.Println("4.2: 'mode' args start")
			g_modememory = value
			mode = false
			continue
		}
		if host {
			fmt.Println("4.3: 'host' args start")
			g_def_ip_port = value
			host = false
			continue
		}
		switch value {
		case "-p", "--port":
			fmt.Println("4.4: 'portcs' args start")
			port = true
		case "-m", "--mode":
			fmt.Println("4.5: 'modecs' args start")
			mode = true
		case "-h", "--host":
			fmt.Println("4.5: 'hostcs' args start")
			mode = true
		}
	}
	fmt.Println("5: 'for' args end")
}

// Функция сервера содержащая в себе и функционал клиента
func ServFunc() {
	fmt.Println("6: Launching server...")
	fmt.Printf("Port->%s\n", g_port)
	li, err := net.Listen(PRTCL_TCP, g_port)
	if err != nil {
		log.Fatal(err)
	}
	defer li.Close()
	defer fmt.Println("6: li.Close.")
	fmt.Println("7: Server is listenning.")

	fmt.Println("8: 'go' cliPos start")
	fmt.Println("8.1:Entered with :" + g_def_ip_port + g_port)
	go ClientFunc()

	for {
		conn, err := li.Accept()
		fmt.Println("9: conn := li.Accept()")
		if err != nil {
			log.Fatal(err)
		}
		ServConnCh := make(chan ServerFlds)
		defer conn.Close()
		defer fmt.Println("9: conn.Close.")
		go ConnInputHandler(ServConnCh, conn)
		fmt.Println("9.1: end ConnInputHandler")
		go ServCmndsHandler(ServConnCh)
		fmt.Println("9.2: end ServCmndsHandler")
		//buffer := make([]byte, BUFF)
		//content := ""
		//fmt.Println("9.3: strt conn.Read(buffer)")
		//for {
		//	fmt.Println("9.4: 'for' conn.Read(buffer)")
		//	length, err := conn.Read(buffer)
		//	if err != nil || length == 0 {
		//		break
		//	}
		//	content += string(buffer[:length])
		//}
		//fmt.Println("Print content: ", content)
	}
}

//Ловит пришедщие данные и отправляет на обработку
func ConnInputHandler(ServConnCh chan ServerFlds, conn net.Conn) {
	fmt.Println("10: ConnInputHandler start")
	defer conn.Close()
	//Определяем обьект
	scnnr := bufio.NewScanner(conn)
	fmt.Println("11: scnnr conn")
	for scnnr.Scan() {
		line := scnnr.Text()
		fmt.Printf("12: scanner ln: %s\n", line)
		inptFlds := strings.Fields(line)
		rslt := make(chan string)
		ServConnCh <- ServerFlds{
			HandFlds: inptFlds,
			Rslt:     rslt,
		}
		if g_modememory != "disk" {
			fmt.Println("13:conn <-rslt")
			io.WriteString(conn, <-rslt)
		} else {
			fmt.Println("14: storageFnc <-rslt")
			storageFnc(rslt)
		}
	}
}

//Обрабатывает команды GET SET DEL, держит данные в cash or disk
func ServCmndsHandler(ServConnCh chan ServerFlds) {
	fmt.Println("15: 'ServCmndsHandler' strt")
	var memData = make(map[string]string)
	for traf := range ServConnCh {
		if len(traf.HandFlds) < 2 {
			traf.Rslt <- "11:At least 2 arguments for start"
			continue
		}

		fmt.Println("15.1:GOT ServerFlds", traf)

		switch traf.HandFlds[0] {
		//GET <key>
		case "GET":
			fmt.Println("15.2 case GET")
			if len(traf.HandFlds) != 2 {
				traf.Rslt <- "Gimme key"
			}
			key := traf.HandFlds[1]
			value := memData[key]
			traf.Rslt <- value
			//SET <key>
		case "SET":
			fmt.Println("15.3: case SET")
			if len(traf.HandFlds) != 3 {
				traf.Rslt <- "Missing value\n"
			}
			key := traf.HandFlds[1]
			value := traf.HandFlds[2]
			memData[key] = value
			traf.Rslt <- " "
			//DEL <KEY>
		case "DEL":
			fmt.Println("15.4: case DEL")
			key := traf.HandFlds[1]
			delete(memData, key)
			traf.Rslt <- " "
		default:
			fmt.Println("15.5: case noCommand default")
			traf.Rslt <- "I don't know this command :" + traf.HandFlds[0] + "\n"
		}
	}
}

//Клиентская функция обрабатывающая команды из командной строки
func ClientFunc() {
	fmt.Println("16: 'go' ClientFunc start")
	scnnrCli := bufio.NewScanner(os.Stdin)
	fmt.Println("17: 'os.Stdin NewScanner' ClientFunc start")
	for scnnrCli.Scan() {
		lineC := scnnrCli.Text()
		fmt.Println("18:ClientFunc scnnrCli: ", lineC)
		CliFlds := strings.Fields(lineC)

		switch CliFlds[0] {

		case EXIT, SHORT_EXIT:
			fmt.Println("19:Exit client. See ya.")
			os.Exit(0)

		case HOST:
			if len(CliFlds) == 2 {
				fmt.Println("20:case HOST")
				var toggle bool = false
				for _, addr := range g_con_to_ip {
					if addr == CliFlds[1] {
						toggle = true
						break
					}
				}
				if toggle {
					conn, err := net.Dial(PRTCL_TCP, CliFlds[1])
					defer conn.Close()
					fmt.Println("21:!toggle case HOST")
					if err != nil {
						g_con_to_ip = append(g_con_to_ip, CliFlds[1])
						g_sz_conn++
						g_connected = true
						fmt.Println("22:conn.Write! !toggle")
						conn.Write([]byte(fmt.Sprintf("[User '%s:%s join']\n", g_def_ip_port, g_port)))
					}
				}
			}
		default:
			fmt.Println("23:default: g_connected")
			if g_connected {
				fmt.Println("24:default: g_connected=true")
				for indx, addr := range g_con_to_ip {
					conn, err := net.Dial(PRTCL_TCP, addr)
					fmt.Println("25:Connected to: " + addr)
					defer conn.Close()
					if err != nil {
						g_con_to_ip = remove(g_con_to_ip, indx)
						g_sz_conn--
						if g_sz_conn == 0 {
							g_connected = false
						}
					} else {
						fmt.Println("26: 'default' ClientFunc mssg")
						conn.Write([]byte(fmt.Sprintf("%s\n", lineC)))
					}
				}
			}
		}
	}
}

// функция записывающая информацию на диск
func storageFnc(rslt chan string) error {
	fmt.Println("27: 'storageFnc' strt")
	//Существующий файл с таким же именем будут перезаписан
	var fl, err = os.Create(g_path)
	if err != nil {
		panic(err)
	}
	defer fl.Close()
	info := <-rslt
	var byteWrttn, errWrt = fl.WriteString(info)
	if errWrt != nil {
		panic(errWrt)
	}
	fmt.Printf("28:Info.txt written: %v\n", byteWrttn)
	return nil
}

//remove func
func remove(list []string, num int) []string {
	fmt.Println("29: 'remove' done")
	return append(list[:num], list[num+1:]...)
}
