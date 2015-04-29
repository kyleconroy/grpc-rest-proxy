package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/kyleconroy/stripe-grpc/proto"

	_ "github.com/stripe/stripe-go"
	stripe "github.com/stripe/stripe-go/client"
)

var (
	port = flag.Int("port", 10000, "The server port")
)

type stripeServer struct {
}

func stripeClient(ctx context.Context) (*stripe.API, error) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return nil, errors.New("missing RPC credentials")
	}
	sc := &stripe.API{}
	sc.Init(md["apikey"], nil)
	return sc, nil
}

// GetFeature returns the feature at the given point.
func (s *stripeServer) GetCharge(ctx context.Context, resource *pb.Resource) (*pb.Charge, error) {
	sc, err := stripeClient(ctx)
	if err != nil {
		return nil, err
	}

	charge, err := sc.Charges.Get(resource.Id, nil)
	if err != nil {
		return nil, err
	}

	return &pb.Charge{Id: charge.ID}, nil
}

// ListFeatures lists all features comtained within the given bounding Rectangle.
func newServer() *stripeServer {
	return new(stripeServer)
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterStripeServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}
