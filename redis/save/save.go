package save

import (
	"fmt"
	"os"
)

//SaveOnDisk - saves data into file on disk
func SaveOnDisk(info string) error {
	Path := "Info.txt"
	fmt.Println("27: 'SaveOnDisk' strt")
	//Существующий файл с таким же именем будут перезаписан
	var fl, err = os.Create(Path)
	if err != nil {
		panic(err)
	}
	defer fl.Close()
	var byteWrttn, errWrt = fl.WriteString(info)
	if errWrt != nil {
		panic(errWrt)
	}
	fmt.Printf("28:Info.txt written: %v\n", byteWrttn)
	return nil
}
