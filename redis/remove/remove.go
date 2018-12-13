package remove

import "fmt"

//Remove - remove element from slice
func Remove(list []string, num int) []string {
	fmt.Println("29: 'remove' done")
	return append(list[:num], list[num+1:]...)
}
