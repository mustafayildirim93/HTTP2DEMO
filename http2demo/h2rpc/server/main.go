/*

package main

import (
	"fmt"
	"image"
	"log"
	"net"
	"path/filepath"

	"file_server/api"
	filepb "file_server/api"

	"gocv.io/x/gocv"
	"google.golang.org/grpc"
)

const (
	frameWidth  = 100
	frameHeight = 100
	port        = ":50051"
)

type fileServiceServer struct {
	api.UnimplementedFileServiceServer
}

func (s *fileServiceServer) FileDownLoad(req *filepb.FileInfo, stream filepb.FileService_FileDownLoadServer) error {
	filePath := filepath.Join(req.FileName + "." + req.FileExtension)

	img := gocv.IMRead(filePath, gocv.IMReadColor)
	defer img.Close()

	for y := 0; y < img.Rows(); y += frameHeight {
		for x := 0; x < img.Cols(); x += frameWidth {
			endX := x + frameWidth
			endY := y + frameHeight

			if endX > img.Cols() {
				endX = img.Cols()
			}
			if endY > img.Rows() {
				endY = img.Rows()
			}

			rect := image.Rect(x, y, endX, endY)
			frame := img.Region(rect)

			// Görüntüyü JPEG formatında encode et
			encodedFrame, err := gocv.IMEncode(".jpg", frame)
			if err != nil {
				return err
			}

			// Byte dizisini protobuf mesajına ekleyerek gönder
			content := &filepb.BytesContent{
				Buffer: encodedFrame.GetBytes(),
				Info:   req,
			}

			if err := stream.Send(content); err != nil {
				return err
			}

			frame.Close()
		}
	}

	return nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	filepb.RegisterFileServiceServer(grpcServer, &fileServiceServer{})

	fmt.Printf("Server is running on port %s\n", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

*/

package main

import (
	"fmt"
	"image"
	"log"
	"net"
	"path/filepath"

	"file_server/api"
	filepb "file_server/api"

	"gocv.io/x/gocv"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type fileServiceServer struct {
	api.UnimplementedFileServiceServer
}

func (s *fileServiceServer) FileDownLoad(req *filepb.FileInfo, stream filepb.FileService_FileDownLoadServer) error {
	filePath := filepath.Join(req.FileName + "." + req.FileExtension)

	img := gocv.IMRead(filePath, gocv.IMReadColor)
	defer img.Close()

	rows := img.Rows()
	cols := img.Cols()

	// Parça boyutları
	frameWidth := cols / 8  // 8'e bölerek eşit parçalara bölüyoruz
	frameHeight := rows / 8 // 8'e bölerek eşit parçalara bölüyoruz

	for y := 0; y < rows; y += frameHeight {
		for x := 0; x < cols; x += frameWidth {
			endX := x + frameWidth
			endY := y + frameHeight

			if endX > cols {
				endX = cols
			}
			if endY > rows {
				endY = rows
			}

			rect := image.Rect(x, y, endX, endY)
			frame := img.Region(rect)

			// Görüntüyü JPEG formatında encode et
			encodedFrame, err := gocv.IMEncode(".jpg", frame)
			if err != nil {
				return err
			}

			// Byte dizisini protobuf mesajına ekleyerek gönder
			content := &filepb.BytesContent{
				Buffer: encodedFrame.GetBytes(),
				Info:   req,
			}

			if err := stream.Send(content); err != nil {
				return err
			}

			frame.Close()
		}
	}

	return nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	filepb.RegisterFileServiceServer(grpcServer, &fileServiceServer{})

	fmt.Printf("Server is running on port %s\n", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
