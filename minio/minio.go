package minio

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func MinioTest() {
	endpoint := "minio-api.jaquelee.com"
	accessKeyID := "jaquelee"
	secretAccessKey := "ljk#19901"
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("%#v\n", minioClient) // minioClient is now set up
	buckets, err := minioClient.ListBuckets(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	for _, bucket := range buckets {
		log.Println(bucket)
	}

	opts := minio.ListObjectsOptions{
		UseV1:     true,
		Prefix:    "",
		Recursive: true,
	}

	// List all objects from a bucket-name with a matching prefix.
	for object := range minioClient.ListObjects(context.Background(), "words", opts) {
		if object.Err != nil {
			fmt.Println(object.Err)
			return
		}
		fmt.Println(object)
	}

	reader, err := minioClient.GetObject(context.Background(), "words", "实现登录+文件上传.docx", minio.GetObjectOptions{})
	if err != nil {
		log.Fatalln(err)
	}
	defer reader.Close()

	localFile, err := os.Create("xxx.docx")
	if err != nil {
		log.Fatalln(err)
	}
	defer localFile.Close()

	stat, err := reader.Stat()
	if err != nil {
		log.Fatalln(err)
	}

	if _, err := io.CopyN(localFile, reader, stat.Size); err != nil {
		log.Fatalln(err)
	}

}
