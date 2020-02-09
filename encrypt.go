package GoDemo

import (
	"io"
	"net"

	"github.com/gorilla/websocket"
)

func DecodeCopy(src *websocket.Conn, dst net.Conn) error {
	for {
		_, buf, errRead := src.ReadMessage()
		readCount := len(buf)
		if errRead != nil {
			if errRead != io.EOF {
				return errRead
			}

			return nil
		}

		if readCount > 0 {
			_, errWrite := dst.Write(buf[0:readCount])
			if errWrite != nil {
				return errWrite
			}
		}
	}
}

func EncodeCopy(src net.Conn, dst *websocket.Conn) error {
	buf := make([]byte, 256)
	for {
		readCount, errRead := src.Read(buf)
		if errRead != nil {
			if errRead != io.EOF {
				return errRead
			}

			return nil
		}

		if readCount > 0 {
			errWrite := dst.WriteMessage(websocket.BinaryMessage, buf[0:readCount])
			if errWrite != nil {
				return errWrite
			}
		}
	}
}
