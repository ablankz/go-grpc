package main

import (
	"context"
	"fmt"
	"go-grpc/internal/config"
	"go-grpc/pb"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	cfg, err := config.Get()
	if err != nil {
		log.Fatalf("Config Error: %v", err)
	}

	conn, err := grpc.Dial("localhost:"+strconv.FormatInt(int64(cfg.Port), 10), grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	defer conn.Close()

	client := pb.NewFileServiceClient(conn)

	// callListFiles(client)
	// callDownload(client)
	// callUpload(client, *cfg)
	callUploadAndNotifyProgress(client, *cfg)
}

func callListFiles(client pb.FileServiceClient) {
	res, err := client.ListFiles(context.Background(), &pb.ListFilesRequest{})

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(res.GetFilenames())
}

func callDownload(client pb.FileServiceClient) {
	req := &pb.DownloadRequest{Filename: "name.txt"}
	stream, err := client.Download(context.Background(), req)
	if err != nil {
		log.Fatalln(err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("Response fro Download(bytes): %v", res.GetData())
		log.Printf("Response fro Download(string): %v", string(res.GetData()))
	}
}

func callUpload(client pb.FileServiceClient, cfg config.Config) {
	filename := "sports.txt"
	path := cfg.LocalRoot + "/storage/" + filename

	file, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	stream, err := client.Upload(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	buf := make([]byte, 5)
	for {
		n, err := file.Read(buf)
		if n == 0 || err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}

		req := &pb.UploadRequest{Data: buf[:n]}
		sendErr := stream.Send(req)
		if sendErr != nil {
			log.Fatalln(sendErr)
		}

		if cfg.Debug {
			time.Sleep(1 * time.Second)
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("received data size: %v", res.GetSize())
}

func callUploadAndNotifyProgress(client pb.FileServiceClient, cfg config.Config) {
	filename := "sports.txt"
	path := cfg.LocalRoot + "/storage/" + filename

	file, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	stream, err := client.UploadAndNotifyProgress(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	// request
	buf := make([]byte, 5)
	go func() {
		for {
			n, err := file.Read(buf)
			if n == 0 || err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalln(err)
			}

			req := &pb.UploadAndNotifyProgressRequest{Data: buf[:n]}
			sendErr := stream.Send(req)
			if sendErr != nil {
				log.Fatalln(sendErr)
			}

			if cfg.Debug {
				time.Sleep(1 * time.Second)
			}
		}

		err := stream.CloseSend()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	// response
	ch := make(chan struct{})
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalln(err)
			}

			log.Printf("received messsage: %v", res.GetMsg())
		}
		close(ch)
	}()
	<-ch
}
