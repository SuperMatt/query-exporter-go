package server

import (
	"net"
	"time"
)

func CheckPort(address string, port string) (bool, error) {

	// check if the port is open
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(address, port), 5*time.Second)
	if err != nil {
		return false, err
	}
	conn.Close()
	return true, nil
}
