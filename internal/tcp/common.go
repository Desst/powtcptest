package tcp

import (
	"encoding/json"
	"fmt"
	"net"
	v1 "powtcptest/internal/protocol/v1"
)

func SendMessage(msg v1.Message, conn net.Conn) error {
	if err := json.NewEncoder(conn).Encode(msg); err != nil {
		return fmt.Errorf("unable to encode message: %w", err)
	}

	return nil
}
