package server

import "net"

func GetIP(address string) (string, error) {
	if net.ParseIP(address) != nil {
		return address, nil
	}
	ips, err := net.LookupHost(address)
	if err != nil {
		return "", err
	}
	return ips[0], nil
}
