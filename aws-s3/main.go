package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync/atomic"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type progressWriter struct {
	written int64
	writer  io.WriterAt
	size    int64
}

func (pw *progressWriter) WriteAt(p []byte, off int64) (int, error) {
	atomic.AddInt64(&pw.written, int64(len(p)))

	percentageDownloaded := float32(pw.written*100) / float32(pw.size)

	fmt.Printf("File size:%d downloaded:%d percentage:%.2f%%\r", pw.size, pw.written, percentageDownloaded)

	return pw.writer.WriteAt(p, off)
}

func byteCountDecimal(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}

func getFileSize(svc *s3.S3, bucket string, prefix string) (filesize int64, error error) {
	params := &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(prefix),
	}

	resp, err := svc.HeadObject(params)
	if err != nil {
		return 0, err
	}

	return *resp.ContentLength, nil
}

func parseFilename(keyString string) (filename string) {
	ss := strings.Split(keyString, "/")
	s := ss[len(ss)-1]
	return s
}

func main() {
	bucket := "test"
	key := "cn_windows_10_consumer_editions_version_20h2_x64_dvd_d4f7a83e.iso"

	filename := parseFilename(key)

	accessKey := ""
	secretKey := ""
	endPoint := ""

	sess, err := session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Endpoint:         aws.String(endPoint),
		Region:           aws.String("us-east-1"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	})
	if err != nil {
		panic(err)
	}

	s3Client := s3.New(sess)
	downloader := s3manager.NewDownloader(sess, func(d *s3manager.Downloader) {
		d.PartSize = 64 * 1024 * 1024
		d.Concurrency = 4
	})
	size, err := getFileSize(s3Client, bucket, key)
	if err != nil {
		panic(err)
	}

	log.Println("Starting download, size:", byteCountDecimal(size))
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	temp, err := ioutil.TempFile(cwd, "download-")
	if err != nil {
		panic(err)
	}
	tempfileName := temp.Name()

	writer := &progressWriter{writer: temp, size: size, written: 0}
	params := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	if _, err := downloader.Download(writer, params); err != nil {
		log.Printf("Download failed! Deleting tempfile: %s", tempfileName)
		os.Remove(tempfileName)
		panic(err)
	}

	if err := temp.Close(); err != nil {
		panic(err)
	}

	if err := os.Rename(temp.Name(), filename); err != nil {
		panic(err)
	}

	fmt.Println()
	log.Println("File downloaded! Available at:", filename)
}
