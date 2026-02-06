package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	// Konfigurasi koneksi ke SeaweedFS Localhost
	endpoint := "127.0.0.1:9001"
	// accessKeyID := "any"     // SeaweedFS default: boleh isi bebas
	// secretAccessKey := "any" // SeaweedFS default: boleh isi bebas
	useSSL := false          // Kita pakai HTTP biasa di localhost

	// 1. Inisialisasi Client
	// Salah karena server belum punya database user
	// Benar: Masuk sebagai anonim (aman karena dilindungi localhost firewall)
	
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4("", "", ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln("Gagal konek:", err)
	}

	// 2. Buat Bucket (Folder) baru
	bucketName := "das-logs"
	ctx := context.Background()
	
	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		// Cek kalau errornya karena bucket sudah ada (itu normal)
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			fmt.Printf("Bucket '%s' sudah ada, lanjut...\n", bucketName)
		} else {
			log.Fatalln("Gagal bikin bucket:", err)
		}
	} else {
		fmt.Printf("Berhasil membuat bucket '%s'\n", bucketName)
	}

	// 3. Upload File Text Sederhana
	objectName := "test-upload.txt"
	fileContent := "Halo, ini data tes dari Golang ke SeaweedFS!"
	reader := strings.NewReader(fileContent)

	info, err := minioClient.PutObject(ctx, bucketName, objectName, reader, int64(reader.Len()), minio.PutObjectOptions{ContentType: "text/plain"})
	if err != nil {
		log.Fatalln("Gagal upload:", err)
	}

	fmt.Printf("SUKSES! File terupload ke %s dengan size %d bytes\n", objectName, info.Size)
}