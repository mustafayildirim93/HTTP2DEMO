/*
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	filepb "file_server/api" // Proto dosyanıza göre güncelle

	"google.golang.org/grpc"








































































	"gocv.io/x/gocv"
)

func main() {
	// Bağlanılacak gRPC sunucu adresi
	serverAddress := "localhost:50051"

	// Bağlantı oluştur
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Bağlantı oluşturulamadı: %v", err)
	}
	defer conn.Close()

	// FileService istemcisini oluştur
	client := filepb.NewFileServiceClient(conn)

	// RPC için istemci talebi oluştur
	request := &filepb.FileInfo{
		FileName:      "image",
		FileExtension: "png",
	}

	// RPC çağrısı için context oluştur
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// RPC çağrısı yap
	stream, err := client.FileDownLoad(ctx, request)
	if err != nil {
		log.Fatalf("FileDownLoad RPC'si başlatılamadı: %v", err)
	}

	// Stream üzerinde döngü ile mesajları al
	frameIndex := 0 // Frame index'i
	for {
		content, err := stream.Recv()
		if err == io.EOF {
			// Stream sona erdiğinde döngüyü sonlandır
			break
		}
		if err != nil {
			log.Fatalf("Mesaj alınamadı: %v", err)
		}

		// Gelen mesajdaki bilgileri ekrana yazdır
		fmt.Printf("Received Content: Size=%d, Read=%d, FileName=%s, FileExtension=%s\n",
		content.FileSize, content.ReadedByte, content.Info.FileName, content.Info.FileExtension)

		// Gelen veriyi işleyerek eğer gerekliyse örneğin bir resim dosyasına yazabilirsiniz
		// Örnek olarak, gocv.IMEncode ve gocv.IMWrite kullanabilirsiniz
		image, err := gocv.IMDecode(content.Buffer, gocv.IMReadColor)
		if err != nil {
			log.Fatalf("Resim okunamadı: %v", err)
		}
		defer image.Close()

		// Frame'i bir dosyaya yazmak istiyorsanız:
		frameFileName := fmt.Sprintf("frames/output_%d.jpg", frameIndex)
		gocv.IMWrite(frameFileName, image)

		frameIndex++
	}

	fmt.Println("İşlem tamamlandı.")
}

*/

package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	filepb "file_server/api" // Proto dosyanıza göre güncelle
	"gocv.io/x/gocv"
	"google.golang.org/grpc"
)

const (
	grpcServerAddress = "localhost:50051"
	httpServerAddress = "http://localhost:8081"
)

func grpcRequest(wg *sync.WaitGroup) {
	defer wg.Done()
	// Bağlantı oluştur
	conn, err := grpc.Dial(grpcServerAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("gRPC sunucusuna bağlanılamadı: %v", err)
	}
	defer conn.Close()

	// FileService istemcisini oluştur
	client := filepb.NewFileServiceClient(conn)

	// RPC için istemci talebi oluştur
	request := &filepb.FileInfo{
		FileName:      "image",
		FileExtension: "png",
	}

	// RPC çağrısı için context oluştur
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// RPC çağrısı yap ve süreyi hesapla
	startTime := time.Now()
	stream, err := client.FileDownLoad(ctx, request)
	if err != nil {
		log.Fatalf("FileDownLoad RPC'si başlatılamadı: %v", err)
	}
	duration := time.Since(startTime)

	// Stream üzerinde döngü ile mesajları al
	for {
		content, err := stream.Recv()
		if err == io.EOF {
			// Stream sona erdiğinde döngüyü sonlandır
			break
		}
		if err != nil {
			log.Fatalf("Mesaj alınamadı: %v", err)
		}

		// Gelen veriyi işleyerek eğer gerekliyse örneğin bir resim dosyasına yazabilirsiniz
		// Örnek olarak, gocv.IMEncode ve gocv.IMWrite kullanabilirsiniz
		image, err := gocv.IMDecode(content.Buffer, gocv.IMReadColor)
		if err != nil {
			log.Fatalf("Resim okunamadı: %v", err)
		}
		defer image.Close()

		//Eğer resmi dosyaya yazmak istiyorsanız:
		//resimAdi := fmt.Sprintf("grpc_output_%d.jpg", time.Now().UnixNano())
		//gocv.IMWrite(resimAdi, image)
	}

	fmt.Printf("gRPC'den gelen veri süresi: %v\n", duration)
}

func httpRequest(wg *sync.WaitGroup) {
	defer wg.Done()
	// HTTP isteği yap ve süreyi hesapla
	startTime := time.Now()
	response, err := http.Get(httpServerAddress)
	if err != nil {
		log.Fatalf("HTTP isteği başlatılamadı: %v", err)
	}
	defer response.Body.Close()
	duration := time.Since(startTime)

	// Gelen veriyi işleyerek eğer gerekliyse örneğin bir resim dosyasına yazabilirsiniz
	// Örnek olarak, gocv.IMEncode ve gocv.IMWrite kullanabilirsiniz
	data, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("HTTP verisi okunamadı: %v", err)
	}

	image, err := gocv.IMDecode(data, gocv.IMReadColor)
	if err != nil {
		log.Fatalf("Resim okunamadı: %v", err)
	}
	defer image.Close()

	//Eğer resmi dosyaya yazmak istiyorsanız:
	resimAdi := fmt.Sprintf("http_output_%d.jpg", time.Now().UnixNano())
	gocv.IMWrite(resimAdi, image)

	fmt.Printf("HTTP'den gelen veri süresi: %v\n", duration)
}

func main() {

	var wg sync.WaitGroup

	wg.Add(2) // Toplamda beklenen goroutine sayısını belirtiyoruz

	go grpcRequest(&wg)
	go httpRequest(&wg)

	// Tüm goroutinelerin bitmesini bekliyoruz
	wg.Wait()

	fmt.Println("Tüm goroutineler tamamlandı.")

}
