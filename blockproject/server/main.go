package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/fatih/color"
)

func init() {
	color.Green("==============")
	color.Red("====启动区块链节点=====")
	color.Green("==============")

	log.SetPrefix("Blockchain: ")
}

func main() {

	port := flag.Uint("port", 9000, "TCP Port Number for Blockchain Server")
	flag.Parse()
	fmt.Printf("port::%v \n", *port)
	app := NewBlockchainServer(uint16(*port))
	app.Run()

}
