package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"path/filepath"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/types/known/emptypb"

	// "userpb"
	// "github.com/scottstensland/ulfhednar/pkg/userpb"
	// userpb "github.com/scottstensland/ulfhednar/user"
	userpb "github.com/scottstensland/ulfhednar/user"
)

var jwtSecret = []byte("supersecretkey")

type server struct {
	userpb.UnimplementedUserServiceServer
	mu    sync.Mutex
	users map[string]*userpb.User
}

// CreateUser creates a new user
func (s *server) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := fmt.Sprintf("%d", time.Now().UnixNano())
	user := &userpb.User{
		Id:    id,
		Name:  req.Name,
		Email: req.Email,
	}
	s.users[id] = user
	return user, nil
}

// UpdateUser updates a user
func (s *server) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.users[req.Id]
	if !exists {
		return nil, errors.New("user not found")
	}

	user.Name = req.Name
	user.Email = req.Email
	return user, nil
}

// DeleteUser deletes a user
func (s *server) DeleteUser(ctx context.Context, req *userpb.DeleteUserRequest) (*emptypb.Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.users, req.Id)
	return &emptypb.Empty{}, nil
}

/*
func (s *server) DeleteUser(ctx context.Context, req *userpb.DeleteUserRequest) (*userpb.Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.users, req.Id)
	return &userpb.Empty{}, nil
}
*/

// GenerateJWT generates a JWT token
func (s *server) GenerateJWT(ctx context.Context, req *userpb.JWTRequest) (*userpb.JWTResponse, error) {
	token, refreshToken, err := generateTokens(req.Email)
	if err != nil {
		return nil, err
	}
	return &userpb.JWTResponse{Token: token, RefreshToken: refreshToken}, nil
}

// RefreshJWT renews a JWT token
func (s *server) RefreshJWT(ctx context.Context, req *userpb.JWTRequest) (*userpb.JWTResponse, error) {
	token, refreshToken, err := generateTokens(req.Email)
	if err != nil {
		return nil, err
	}
	return &userpb.JWTResponse{Token: token, RefreshToken: refreshToken}, nil
}

// Helper: Generate JWT and Refresh Token
func generateTokens(email string) (string, string, error) {
	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour).Unix(),
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}

	// Generate Refresh Token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	})
	refreshTokenString, err := refreshToken.SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}

	return tokenString, refreshTokenString, nil
}

func main() {

	certDir := "cert"
	serverCert := filepath.Join(certDir, "server-cert.pem")
	serverKey := filepath.Join(certDir, "server-key.pem")
	CACert := filepath.Join(certDir, "ca-cert.pem")
	// Load certificates for mTLS
	cert, err := tls.LoadX509KeyPair(serverCert, serverKey)
	if err != nil {
		log.Fatalf("Failed to load server cert/key: %v", err)
	}

	caCert, err := ioutil.ReadFile(CACert)
	if err != nil {
		log.Fatalf("Failed to load CA cert: %v", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	})

	// Start gRPC server
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.Creds(creds))
	userpb.RegisterUserServiceServer(grpcServer, &server{users: make(map[string]*userpb.User)})
	log.Println("Server running on port 50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
