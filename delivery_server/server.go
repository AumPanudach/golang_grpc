package main

import (
	"context"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"

	// This will need to be updated once the protoc-generated code is available
	pb "go_grpc/grpc_server/delivery"
)

const (
	// Port is the port on which the server will listen
	Port = ":50052"
)

// DeliveryServer implements the DeliveryService gRPC service
type DeliveryServer struct {
	pb.UnimplementedDeliveryServiceServer
	mu     sync.Mutex
	orders map[int32]*Order
	nextID int32
}

// Order represents a delivery order in the system
type Order struct {
	ID           int32
	CustomerName string
	Address      string
	Item         string
	Status       string
}

// NewDeliveryServer creates a new delivery server instance
func NewDeliveryServer() *DeliveryServer {
	return &DeliveryServer{
		orders: make(map[int32]*Order),
		nextID: 1,
	}
}

// CreateOrder implements the CreateOrder RPC method
func (s *DeliveryServer) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Create a new order
	orderID := s.nextID
	s.nextID++

	// Store the order
	s.orders[orderID] = &Order{
		ID:           orderID,
		CustomerName: req.CustomerName,
		Address:      req.Address,
		Item:         req.Item,
		Status:       "PENDING",
	}

	log.Printf("Created order %d for customer %s", orderID, req.CustomerName)
	return &pb.CreateOrderResponse{
		OrderId: orderID,
		Status:  "PENDING",
	}, nil
}

// GetOrderStatus implements the GetOrderStatus RPC method
func (s *DeliveryServer) GetOrderStatus(ctx context.Context, req *pb.GetOrderStatusRequest) (*pb.GetOrderStatusResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Look up the order
	order, exists := s.orders[req.OrderId]
	if !exists {
		log.Printf("Order %d not found", req.OrderId)
		return &pb.GetOrderStatusResponse{
			Status: "NOT_FOUND",
		}, nil
	}

	log.Printf("Order %d status: %s", req.OrderId, order.Status)
	return &pb.GetOrderStatusResponse{
		Status: order.Status,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", Port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Register our delivery service implementation
	deliveryServer := NewDeliveryServer()
	pb.RegisterDeliveryServiceServer(grpcServer, deliveryServer)

	log.Printf("Delivery server is listening on port %s", Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
