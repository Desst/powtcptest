package server

import (
	"net"
)

type connInfo struct {
	conn                net.Conn
	challenge           string
	challengeDifficulty int
}

func newConnInfo(conn net.Conn, challenge string, difficulty int) *connInfo {
	return &connInfo{
		conn:                conn,
		challenge:           challenge,
		challengeDifficulty: difficulty,
	}
}
