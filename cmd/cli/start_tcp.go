package main

import (
	"bufio"
	"errors"
	"github.com/isichei/recipe-book/internal/filesyncer"
	"log"
	"net"
)

func startMainTCP(addr, apiKey, directory string, pingOnly, tls bool) {
	conn, err := filesyncer.CreateMainSenderConn(addr, apiKey, tls)
	if err != nil {
		log.Fatalf("Can't chat on addr: %s. Error: %s\n", addr, err)
	}

	if pingOnly {
		log.Println("Connection complete. Ping only so closing connection.")
	} else {
		fc, err := filesyncer.CreateRawMdFileCache(directory)
		if err != nil {
			log.Fatalf("Failed to create the file cache: %s\n", err)
		}
		mainSyncer := filesyncer.Syncer{Replica: false, Conn: conn, FileCache: fc}
		err = mainSyncer.Run()
		if err != nil {
			log.Fatalf("Main syncer failed: %s\n", err)
		}
		log.Println("Main syncer completed.")
	}
}

func handleTCP(l net.Listener, port int) {
	for {
		conn, err := l.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				log.Printf("TCP listener on port %d closed", port)
				return
			}

			log.Printf("error accepting on %d/tcp: %s", port, err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(c net.Conn) {
	defer c.Close()

	log.Printf("New connection from %s", c.RemoteAddr())

	lines := bufio.NewReader(c)

	for {
		line, err := lines.ReadString('\n')
		if err != nil {
			log.Printf("Connection from %s closed: %v", c.RemoteAddr(), err)
			return
		}

		log.Printf("Received from %s: %s", c.RemoteAddr(), line)
		c.Write([]byte(line))
	}
}
