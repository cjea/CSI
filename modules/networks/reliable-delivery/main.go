package main

import (
	"fmt"
	"livery/pkg/livery"
	"net"
	"time"
)

/*
'livery protocol

Goal:
=====
  * [X] peer to peer
  * [ ] reliable delivery of stream of packets
  * [ ] method of detecting corruption
  * [X] method of ending the conversation

Problems:
=========
  * corrupted data
  * dropped packets

STATES:
=======

Alice: --SEND HELLO-->  | wait_hello_recv |
       --RECV HELLO-->  | ready_to_send | ,___
          --SEND-->  | wait_send_recv |      |
          --RECV-->  | ready_to_send | -------
                -->  | ready_to_close |
          --SAY BYE--> | wait_bye_recv |
                    --RECV BYE--> | teardown |
                    --TIMEOUT-->  | teardown |
          --TEARDOWN--

Bob: --LISTEN-->  | wait_hello |
       --RECV HELLO-->  | ready_to_confirm_hello
       --SEND HELLO-->  | wait_data |   <-=========
          --RECV-->  | ready_to_confirm_data |     }
              --CONFIRM DATA--> | wait_data | ======
                -->  | ready_to_close |
          --SAY BYE--> | wait_bye_recv |
                    --RECV BYE--> | teardown |
                    --TIMEOUT-->  | teardown |
          --TEARDOWN--

FLOW:
=====

Alice: HELLO      (TTL 10)
Bob:   RECV HELLO (TTL 10)

Alice: DATA A
Bob:   RECV A
Alice: DATA B

---- DROP ----

Bob:   RECV A
Alice: DATA B
Bob:   RECV B
Alice: DATA C

---- CORRUPT DATAGRAM ----

Bob:   RECV B
Alice: DATA C
Bob:   RECV C

Alice: GOODBYE (TTL 60)
BOB:   GOODBYE (TTL 60)
*/

const serverHostname = "127.0.0.1:53865"

func main() {
	go func() {
		fmt.Printf("Listening on %s\n", serverHostname)
		l, err := net.ListenUDP("udp", &net.UDPAddr{
			IP:   net.IP{127, 0, 0, 1},
			Port: 53865,
		})
		must(err)
		buf := []byte{}
		for {
			n, _, err := l.ReadFrom(buf)
			if n > 0 {
				fmt.Printf("Listener received message '%s'\n", string(buf))
				continue
			}

			fmt.Printf("Received %d bytes\n", n)
			must(err)
			time.Sleep(1 * time.Second)
		}
	}()

	actions := []livery.Action{
		{Op: livery.SEND_HELLO, Body: nil},
		{Op: livery.RECV_HELLO, Body: nil},
		{Op: livery.SEND_DATA, Body: []byte("the really secret message")},
		{Op: livery.RECV_DATA_CONFIRMATION, Body: nil},
	}

	c := livery.StartConversation(
		&net.UDPAddr{IP: net.IP{127, 0, 0, 1}, Port: 53865},
	)

	for i := 0; i < len(actions) && c.Err == nil; i++ {
		c.RunOp(actions[i])
	}

	fmt.Printf("Conversation: %s\n", c.String())
	c.Close()
	// fmt.Printf("Number of statuses %+v\n", livery.NUM_LIVERY_STATUSES)
	// fmt.Printf("Number of actions %+v\n", livery.NUM_LIVERY_OPS)
}

func main2() {
	buf := make([]byte, 1<<10)
	host := "127.0.0.1:8090"
	fmt.Printf("Listening on %s\n", host)
	c, err := net.Dial("udp", host)
	fmt.Printf("Listening on %s (%s)\n", host, c.LocalAddr().String())
	must(err)
	// _, err = c.Read(buf)
	packet, err := net.ListenPacket("udp", c.LocalAddr().String())
	must(err)
	packet.ReadFrom(buf)
	fmt.Printf("Buf: %v\n", buf)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
