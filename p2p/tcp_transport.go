package p2p

import (
	"fmt"
	"net"
	"sync"
)

// TCPPeer represents the remote node over a TCP established connection.
type TCPPeer struct {
	// The underlying connection of the peer. Which in this case
	// is a TCP connection.
	connection net.Conn
	// if we dial and retrieve a conn => outbound == true
	// if we accept and retrieve a conn => outbound == false
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		connection:     conn,
		outbound: outbound,
	}
}

type TCPTransport struct {
	listenAddress string
	listener    net.Listener
	mu sync.RWMutex
	peers map[net.Addr]Peer
}

func NewTCPTransport(listenAddr string) *TCPTransport {
	return &TCPTransport{
		listenAddress: listenAddr,
		
	}
}

func (transport *TCPTransport) ListenAndAccept() error{
	var err error
	transport.listener,err=net.Listen("tcp",transport.listenAddress)
	if err != nil{
      return err
	}
	go transport.startAcceptLoop();
	return nil;
}

func (transport *TCPTransport) startAcceptLoop(){
	for{ 
	connection,err:=transport.listener.Accept()
	 if err != nil{
		fmt.Printf("Error during accepting TCP connection: %s\n", err)
	}
	go transport.handleConnection(connection);
	}
}

func (transport *TCPTransport) handleConnection(connection net.Conn){
  peer := NewTCPPeer(connection, true)
  
  fmt.Printf("New connection established with connection %v\n", peer)
}