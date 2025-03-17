package p2p

import (
	"fmt"
	"net"
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


type TCPTransportOpts struct {
	ListenAddress    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
	OnPeer        func(Peer) error
}

type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener
	rpcch    chan RPC
}
func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcch:            make(chan RPC, 1024),
	}
}

func (transport *TCPTransport) ListenAndAccept() error{
	var err error
	transport.listener,err=net.Listen("tcp",transport.ListenAddress)
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
  var err error

  defer func() {
		fmt.Printf("dropping peer connection: %s", err)
		connection.Close()
	}()

  peer := NewTCPPeer(connection, true)
  if err= transport.HandshakeFunc(peer); err != nil {
		return
	}
  fmt.Printf("New connection established with connection %v\n", peer)
  
  //Read Loop
  for {
		rpc := RPC{}
		err = transport.Decoder.Decode(connection, &rpc)
		if err != nil {
			return
		}

		rpc.From = connection.RemoteAddr().String()
        fmt.Printf("Message received%s",rpc.Payload)
		
 
	}
}