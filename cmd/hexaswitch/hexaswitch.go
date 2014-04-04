package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/morriswinkler/hexabus"
	"os"
)

var opts struct {
	Version   bool   `long:"version" description:"print libhexabus version and exit"`
	Command   string `short:"c" long:"command" description:"{get|set|epquery|send|listen|on|off|status|power|devinfo}"`
	Ip        string `short:"i" long:"ip" description:"the hostname to connect to"`
	Bind      string `short:"b" long:"bind" description:"local IP address to use"`
	Interface string `short:"I" long:"interface" description:"for listen: interface to listen on. otherwise: outgoing interface for multicast"`
	Eid       uint32 `short:"e" long:"eid" description:"Endpoint ID (EID)"`
	Dtype     uint   `shor:"d" long:"datatype" description:"{1: Bool | 2: UInt8 | 3: UInt32 | 4: HexaTime | 5:Float | 6: String}"`
	Value     string `short:"v" long:"value" description:"Value"`
	Oneline   bool   `long:"oneline" description:"Print each receive packet on one line"`
}

func main() {

	flags.Parse(&opts)

	if opts.Version {
		fmt.Println("go hexabus library: " + hexabus.VERSION)
		os.Exit(0)
	}

}
