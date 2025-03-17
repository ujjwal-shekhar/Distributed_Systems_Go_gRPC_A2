package server

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"sync"

	pb "github.com/ujjwal-shekhar/bft/services/common/genproto/comms"
	"github.com/ujjwal-shekhar/bft/services/common/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	pb.UnimplementedOMServer
	GeneralConns 	[]pb.OMClient

	ID   			int
	TYPE 			string
	N    			int
	T    			int

	valuesLock sync.Mutex
	values     map[int32][]bool // Stores values received for each round
}

func NewServer(ID int, TYPE string, N int, T int) (*Server, error) {
	server := &Server{
		GeneralConns: make([]pb.OMClient, N),
		ID: ID,
		TYPE: TYPE,
		N: N,
		T: T,

		valuesLock: sync.Mutex{},
		values:     make(map[int32][]bool),
	}

	for round := range N+1 {
		server.values[int32(round)] = make([]bool, N)
		for i := range N {
			server.values[int32(round)][i] = false
		}
	}

	return server, nil
}

func (s *Server) ConnectToEveryone() {
	s.valuesLock.Lock()
	defer s.valuesLock.Unlock()

	log.Printf("Server %d connecting to everyone", s.ID)

	// Connect to everyone
	for i := range s.N {
		if i == s.ID { continue } // Duh

		// Connect, and then store
		conn, err := grpc.NewClient(
			"localhost:"+strconv.Itoa(utils.PORT_BASE + i), 
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil { panic(err) }

		s.GeneralConns[i] = pb.NewOMClient(conn)
		log.Printf("Server %d connected to %d, on port %s", s.ID, i, strconv.Itoa(utils.PORT_BASE + i))
	}

	log.Printf("Server %d connected to everyone", s.ID)
}

func (s *Server) SendValue(ctx context.Context, value *pb.Value) (*emptypb.Empty, error) {
	log.Printf("Server %d received value from %d: %v", s.ID, value.Sender, value.Attack)

	// If this is the last round, compute the majority value
	// and stop the protocol
	majorityValue := s.computeMajority(value.Round)
	s.logMajorityValue(value.Round, majorityValue)
	if value.Round == int32(s.T) {
		return &emptypb.Empty{}, nil
	}

	if value.IsCommander { // Commander sent this value
		s.valuesLock.Lock()
		s.values[value.Round][value.Sender] = value.Attack
		s.valuesLock.Unlock()

		// Forward this value to all lieutenants (excluding self and commander)
		for i := range s.N {
			if i == s.ID || i == 0 { continue }

			// Poison the value if traitor
			attack := value.Attack
			if s.TYPE == "traitor" {
				attack = rand.Intn(2) == 0
			}

			s.GeneralConns[i].SendValue(
				ctx,
				&pb.Value{
					Round:  value.Round+1,
					Sender: int32(s.ID),
					Attack: attack,
					IsCommander: true,
				},
			)
		}
	} else {
		// Forward this value to all lieutenants (excluding self and commander)
		for i := range s.N {
			if i == s.ID || i == 0 { continue }

			s.valuesLock.Lock()
			s.values[value.Round][value.Sender] = s.computeMajority(value.Round)
			s.valuesLock.Unlock()

			// s.GeneralConns[i].SendValue(
			// 	ctx,
			// 	&pb.Value{
			// 		Round:  value.Round+1,
			// 		Sender: int32(s.ID),
			// 		Attack: s.getValue(value.Round - 1),
			// 		IsCommander: false,
			// 	},
			// )
		}
	}

	log.Printf("Server %d sent value to all", s.ID)

	return &emptypb.Empty{}, nil
}


func (s *Server) computeMajority(round int32) bool {
	values := s.values[round]
	count := 0
	for _, v := range values {
		if v {
			count++
		}
	}
	return count > len(values)/2
}

// func (s *Server) getValue(roundID int32) bool {
// 	if s.TYPE == "honest" {
// 		// Wait for all values for this round
// 		for len(s.values[roundID]) < s.N-2 {
// 			// Wait for all lieutenants to send their values
// 		}
// 		return s.computeMajority(roundID)
// 	} else {
// 		// Traitor randomly decides
// 		return rand.Intn(2) == 0
// 	}
// }

func (s *Server) logMajorityValue(round int32, majorityValue bool) {
	// Create or open the log file for writing
	fileName := fmt.Sprintf("%d.out", s.ID)
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	// Write the majority value to the file
	_, err = file.WriteString(fmt.Sprintf("Round %d: [%s] Majority Value = %v\n", round, s.TYPE, majorityValue))
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
	}
}