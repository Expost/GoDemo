package GoDemo

import (
	"io"
	"net"
)

func DecodeCopy(src net.Conn, dst net.Conn) error {
	buf := make([]byte, 1024)
	for {
		readCount, errRead := src.Read(buf)
		if errRead != nil {
			if errRead != io.EOF {
				return errRead
			} else {
				return nil
			}
		}

		if readCount > 0 {
			writeCount, errWrite := dst.Write(buf[0:readCount])
			if errWrite != nil {
				return errWrite
			}

			if readCount != writeCount {
				return io.ErrShortWrite
			}
		}
	}
}

func EncodeCopy(src net.Conn, dst net.Conn) error {
	buf := make([]byte, 1024)
	for {
		readCount, errRead := src.Read(buf)
		if errRead != nil {
			if errRead != io.EOF {
				return errRead
			} else {
				return nil
			}
		}

		if readCount > 0 {
			writeCount, errWrite := dst.Write(buf[0:readCount])
			if errWrite != nil {
				return errWrite
			}

			if readCount != writeCount {
				return io.ErrShortWrite
			}
		}
	}
}
