package bootstrap

import (
	"expenses/internal/application/expense"

	grpcTransport "expenses/handler/grpc"
	pb "expenses/proto/pb"
	"log"
	"net"

	"google.golang.org/grpc"
)

func StartGRPCServer(port string, service expense.ServiceInterface) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", port, err)
	}

	s := grpc.NewServer()

	pb.RegisterExpenseServiceServer(s, grpcTransport.NewGrpcHandler(service))

	log.Println("gRPC server running", port)

	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
