package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	// "userpb"
	userpb "github.com/scottstensland/ulfhednar/user"
)

func main() {

	certDir := "cert"
	clientCert := filepath.Join(certDir, "client-cert.pem")
	clientKey := filepath.Join(certDir, "client-key.pem")
	caCert := filepath.Join(certDir, "ca-cert.pem")

	// Load certificates for mTLS
	// cert, err := tls.LoadX509KeyPair("client-cert.pem", "client-key.pem")
	cert, err := tls.LoadX509KeyPair(clientCert, clientKey)
	if err != nil {
		log.Fatalf("Failed to load client cert/key: %v", err)
	}

	caCertBytes, errCA := ioutil.ReadFile(caCert)
	if err != nil {
		log.Fatalf("Failed to load CA cert: %v", errCA)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCertBytes)

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	})

	// Connect to gRPC server
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := userpb.NewUserServiceClient(conn)

	// Create a user
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.CreateUser(ctx, &userpb.CreateUserRequest{
		Name:  "Alice",
		Email: "alice@example.com",
	})
	if err != nil {
		log.Fatalf("CreateUser failed: %v", err)
	}
	fmt.Printf("Created User: %v\n", resp)

	// Generate JWT
	jwtResp, err := client.GenerateJWT(ctx, &userpb.JWTRequest{
		Email: "alice@example.com",
	})
	if err != nil {
		log.Fatalf("GenerateJWT failed: %v", err)
	}
	fmt.Printf("Generated JWT: %v\n", jwtResp)
}
