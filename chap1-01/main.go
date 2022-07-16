package main

import "fmt"

func main() {
	stringlist := [5]string{"I", "am", "stupid", "and", "weak"}
	for i, value := range stringlist {
		switch value {
		case "stupid":
			stringlist[i] = "smart"
		case "weak":
			stringlist[i] = "strong"
		}
	}
	for _, value := range stringlist {
		fmt.Println(value)
	}
}
