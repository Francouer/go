package argscheck

import "fmt"

//Start - update gPort, gMode, gIP variable if it's need to be changed
func Start(args []string, gPort string, gMemoryMode string, gIP string) (string, string, string) {
	var (
		host = false
		mode = false
		port = false
	)
	fmt.Println("'for' args start")
	for _, value := range args[1:] {
		if port {
			fmt.Println("'port' args start")
			gPort = ":" + value
			port = false
			continue
		}
		if mode {
			fmt.Println("'mode' args start")
			gMemoryMode = value
			mode = false
			continue
		}
		if host {
			fmt.Println("'connectTo' args start")
			gIP = value
			host = false
			continue
		}
		switch value {
		case "-p", "--port":
			fmt.Println("'portcs' args start")
			port = true
		case "-m", "--mode":
			fmt.Println("'modecs' args start")
			mode = true
		case "-h", "--host":
			fmt.Println("'connectTocs' args start")
			mode = true
		}
	}
	fmt.Println("'for' args end")
	return gPort, gIP, gMemoryMode
}
