package main

import (
	"bytes"
	"context"
	"fmt"
	"go-grpc/internal/config"
	"go-grpc/pb"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedFileServiceServer
	cfg config.Config
}

func (s *server) ListFiles(ctx context.Context, req *pb.ListFilesRequest) (*pb.ListFilesResponse, error) {
	fmt.Println("ListFiles was invoked")

	dir := s.cfg.LocalRoot + "/storage"

	paths, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var filenames []string
	for _, path := range paths {
		if !path.IsDir() {
			filenames = append(filenames, path.Name())
		}
	}

	res := &pb.ListFilesResponse{
		Filenames: filenames,
	}
	return res, nil
}

func (s *server) Download(req *pb.DownloadRequest, stream pb.FileService_DownloadServer) error {
	fmt.Println("Download was invoked")

	filename := req.GetFilename()
	path := s.cfg.LocalRoot + "/storage/" + filename

	file, err := os.Open(path)
	if err != nil {
		return err
	}

	defer file.Close()

	buf := make([]byte, 5)
	for {
		n, err := file.Read(buf)
		if n == 0 || err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		res := &pb.DownloadResponse{Data: buf[:n]}
		sendErr := stream.Send(res)
		if sendErr != nil {
			return sendErr
		}

		if s.cfg.Debug {
			time.Sleep(1 * time.Second)
		}
	}

	return nil
}

func (*server) Upload(stream pb.FileService_UploadServer) error {
	fmt.Println("Upload was invoked")

	var buf bytes.Buffer
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			res := &pb.UploadResponse{Size: int32(buf.Len())}
			return stream.SendAndClose(res)
		}
		if err != nil {
			return err
		}

		data := req.GetData()
		log.Printf("received data(bytes): %v", data)
		log.Printf("received data(string): %v", string(data))
		buf.Write(data)
	}
}

func (*server) UploadAndNotifyProgress(stream pb.FileService_UploadAndNotifyProgressServer) error {
	fmt.Println("UploadAndNotifyProgress was invoked")

	size := 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		data := req.GetData()
		log.Printf("received data: %v", data)

		size += len(data)

		res := &pb.UploadAndNotifyProgressResponse{
			Msg: []byte(fmt.Sprintf("received %vbytes", size)),
		}
		err = stream.Send(res)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	cfg, err := config.Get()
	if err != nil {
		log.Fatalf("Config Error: %v", err)
	}

	lis, err := net.Listen("tcp", "localhost:"+strconv.FormatInt(int64(cfg.Port), 10))
	if err != nil {
		log.Fatalf("Failed to Listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterFileServiceServer(s, &server{cfg: *cfg})

	fmt.Println("server is running")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
