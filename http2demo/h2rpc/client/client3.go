package main

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"

	filepb "file_server/api" // proto dosyanıza göre düzenleyin

	"google.golang.org/grpc"
)

func main() {
	// Sunucu adresi
	serverAddress := "localhost:50051" // Sunucu adresini güncelleyin

	// Bağlantı oluştur
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Sunucuya bağlanılamadı: %v", err)
	}
	defer conn.Close()

	// gRPC istemcisini oluştur
	client := filepb.NewFileServiceClient(conn)

	// Sunucudan dosya bilgilerini ve içeriğini al
	fileInfo := &filepb.FileInfo{
		FileName:      "image", // Dosya adını güncelleyin
		FileExtension: "png",   // Dosya uzantısını güncelleyin
	}

	stream, err := client.FileDownLoad(context.Background(), fileInfo)
	if err != nil {
		log.Fatalf("Dosya indirme isteği başlatılamadı: %v", err)
	}

	// Dosya içeriğini al
	var content []byte
	for {
		data, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Dosya içeriği alınamadı: %v", err)
		}

		// Gelen veriyi content'e ekle
		content = append(content, data.Buffer...)
	}

	// HTTP sunucusunu başlat ve veriyi yayınla
	startHTTPServer(content)
}

func startHTTPServer(content []byte) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.New("index").Parse(`
			<!--<!DOCTYPE html>
			<html>
				<head>
					<title>GRPC Data Viewer</title>
				</head>
				<body style="background-color: #1e3799; color: white; text-align: center;">
					<h1>GRPC Data Viewer</h1>
					<div style="padding: 20px;">
						<pre>{{ . }}</pre>
					</div>
				</body>
			</html>-->
			<head>
			<title>GRPC Data Viewer</title>
<meta name="viewport" content="width=device-width, initial-scale=1">
<style>
* {
  box-sizing: border-box;
}

/* Create two equal columns that floats next to each other */
.column {
  float: left;
  width: 50%;
  padding: 10px;
  margin: %5;
  height: 100%; /* Should be removed. Only for demonstration */
  border-radius: 10%;
  text-align:center;
}
/* Example styling for buttons inside .column */
.column button {
  background-color: #4CAF50;
  color: smokegray;
  padding: 10px 20px;
  border: none;
  transition: all 0.3s ease;
  border-radius: 5px;
  cursor: pointer;
}
.column button:hover {
	
	transform: scale(1.1); /* Increase the size on hover */
  }


/* Clear floats after the columns */
.row:after {
  content: "";
  display: table;
  clear: both;
}

/* Responsive layout - makes the two columns stack on top of each other instead of next to each other */
@media screen and (max-width: 600px) {
  .column {
    width: 100%;
  }
}
</style>
</head>
<body>

<h2>Responsive Two Column Layout</h2>
<p>Resize the browser window to see the responsive effect (the columns will stack on top of each other instead of floating next to each other, when the screen is less than 600px wide).</p>

<div class="row">
  <div class="column" style="background-color:rgb(252,165,165)">
  <h2>TEST HTTP/2 DEMO</h2>
  	<button onclick="sendHTTPRequest()">Send HTTP Request</button>
  </div>
  <div class="column" style="background-color:#bbb;">
  	<h2>TEST HTTP/1.1 DEMO</h1>
  	<button onclick="downloadImage()">Download Image</button>
  </div>
</div>

</body>
</html>
		`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// HTML sayfasını oluştur
		err = tmpl.Execute(w, string(content))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	fmt.Println("HTTP Server is running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("HTTP Server couldn't start: %v", err)
	}
}
