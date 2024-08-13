/*
Copyright © 2024 Arnav Surve arnav@surve.dev
*/
package main

import (
	"taskman/cmd"
)

func main() {
	defer cmd.DB.Close()
	cmd.Execute()
}
