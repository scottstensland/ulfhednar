package main

/*
	client gRPC ___________________________

	#  For uploading a JPEG image:


go run client.go  -mode=image -image=./media_client/Albrecht_Durer_De_Heilige_Familie_met_de_libelle.jpeg  -float=3.14 -int=42 -string="example"


	#  For uploading a WAV file:

go run client.go  -mode=sound -sound=./media_client/bach_three_lute_pieces_andres_segovia_mono_simple_flow.wav -another_float=2.71 -another_int=100 -another_string="another_example"



*/

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	pb "github.com/scottstensland/ulfhednar/enablers/helgustadanma/image_to_audio"

	"google.golang.org/grpc"
)

func main() {

	var mode string
	var imagePath, soundPath string
	var someFloat, anotherFloat float64
	var someInt, anotherInt int
	var someString, anotherString string

	flag.StringVar(&mode, "mode", "", "Mode of operation: 'image' or 'sound'")
	flag.StringVar(&imagePath, "image", "", "Path to the JPEG image file")
	flag.StringVar(&soundPath, "sound", "", "Path to the WAV sound file")
	flag.Float64Var(&someFloat, "float", 0.0, "A float value for image upload")
	flag.IntVar(&someInt, "int", 0, "An int value for image upload")
	flag.StringVar(&someString, "string", "", "A string value for image upload")
	flag.Float64Var(&anotherFloat, "another_float", 0.0, "A float value for sound upload")
	flag.IntVar(&anotherInt, "another_int", 0, "An int value for sound upload")
	flag.StringVar(&anotherString, "another_string", "", "A string value for sound upload")
	flag.Parse()

	if mode == "" {
		log.Fatal("Please specify a mode with -mode=image or -mode=sound")
	}

	serverAddr := os.Getenv("SERVER_PORT_GRPC_IMAGE_TO_SOUND")
	if serverAddr == "" {
		serverAddr = "localhost:50051"
		fmt.Println("SERVER_PORT_GRPC_IMAGE_TO_SOUND not set, using default value of local box")
	}

	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewImageToSoundServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	switch mode {
	case "image":
		if imagePath == "" {
			log.Fatal("Please specify an image path with -image=path/to/image.jpg")
		}
		imageData, err := ioutil.ReadFile(imagePath)
		if err != nil {
			log.Fatalf("could not open image: %v", err)
		}
		// r := &pb.ImageRequest{ImageData: imageData}

		r := &pb.ImageRequestWithMetadata{
			ImageData:  imageData,
			SomeFloat:  someFloat,
			SomeInt:    int32(someInt),
			SomeString: someString,
		}
		res, err := c.UploadImage(ctx, r)
		if err != nil {
			log.Fatalf("could not upload image: %v", err)
		}
		if err := ioutil.WriteFile("output.wav", res.SoundData, 0644); err != nil {
			log.Fatalf("could not write sound file: %v", err)
		}
		log.Println("WAV file written successfully")

	case "sound":
		if soundPath == "" {
			log.Fatal("Please specify a sound path with -sound=path/to/sound.wav")
		}
		soundData, err := ioutil.ReadFile(soundPath)
		if err != nil {
			log.Fatalf("could not open sound file: %v", err)
		}
		// soundReq := &pb.SoundRequest{SoundData: soundData}

		soundReq := &pb.SoundRequestWithMetadata{
			SoundData:     soundData,
			AnotherFloat:  anotherFloat,
			AnotherInt:    int32(anotherInt),
			AnotherString: anotherString,
		}
		res, err := c.UploadSound(ctx, soundReq)
		if err != nil {
			log.Fatalf("could not upload sound: %v", err)
		}
		if err := ioutil.WriteFile("output.jpeg", res.ImageData, 0644); err != nil {
			log.Fatalf("could not write image file: %v", err)
		}
		log.Println("JPEG file written successfully")

	default:
		log.Fatalf("Invalid mode: %s. Use 'image' or 'sound'.", mode)
	}
}
