package main

import (
	"context"
	"io/ioutil"
	"log"
	"net"

	pb "github.com/scottstensland/ulfhednar/enablers/helgustadanma/image_to_audio"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedImageToSoundServiceServer
}

// UploadImage handles the upload of an image and returns a sound file
// func (s *server) UploadImage(ctx context.Context, req *pb.ImageRequest) (*pb.SoundResponse, error) {
func (s *server) UploadImage(ctx context.Context, req *pb.ImageRequestWithMetadata) (*pb.SoundResponse, error) {
	// Same as before, just reading a static WAV file
	soundData, err := ioutil.ReadFile("./media_server/ostrich_chick_pulses.wav") // respond to client with this sound
	if err != nil {
		return nil, err
	}
	return &pb.SoundResponse{SoundData: soundData}, nil
}

// UploadSound handles the upload of a sound file and returns an image
// func (s *server) UploadSound(ctx context.Context, req *pb.SoundRequest) (*pb.ImageResponse, error) {
func (s *server) UploadSound(ctx context.Context, req *pb.SoundRequestWithMetadata) (*pb.ImageResponse, error) {
	// Here you would process the sound data to generate an image. For simplicity, we'll just read a static JPEG file.
	imageData, err := ioutil.ReadFile("./media_server/knights_templar_flag.jpeg") // respond to client with this image
	if err != nil {
		return nil, err
	}
	return &pb.ImageResponse{ImageData: imageData}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// s := grpc.NewServer()
	s := grpc.NewServer(
		grpc.MaxRecvMsgSize(1024*1024*10), // 10MB limit
		grpc.MaxSendMsgSize(1024*1024*10), // 10MB limit
	)
	pb.RegisterImageToSoundServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
