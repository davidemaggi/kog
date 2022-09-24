/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"github.com/davidemaggi/kog/cmd/kog"
	"os"
	"os/signal"
	"syscall"
)

var (
	Version string = "0.0.0"
)

func main() {

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		os.Exit(1)
	}()

	kog.Execute(Version)
}

func cleanup() {
	fmt.Println("cleanup")
}
