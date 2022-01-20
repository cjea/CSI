package livery

import (
	"fmt"
	"net"
	"time"
)

const (
	MAX_DATAGRAM_SIZE_BYTES = 1 << 10
	WAIT_TTL                = 1 * time.Second
	NUM_EMPTY_READS_ALLOWED = 5
)

var (
	// HELLO_MSG = []byte{0xff}
	// HELLO_RECV_MSG = []byte{0x66}
	HELLO_MSG      = []byte("HELLO\n")
	HELLO_RECV_MSG = []byte("HELLO RECEIVED")
)

func StartConversation(addr *net.UDPAddr) *Sender {
	conn, err := net.DialUDP(addr.Network(), nil, addr)
	sendCh := make(chan []byte)
	recvCh := make(chan []byte)

	if err != nil {
		return &Sender{Destination: addr, Err: err}
	}

	return &Sender{
		Conn:        conn,
		Destination: addr,
		SendCh:      sendCh,
		RecvCh:      recvCh,
	}
}

type Sender struct {
	Conn        *net.UDPConn
	Destination *net.UDPAddr
	Status      Status
	Err         error
	SendCh      chan []byte
	RecvCh      chan []byte
	NumFailures int
	Message     string
	// TODO: Include timeouts or consecutive failures for each status.
}

func (s *Sender) String() string {
	return fmt.Sprintf(
		"\n  {\n   Status: %d,\n   Err: %v\n  }\n",
		s.Status, s.Err,
	)
}

func (s *Sender) Error(err error) {
	if s.Err != nil {
		return
	}
	s.Err = err
}

// TODO: Append fencing tokens and store token <-> payload,
//			 and highest token for fast lookup.
func (s *Sender) Send(bs []byte) {
	_, _, err := s.Conn.WriteMsgUDP(bs, nil, nil)
	// _, err := s.Conn.WriteToUDP(bs, s.Destination)
	s.Error(err)
}

// Read has multiple possible next stages
func (s *Sender) Read(buf []byte) {
	emptyReads := 0
	fmt.Printf("Client is beginning a read\n")

	s.Conn.SetReadDeadline(time.Now().Add(WAIT_TTL))
	n, _, err := s.Conn.ReadFrom(buf)
	fmt.Printf("Client completed a read\n")

	for n == 0 && emptyReads < NUM_EMPTY_READS_ALLOWED {
		if err != nil {
			fmt.Printf("Errored with '%s', but continuing.\n", err.Error())
			continue
		}
		fmt.Printf("Failed read #%d (%d bytes)", emptyReads, n)
		emptyReads += 1
		time.Sleep(1 * time.Second)
		n, _, err = s.Conn.ReadFrom(buf)
	}
	fmt.Printf("Client stopped reading: %#v\n", err)
	s.Error(err)
}

func (s Sender) Close() error {
	return s.Conn.Close()
}

type Action struct {
	Op   Op
	Body []byte
}

func (s *Sender) SendHello() {
	fmt.Printf("Sending %s\n", HELLO_MSG)
	s.Send(HELLO_MSG)
}

func isHelloRecv(bs []byte) bool {
	for i := range bs {
		if bs[i] != HELLO_RECV_MSG[i] {
			return false
		}
	}
	return true
}

// TODO: Allow overriding specific statuses.
func (s *Sender) RunOp(action Action) *Sender {
	fmt.Printf("Running op: %v\n", action)

	if s.Err != nil {
		fmt.Printf("Already errored. Cant' run action %v\n", action)
		return s
	}
	op := action.Op
	assert(NUM_OPS == 7, "unexpected operations")
	switch op {
	case SEND_HELLO:
		s.SendHello()
		s.Status = WAITING_HELLO_RECV
	case RECV_HELLO:
		buf := make([]byte, len(HELLO_RECV_MSG))
		fmt.Printf("About to read\n")
		s.Read(buf)
		fmt.Printf("Read into buffer: '%s'\n", string(buf))
		if isHelloRecv(buf) {
			s.Status = READY_TO_SEND
		} else {
			s.Error(fmt.Errorf(
				"Expected %s, but got %s\n", string(HELLO_RECV_MSG), string(buf)))
		}
	case SEND_DATA:
		s.Send(action.Body)
		s.Status = WAITING_SEND_RECV
	case RECV_DATA_CONFIRMATION:
		// TODO: RecvData should handle retries.
		buf := []byte{}
		s.Read(buf)
		s.Status = READY_TO_SEND
	case TIMEOUT:
		s.Error(fmt.Errorf("TODO: TIMEOUT"))
	case SAY_BYE:
		s.Error(fmt.Errorf("TODO: SAY_BYE"))
	case RECV_BYE:
		s.Error(fmt.Errorf("TODO: RECV_BYE"))
	case NUM_OPS:
		s.Error(fmt.Errorf("TODO: NUM_OPS"))
	default:
		s.Error(fmt.Errorf("Unexpected operation: %v\n", op))
	}
	if s.Err != nil {
		s.Status = FAILED
	}
	return s
}

type Status = int

const (
	WAITING_HELLO_RECV Status = iota
	READY_TO_SEND
	WAITING_SEND_RECV
	READY_TO_CLOSE
	WAITING_BYE_RECV
	READY_TO_TEARDOWN
	FAILED
	NUM_STATUSES
)

type Op = int

const (
	SEND_HELLO Op = iota
	RECV_HELLO
	SEND_DATA
	RECV_DATA_CONFIRMATION
	TIMEOUT
	SAY_BYE
	RECV_BYE
	NUM_OPS
)

func assert(b bool, msg string) {
	if !b {
		panic(msg)
	}
}
