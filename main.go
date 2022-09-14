/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/davidemaggi/kog/cmd/kog"
)

var (
	Version string = "0.0.0"
)

func main() {
	kog.Execute(Version)
}
