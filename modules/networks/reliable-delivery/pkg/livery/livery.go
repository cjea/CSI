package livery

import (
	"fmt"
	"net"
	"time"
)

const (
	MAX_DATAGRAM_SIZE_BYTES = 1 << 10
	WAIT_TTL                = 10000 * time.Second
)

var (
	// HELLO_MSG = []byte{0xff}
	// HELLO_RECV_MSG = []byte{0x66}
	HELLO_MSG      = []byte("HELLO")
	HELLO_RECV_MSG = []byte("HELLO RECEIVED")
)

func StartConversation(hostname string) (*Sender, error) {
	conn, err := net.Dial("udp", hostname)
	if err != nil {
		return nil, err
	}
	sendCh := make(chan []byte)
	recvCh := make(chan []byte)
	return &Sender{
		Conn:   conn,
		Status: WAITING_HELLO_RECV,
		SendCh: sendCh,
		RecvCh: recvCh,
	}, nil
}

type Sender struct {
	Conn        net.Conn
	Status      LiveryStatus
	Err         error
	SendCh      chan []byte
	RecvCh      chan []byte
	NumFailures int
	Message     string
	// TODO: Include timeouts or consecutive failures for each status.
}

// TODO: Append fencing tokens and store token <-> payload,
//			 and highest token for fast lookup.
func (s Sender) Send(bs []byte) error {
	_, err := s.Conn.Write(bs)
	return err
}

func (s Sender) RecvInto(buf []byte) {
	// buf := make([]byte, MAX_DATAGRAM_SIZE_BYTES)
	s.Conn.SetReadDeadline(time.Now().Add(WAIT_TTL))
	fmt.Printf("RecvInto with buf %v\n", buf)
	n, err := s.Conn.Read(buf)
	fmt.Printf("Num bytes read: %d\n", n)
	s.Err = err
}

func (s Sender) Close() error {
	return s.Conn.Close()
}

type Action struct {
	Op   LiveryOperation
	Body []byte
}

func (s *Sender) ErrIllegalStatusChange(action Action) {
	err := fmt.Errorf(
		"cannot perform %v while having status %v\n",
		s.Status, action.Op,
	)
	s.Err = err
}

func (s *Sender) SendHello() {
	s.Err = s.Send(HELLO_MSG)
}

func isHelloRecv(bs []byte) bool {
	fmt.Printf("Running isHelloRecv with %v\n", bs)

	return string(bs) == string(HELLO_RECV_MSG)
}

// TODO: Allow overriding specific statuses.
func (s *Sender) RunOp(action Action) *Sender {
	op := action.Op
	assert(NUM_LIVERY_OPS == 7, "unexpected operations")
	switch op {
	case SEND_HELLO:
		s.SendHello()
		s.Status = WAITING_HELLO_RECV
	case RECV_HELLO:
		buf := make([]byte, len(HELLO_RECV_MSG))
		s.RecvInto(buf)
		if isHelloRecv(buf) {
			s.Status = READY_TO_SEND
		}
	case SEND_DATA:
		s.Send(action.Body)
		s.Status = WAITING_SEND_RECV
	case RECV_DATA_CONFIRMATION:
		// TODO: RecvData should handle retries.
		// s.RecvData()
		s.Err = fmt.Errorf("TODO")
	case TIMEOUT:
		s.Err = fmt.Errorf("TODO")
	case SAY_BYE:
		s.Err = fmt.Errorf("TODO")
	case RECV_BYE:
		s.Err = fmt.Errorf("TODO")
	case NUM_LIVERY_OPS:
		s.Err = fmt.Errorf("TODO")
	default:
		s.Err = fmt.Errorf("Unexpected operation: %v\n", op)
	}
	if s.Err != nil {
		s.Status = FAILED
	}
	return s
}

type LiveryStatus = int

const (
	WAITING_HELLO_RECV LiveryStatus = iota
	READY_TO_SEND
	WAITING_SEND_RECV
	READY_TO_CLOSE
	WAITING_BYE_RECV
	READY_TO_TEARDOWN
	FAILED
	NUM_LIVERY_STATUSES
)

type LiveryOperation = int

const (
	SEND_HELLO LiveryOperation = iota
	RECV_HELLO
	SEND_DATA
	RECV_DATA_CONFIRMATION
	TIMEOUT
	SAY_BYE
	RECV_BYE
	NUM_LIVERY_OPS
)

func assert(b bool, msg string) {
	if !b {
		panic(msg)
	}
}
