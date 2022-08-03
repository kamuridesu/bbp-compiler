package main

import "fmt"

func input(message string) string {
	fmt.Print(message)
	var usr_input string
	fmt.Scanln(&usr_input)
	return usr_input
}
