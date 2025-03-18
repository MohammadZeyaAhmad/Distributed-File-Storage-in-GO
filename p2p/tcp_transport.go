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
	wg *sync.WaitGroup
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		connection:     conn,
		outbound: outbound,
		wg:       &sync.WaitGroup{},
	}
}

func (p *TCPPeer) CloseStream() {
	p.wg.Done()
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
	rpcChannel    chan RPC
}
func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcChannel:            make(chan RPC, 1024),
	}
}

// Consume implements the Tranport interface, which will return read-only channel
// for reading the incoming messages received from another peer in the network.
func (transport *TCPTransport) Consume() <-chan RPC {
	return transport.rpcChannel
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

  if transport.OnPeer != nil {
		if err = transport.OnPeer(peer); err != nil {
			return
		}
	}
  
  //Read Loop
  for {
		rpc := RPC{}
		err = transport.Decoder.Decode(connection, &rpc)
		if err != nil {
			return
		}

		rpc.From = connection.RemoteAddr().String()

		if rpc.Stream {
			peer.wg.Add(1)
			fmt.Printf("[%s] incoming stream, waiting...\n", connection.RemoteAddr())
			peer.wg.Wait()
			fmt.Printf("[%s] stream closed, resuming read loop\n", connection.RemoteAddr())
			continue
		}

		transport.rpcChannel <- rpc
        fmt.Printf("Message received%s",rpc.Payload)
		
 
	}
}