package main

import (
	"fmt"
	"net"
	"net/url"

	"github.com/Expost/GoDemo"
	"github.com/gorilla/websocket"
)

// func connHandler(conn net.Conn) {
// 	if conn == nil {
// 		return
// 	}

// 	buf := make([]byte, 4096)
// 	for {
// 		cnt, err := conn.Read(buf)
// 		if err != nil || cnt == 0 {
// 			conn.Close()
// 			break
// 		}

// 		fmt.Printf("Recv length is %d, and data is '%s'\n", cnt, string(buf))
// 	}
// }

// func proxy() {
// 	conn, err := net.Dial("tcp", "localhost:1090")
// 	if err != nil {
// 		fmt.Printf("failed to connect, %s\n", err)
// 		return
// 	}

// 	defer conn.Close()

// 	reader := bufio.NewReader(os.Stdin)
// 	buf := make([]byte, 1024)
// 	for {
// 		input, _ := reader.ReadString('\n')
// 		input = strings.TrimSpace(input)
// 		if input == "quit" {
// 			return
// 		}

// 		conn.Write([]byte(input))
// 		cnt, err := conn.Read(buf)
// 		if err != nil {
// 			fmt.Printf("failed to read data, %s\n", err)
// 			continue
// 		}

// 		fmt.Printf("read data is %s\n", string(buf[0:cnt]))
// 	}
// }

func localHandle(localConn net.Conn) {
	defer localConn.Close()

	u := url.URL{Scheme: "ws", Host: "localhost:12345", Path: "/ws"}
	dialer := websocket.Dialer{}
	remoteConn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		err := GoDemo.DecodeCopy(remoteConn, localConn)
		if err != nil {
			localConn.Close()
			remoteConn.Close()
		}
	}()

	GoDemo.EncodeCopy(localConn, remoteConn)
}

func main() {
	addr, err := net.ResolveTCPAddr("tcp", ":1090")
	if err != nil {
		fmt.Printf("resove tcp addr failed with %s\n", err)
		return
	}

	server, err := net.ListenTCP("tcp", addr)
	if err != nil {
		fmt.Printf("Failed to start server, %s\n", err)
		return
	}

	defer server.Close()
	fmt.Println("Server starting...")

	for {
		conn, err := server.AcceptTCP()
		if err != nil {
			fmt.Printf("Failed to connect, %s\n", err)
			break
		}

		conn.SetLinger(0)
		go localHandle(conn)
	}
}
