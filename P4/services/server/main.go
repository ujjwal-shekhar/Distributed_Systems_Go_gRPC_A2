package main

import (
	"context"
	"flag"
	"log"
	"time"

	// "time"

	pb "github.com/ujjwal-shekhar/bft/services/common/genproto/comms"
	"github.com/ujjwal-shekhar/bft/services/server/handler/server"
)

func main() {
	ID := flag.Int("ID", 1, "ID of the server")
	PORT := flag.Int("PORT", 5000, "Port of the server")
	TYPE := flag.String("TYPE", "honest", "Type of the server")
	N := flag.Int("N", 1, "Number of Voters")
	T := flag.Int("T", 1, "Number of Traitors")
	flag.Parse()

	log.Printf("args: %d %d %s %d %d", *ID, *PORT, *TYPE, *N, *T)

	// time.Sleep(5* time.Second) // Wait for everyone to get ready

	// Get a server type
	server, err := server.NewServer(*ID, *TYPE, *N, *T)
	if err != nil { panic(err) }
	log.Printf("Server %d of type %s created", *ID, *TYPE)

	go StartServer(*PORT, server)
	time.Sleep(5 * time.Second) // Wait for everyone to get ready
	
	server.ConnectToEveryone()
	time.Sleep(5 * time.Second) // Wait for everyone to get ready
	log.Printf("Server %d sent a connection request to all", *ID)
	
	if *ID == 0 {
		log.Printf("Server %d sending attack", *ID)
		for i, conn := range server.GeneralConns {
			if i == *ID { continue }

			log.Printf("Server %d sending attack to %v", *ID, conn)
			_, err := conn.SendValue(context.Background(), &pb.Value{
				Round: int32(0),
				Sender: int32(0),
				Attack: true,
				IsCommander: true,
			})

			if err != nil {
				log.Printf("Server %d failed to send attack to %v: %v", *ID, conn, err)
			}
		}
	}

	select {}
	// time.Sleep(5* time.Second) // Wait for everyone to get ready
}