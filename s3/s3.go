// infra/s3_connection.go
package s3

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"mime/multipart"
	"net/url"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3Resp struct {
	FileName    string  `json:"file_name"`
	ContentType string  `json:"content_type"`
	FileUrl     string  `json:"file_url"`
	Size        float64 `json:"size"`
}

type S3Connection struct {
	AccessKey string
	SecretKey string
	Bucket    string
	Endpoint  string
}

func NewS3Connection(conf S3Connection) *S3Connection {
	return &S3Connection{
		AccessKey: conf.AccessKey,
		SecretKey: conf.SecretKey,
		Bucket:    conf.Bucket,
		Endpoint:  conf.Endpoint,
	}
}

func (s3 *S3Connection) getClient() (*minio.Client, error) {
	return minio.New(s3.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(s3.AccessKey, s3.SecretKey, ""),
		Secure: true,
	})
}

func (s3 *S3Connection) Upload(file *bytes.Reader, header *multipart.FileHeader, contentType string, collection string) (*S3Resp, error) {

	if collection == "" {
		return nil, fmt.Errorf("invalid collection name")
	}

	minioClient, err := s3.getClient()
	if err != nil {
		return nil, err
	}

	objectName := s3.createUniqueFilename(header.Filename)

	_, err = minioClient.PutObject(context.TODO(), s3.Bucket, objectName, file, header.Size, minio.PutObjectOptions{
		ContentType: contentType,
		UserMetadata: map[string]string{
			"x-amz-acl": "public-read",
		},
	})
	if err != nil {
		return nil, err
	}

	// fmt.Printf("File %s uploaded successfully to S3!\n", objectName)
	url, err := s3.PrintPublicURL(objectName)
	if err != nil {
		return nil, err
	}

	return &S3Resp{
		FileName:    objectName,
		ContentType: contentType,
		FileUrl:     url,
		Size:        float64(s3.convertByteToKB(int(header.Size))),
	}, nil
}

func (s3 *S3Connection) convertByteToKB(byteSize int) int {
	sizeInKB := float64(byteSize) / 1024
	return int(math.Round(sizeInKB))
}

func (s3 *S3Connection) createUniqueFilename(originalFilename string) string {
	extension := filepath.Ext(originalFilename)
	filename := fmt.Sprintf("%v-%v", strings.TrimSuffix(originalFilename, extension), uuid.New().String())
	// Menambahkan timestamp (Unix) ke nama file untuk membuatnya unik
	uniquePart := time.Now().Unix()
	// Menggabungkan nama file baru dengan ekstensi
	uniqueFilename := fmt.Sprintf("%s_%d%s", filename, uniquePart, extension)
	return uniqueFilename
}

func (s3 *S3Connection) PrintPublicURL(objectName string) (string, error) {
	minioClient, err := s3.getClient()
	if err != nil {
		return "", err
	}

	endPoint := minioClient.EndpointURL()
	return fmt.Sprintf("%v/%v/%v", endPoint, s3.Bucket, objectName), nil
}

func (s3 *S3Connection) PrintURL(objectName string) (string, error) {
	minioClient, err := s3.getClient()
	if err != nil {
		return "", err
	}

	// Mendapatkan URL publik objek dengan waktu kadaluwarsa 1 detik
	presignedURL, err := minioClient.PresignedGetObject(context.TODO(), s3.Bucket, objectName, time.Second*3600, nil)
	if err != nil {
		return "", err
	}

	url := presignedURL.String()
	return url, nil
}

func (s3 *S3Connection) GetObjectFromURL(presignedURL string) (string, error) {
	u, err := url.Parse(presignedURL)
	if err != nil {
		return "", err
	}

	// Mengambil path dari URL
	objectPath := u.Path
	objectName := path.Base(objectPath)
	return objectName, nil
}
