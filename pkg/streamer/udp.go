package streamer

import "net"

type UDPWriter struct {
	net.Conn
	address string
}

func NewUDPWriter(address string) (UDPWriter, error) {
	conn, err := net.Dial("udp", address)
	if err != nil {
		return UDPWriter{}, nil
	}

	return UDPWriter{
		Conn:    conn,
		address: address,
	}, nil
}
