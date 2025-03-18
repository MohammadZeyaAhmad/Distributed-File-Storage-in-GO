package main

import (
	"fmt"
	"log"

	"github.com/MohammadZeyaAhmad/DFS/p2p"
)

func main() {
	tcptransportOpts := p2p.TCPTransportOpts{
		ListenAddress:    ":3000",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
	}
	tcpTransport := p2p.NewTCPTransport(tcptransportOpts)
     
    go func () {
		for {
			msg:=tcpTransport.Consume();
			fmt.Printf("%v\n", msg);
		}
	}()
	if err:=tcpTransport.ListenAndAccept();err!=nil {
		log.Fatal(err)
	}
	select{}
}