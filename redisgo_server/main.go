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
	CHECK_FILE = "[:~:-file-:~:]"
	CHECK_NAME = "[:~:-name-:~:]"

	PRTCL_TCP = "tcp"
	DNS       = "8.8.4.4:8888"

	EXT     = "--exit"
	SRT_EXT = "-q"
	HST     = "--host"
	SRT_HST = "-h"

	SPRTR = "->"
	BUFF  = 1024
)

var (
	g_cnctd     bool   = false
	g_mn_prt    string = ":9090"
	g_mdplc     string = "..."
	g_strg      string = "disk"
	g_hst       string = "127.0.0.1"
	g_snd_to_ip        = []string{g_hst}
	g_sz_conn   uint16 = 0
	g_path      string = "/src/go/redisgolang/redisgo_server/Info.txt"
)

type Nero struct {
	Hits []string
	Rslt chan string
}

//Главная горутина
func main() {
	inpt_rgs(os.Args)
	srv_pos()
}

//Функция проверяет входные аргументы и определяет порт
func inpt_rgs(args []string) {
	var (
		host bool = false
		mode bool = false
		port bool = false
	)
	for _, value := range args[1:] {
		if port {
			g_mn_prt = ":" + value
			port = false
			continue
		}
		if mode {
			g_mdplc = value
			mode = false
			continue
		}
		if host {
			g_hst = value
			host = false
			continue
		}
		switch value {
		case "-p", "--port":
			port = true
		case "-m", "--mode":
			mode = true
		case "-h", "--host":
			mode = true
		}
	}
}

// Функция сервера содержащая в себе и функционал клиента
func srv_pos() {
	var (
		buffer  []byte = make([]byte, BUFF)
		splited []string
		content string = ""
	)
	go cliPos()
	fmt.Println("Launching server...")
	fmt.Printf("Port->%s\n", g_mn_prt)
	li, err := net.Listen(PRTCL_TCP, g_mn_prt)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Server is listenning.")
	Combos := make(chan Nero)
	go srvrRds_go(Combos)
	for {
		conn, err := li.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go hndlFnc(Combos, conn)
	}
}

//Оснавная нагрузка описана здесь
func srvrRds_go(Combos chan Nero) {
	var dict = make(map[string]string)
	for cmbs := range Combos {
		if len(cmbs.Hits) < 2 {
			cmbs.Rslt <- "At least 2 arguments for start"
			continue
		}

		//fmt.Println("GOT Nero", cmbs)

		switch strings.ToUpper(cmbs.Hits[0]) {
		//GET <key>
		case "GET":
			if len(cmbs.Hits) != 2 {
				cmbs.Rslt <- "Gimme key"
			}
			key := cmbs.Hits[1]
			value := dict[key]
			cmbs.Rslt <- value
			//SET <key>
		case "SET":
			if len(cmbs.Hits) != 3 {
				cmbs.Rslt <- "Missing value"
			}
			key := cmbs.Hits[1]
			value := cmbs.Hits[2]
			dict[key] = value
			cmbs.Rslt <- " "
			//DEL <KEY>
		case "DEL":
			key := cmbs.Hits[1]
			delete(dict, key)
			cmbs.Rslt <- " "
		default:
			cmbs.Rslt <- "I don't know this command :" + cmbs.Hits[0] + "\n"
		}
	}
}

func hndlFnc(Combos chan Nero, conn net.Conn) {
	defer conn.Close()
	//Определяем обьект
	scnnr := bufio.NewScanner(conn)
	for scnnr.Scan() {
		ln := scnnr.Text()
		fmt.Printf("scanner ln: %s", ln)
		fs := strings.Fields(ln)
		rslt := make(chan string)
		Combos <- Nero{
			Hits: fs,
			Rslt: rslt,
		}
		if g_mdplc != g_strg {
			io.WriteString(conn, <-rslt)
		} else {
			storageFnc(rslt)
		}
	}
}

//client func and clisrv
func cliPos() {
	fmt.Println("Entered with :" + g_hst)
	go redisCli()
}

//redisClient_go
func redisCli() {
	//cmmnd, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	//mssg := strings.Replace(cmmnd, "\n", "", -1)
	//cmd := strings.Split(strings.Replace(mssg, " ", "", -1), SPRTR)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		scnln := scanner.Text()
		fmt.Println(scnln)
		scnfs := strings.Fields(scnln)

		switch scnfs[0] {

		case EXT, SRT_EXT:
			fmt.Println("Exit client. Good day.")
			os.Exit(0)

		case SRT_HST, HST:
			if len(scnfs) == 2 {
				var flg bool = false
				for _, vle := range g_snd_to_ip {
					if vle == scnfs[1] {
						flg = true
						break
					}
				}
				if !flg {
					conn, err := net.Dial(PRTCL_TCP, scnfs[1])
					defer conn.Close()
					if err != nil {
						g_snd_to_ip = append(g_snd_to_ip, scnfs[1])
						g_sz_conn++
						g_cnctd = true
						conn.Write([]byte(fmt.Sprintf("[User '%s:%s join']\n", g_hst, g_mn_prt)))
					}
				}
			}
		default:
			if g_cnctd {
				for indx, vle := range g_snd_to_ip {
					conn, err := net.Dial(PRTCL_TCP, vle)
					fmt.Println("Connected to: " + vle)
					defer conn.Close()
					if err != nil {
						g_snd_to_ip = remove(g_snd_to_ip, indx)
						g_sz_conn--
						if g_sz_conn == 0 {
							g_cnctd = false
						}
					} else {
						conn.Write([]byte(fmt.Sprintf("[%s:%s]: %s\n", g_hst, g_mn_prt, scnfs)))
					}
				}
			}
		}
	}
}

// функция записывающая информацию на диск
func storageFnc(rslt chan string) error {
	//Существующий файл с таким же именем будут перезаписан
	var fl, err = os.Create(g_path)
	if err != nil {
		panic(err)
	}
	defer fl.Close()
	var byteWrttn, errWrt = fl.WriteString(<-rslt)
	if errWrt != nil {
		panic(errWrt)
	}
	fmt.Printf("Info.txt written: %v\n", byteWrttn)
	return nil
}

//remove func
func remove(list []string, num int) []string {
	return append(list[:num], list[num+1:]...)
}
