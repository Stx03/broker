package main

import (
	"flag"
	"fmt"
	"net"
	"net/rpc"
	"uk.ac.bris.cs/distributed3/pairbroker/stubs"
)

type Factory struct{}

//TODO: Define a Multiply function to be accessed via RPC.
//Check the previous weeks' examples to figure out how to do this.

func (f *Factory) Multiply(req stubs.Pair, res *stubs.JobReport) (err error) {
	x := req.X
	y := req.Y
	res.Result = x * y
	return
}

func main() {
	pAddr := flag.String("ip", "127.0.0.1:8050", "IP and port to listen on")
	brokerAddr := flag.String("broker", "127.0.0.1:8030", "Address of broker instance")
	flag.Parse()
	//TODO: You'll need to set up the RPC server, and subscribe to the running broker instance.
	rpc.Register(&Factory{})
	client, _ := rpc.Dial("tcp", *brokerAddr)
	defer client.Close()
	req := stubs.Subscription{Topic: "multiply", FactoryAddress: *pAddr, Callback: "Factory.Multiply"}
	res := new(stubs.StatusReport)
	listen, err := net.Listen("tcp", *pAddr)
	if err != nil {
		fmt.Println(err)
	}
	client.Call(stubs.Subscribe, req, res)
	defer listen.Close()
	rpc.Accept(listen)
}
