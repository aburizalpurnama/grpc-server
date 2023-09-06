package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/aburizalpurnama/grpc-server/proto"
	"github.com/jinzhu/copier"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("APP_PORT")))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	proto.RegisterAccountsServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type server struct {
	proto.UnimplementedAccountsServer
}

func (s *server) SelectAccount(ctx context.Context, req *proto.SelectAccountRequest) (*proto.SelectAccountResponse, error) {

	db, err := sqlx.Connect("postgres", "postgres://test:test@localhost:5435/test?&sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	var accounts []Account
	var protoAccounts []*proto.Account

	err = db.SelectContext(ctx, &accounts, "SELECT id, name, balance FROM accounts")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("accounts: %v\n", accounts)

	copier.Copy(&protoAccounts, &accounts)

	fmt.Printf("protoAccounts: %v\n", protoAccounts)

	return &proto.SelectAccountResponse{Accounts: protoAccounts}, nil
}

type Account struct {
	Id      int32   `db:"id"`
	Name    string  `db:"name"`
	Balance float64 `db:"balance"`
}
