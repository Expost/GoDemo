package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"net/http"

	"github.com/Expost/GoDemo"
	"github.com/gorilla/websocket"
)

// func remoteHandle(localConn net.Conn) {
// 	defer localConn.Close()

// 	buf := make([]byte, 256)
// 	_, readErr := localConn.Read(buf)
// 	if readErr != nil || buf[0] != 0x05 {
// 		return
// 	}

// 	localConn.Write([]byte{0x05, 0x00})

// 	n, readErr := localConn.Read(buf)
// 	if readErr != nil || n < 7 {
// 		return
// 	}

// 	if buf[1] != 0x1 {
// 		return
// 	}

// 	var dIP []byte
// 	switch buf[3] {
// 	case 0x01:
// 		dIP = buf[4 : 4+net.IPv4len]
// 		fmt.Printf("0x01 -> ip:%s\n", string(dIP))
// 	case 0x03:
// 		ipAddr, err := net.ResolveIPAddr("ip", string(buf[5:n-2]))
// 		if err != nil {
// 			return
// 		}

// 		dIP = ipAddr.IP
// 		fmt.Printf("0x03 -> ip:%d-%d.%d.%d.%d(%s)\n", len(dIP), dIP[0], dIP[1], dIP[2], dIP[3], buf[5:n-2])
// 	case 0x04:
// 		dIP = buf[4 : 4+net.IPv6len]
// 		fmt.Printf("0x04 -> ip:%s\n", string(dIP))
// 	default:
// 		return
// 	}

// 	dPort := buf[n-2:]
// 	dstAddr := &net.TCPAddr{
// 		IP:   dIP,
// 		Port: int(binary.BigEndian.Uint16(dPort)),
// 	}

// 	fmt.Printf("Port:%d\n", dstAddr.Port)

// 	remoteConn, err := net.DialTCP("tcp", nil, dstAddr)
// 	if err != nil {
// 		return
// 	}

// 	defer remoteConn.Close()
// 	localConn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

// 	go func() {
// 		err := GoDemo.DecodeCopy(localConn, remoteConn)
// 		if err != nil {
// 			localConn.Close()
// 			remoteConn.Close()
// 		}
// 	}()

// 	GoDemo.EncodeCopy(remoteConn, localConn)
// }

// func use() {
// 	addr, err := net.ResolveTCPAddr("tcp", ":1091")
// 	if err != nil {
// 		fmt.Printf("resove tcp addr failed with %s\n", err)
// 		return
// 	}

// 	server, err := net.ListenTCP("tcp", addr)
// 	if err != nil {
// 		fmt.Printf("Failed to start server, %s\n", err)
// 		return
// 	}

// 	defer server.Close()
// 	fmt.Println("Server starting...")

// 	for {
// 		conn, err := server.AcceptTCP()
// 		if err != nil {
// 			fmt.Printf("Failed to connect, %s\n", err)
// 			break
// 		}

// 		conn.SetLinger(0)
// 		go remoteHandle(conn)
// 	}
// }

func webRemoteHandle(localConn *websocket.Conn) {
	defer localConn.Close()

	//buf := make([]byte, 256)
	_, buf, readErr := localConn.ReadMessage()
	if readErr != nil || buf[0] != 0x05 {
		return
	}

	localConn.WriteMessage(websocket.BinaryMessage, []byte{0x05, 0x00})

	_, buf, readErr = localConn.ReadMessage()
	n := len(buf)
	if readErr != nil || n < 7 {
		return
	}

	if buf[1] != 0x1 {
		return
	}

	var dIP []byte
	switch buf[3] {
	case 0x01:
		dIP = buf[4 : 4+net.IPv4len]
		fmt.Printf("0x01 -> ip:%s\n", string(dIP))
	case 0x03:
		ipAddr, err := net.ResolveIPAddr("ip", string(buf[5:n-2]))
		if err != nil {
			return
		}

		dIP = ipAddr.IP
		fmt.Printf("0x03 -> ip:%d-%d.%d.%d.%d(%s)\n", len(dIP), dIP[0], dIP[1], dIP[2], dIP[3], buf[5:n-2])
	case 0x04:
		dIP = buf[4 : 4+net.IPv6len]
		fmt.Printf("0x04 -> ip:%s\n", string(dIP))
	default:
		return
	}

	dPort := buf[n-2:]
	dstAddr := &net.TCPAddr{
		IP:   dIP,
		Port: int(binary.BigEndian.Uint16(dPort)),
	}

	fmt.Printf("Port:%d\n", dstAddr.Port)

	remoteConn, err := net.DialTCP("tcp", nil, dstAddr)
	if err != nil {
		return
	}

	defer remoteConn.Close()
	localConn.WriteMessage(websocket.BinaryMessage, []byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	go func() {
		err := GoDemo.DecodeCopy(localConn, remoteConn)
		if err != nil {
			localConn.Close()
			remoteConn.Close()
		}
	}()

	GoDemo.EncodeCopy(remoteConn, localConn)
}

func wsPage(res http.ResponseWriter, req *http.Request) {
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
	if err != nil {
		http.NotFound(res, req)
		return
	}

	go webRemoteHandle(conn)
}

func main() {
	http.HandleFunc("/ws", wsPage)
	http.ListenAndServe(":12345", nil)
}
