package s3

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/rs/zerolog/log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// ObjectConf for S3 object upload
type ObjectConf struct {
	Region *string
	Bucket *string
	Key    *string
	Body   *string
	GZip   bool
}

func sess(region *string) *session.Session {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Region: aws.String(*region)},
	}))
	return sess
}

func initUploader(region *string) *s3manager.Uploader {
	s := sess(region)
	return s3manager.NewUploader(s)
}

// Upload contents of body to an S3 bucket based on key
func Upload(conf *ObjectConf) error {
	myUploader := initUploader(conf.Region)

	if conf.GZip {
		gzipObject(conf)
	}

	// Upload the file to S3.
	_, err := myUploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(*conf.Bucket),
		Key:    aws.String(*conf.Key),
		Body:   strings.NewReader(*conf.Body),
	})
	if err != nil {
		return fmt.Errorf("failed to upload file, %v", err)
	}
	return nil
}

func initDownloader(region *string) *s3manager.Downloader {
	s := sess(region)
	return s3manager.NewDownloader(s)
}

func Download(conf *ObjectConf) (err error) {
	myDownloader := initDownloader(conf.Region)
	var buf []byte
	writeAt := aws.NewWriteAtBuffer(buf)
	_, err = myDownloader.Download(
		writeAt,
		&s3.GetObjectInput{
			Bucket: aws.String(*conf.Bucket),
			Key:    aws.String(*conf.Key),
		},
	)

	buf = writeAt.Bytes()
	if conf.GZip {
		gunzipObject(&buf)
	}

	body := string(buf)
	conf.Body = &body

	return
}

func gunzipObject(b *[]byte) {
	// TODO: implement
	return
}

func gzipObject(conf *ObjectConf) {
	addGzFileExtension(conf.Key)
	compressBody(conf.Body)
}

func compressBody(body *string) {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	if _, err := zw.Write([]byte(*body)); err != nil {
		log.Panic().Msg(err.Error())
	}
	if err := zw.Close(); err != nil {
		log.Panic().Msg(err.Error())
	}
	*body = string(buf.Bytes())
}

func addGzFileExtension(key *string) {
	*key = fmt.Sprintf("%s.gz", *key)
}
