package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	// HTTP handler'ı tanımla
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Dosya yolu
		filePath := "image.png"

		// Dosyayı oku
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			http.Error(w, fmt.Sprintf("Dosya okuma hatası: %v", err), http.StatusInternalServerError)
			return
		}

		// Content-Type başlığını ayarla
		w.Header().Set("Content-Type", "image/png")

		// Byte dizisini yanıta yaz
		_, err = w.Write(data)
		if err != nil {
			http.Error(w, fmt.Sprintf("Yanıt yazma hatası: %v", err), http.StatusInternalServerError)
			return
		}
	})

	// Sunucuyu başlat
	port := 8081
	fmt.Printf("Server listening on :%d...\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Printf("Sunucu başlatma hatası: %v\n", err)
		os.Exit(1)
	}
}
