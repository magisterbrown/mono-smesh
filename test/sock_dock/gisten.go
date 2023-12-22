package main

import (
    "net"
    "log"
)

func main() {
    socket, err := net.Listen("unix", "./mysocket.sock")
    if err != nil {
        panic(err)
    }
    defer socket.Close()

    for {
        conn, err := socket.Accept()
        if err != nil {
            panic(err)
        }
        go func(conn net.Conn) {
            defer conn.Close()
            // Create a buffer for incoming data.
            buf := make([]byte, 4096)

            // Read data from the connection.
            n, err := conn.Read(buf)
            if err != nil {
                log.Fatal(err)
            }

            // Echo the data back to the connection.
            _, err = conn.Write(buf[:n])
            if err != nil {
                log.Fatal(err)
            }
        }(conn) 

    }

}
