package main

import (
	"flag"
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/kyleconroy/stripe-grpc/proto"
)

var (
	serverAddr = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	apiKey     = flag.String("api_key", "sk_test_xxxx", "The Stripe API key you'd like to use")
	charge     = flag.String("charge", "ch_123", "The Stripe charge to fetch")
)

type StripeCredentials struct {
	Key string
}

func (sc StripeCredentials) GetRequestMetadata(ctx context.Context) (map[string]string, error) {
	return map[string]string{"apikey": sc.Key}, nil
}

func main() {
	flag.Parse()
	var opts []grpc.DialOption
	if *apiKey != "" {
		creds := StripeCredentials{*apiKey}
		opts = append(opts, grpc.WithPerRPCCredentials(creds))
	}
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewStripeClient(conn)

	charge, err := client.GetCharge(context.Background(), &pb.Resource{*charge})
	if err != nil {
		log.Fatalf("%v.GetCharge(_) = _, %v: ", client, err)
	}
	log.Println(charge)
}
