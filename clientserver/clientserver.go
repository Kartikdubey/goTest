package clientserver

import (
	"fmt"

	"golang.org/x/net/context"
)

type Server struct {
}

func (s *Server) SendData(ctx context.Context, person *clientserver.Person) (*clientserver.Person, error) {
	fmt.Println("Received message body from client--", person)
	return &Person{Name: "ok"}, nil
}
