package aws

import (
	"context"
	"fmt"
	"io"
	"log"
	entity "quiz-app/internal/domain/entities"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileAWSRepository struct {
	S3Client *s3.S3
	Bucket   string
}

func NewFileAWSRepository(bucket string, region string) *FileAWSRepository {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		log.Fatal("Failed to create AWS session", err)
	}
	return &FileAWSRepository{
		S3Client: s3.New(sess),
		Bucket:   bucket,
	}
}

func (r *FileAWSRepository) CreateFile(ctx context.Context, file *entity.File, body io.ReadSeeker) (primitive.ObjectID, error) {
	fileID := primitive.NewObjectID()
	userPrefix := file.Metadata.Email + "/"
	fileKey := userPrefix + file.Filename

	input := &s3.PutObjectInput{
		Bucket: aws.String(r.Bucket),
		Key:    aws.String(fileKey),
		Body:   body,
		Metadata: map[string]*string{
			"email": aws.String(file.Metadata.Email),
		},
	}

	_, err := r.S3Client.PutObjectWithContext(ctx, input)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("failed to upload file: %v", err)
	}

	return fileID, nil
}

func (r *FileAWSRepository) GetFile(ctx context.Context, email string, filename string) (io.ReadCloser, error) {
	userPrefix := email + "/"
	fileKey := userPrefix + filename

	input := &s3.GetObjectInput{
		Bucket: aws.String(r.Bucket),
		Key:    aws.String(fileKey),
	}

	result, err := r.S3Client.GetObjectWithContext(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get file: %v", err)
	}

	return result.Body, nil
}

func (r *FileAWSRepository) UpdateFile(ctx context.Context, email string, file *entity.File, body io.ReadSeeker) (primitive.ObjectID, error) {
	err := r.DeleteFile(ctx, email, file.Filename)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("failed to delete old file: %v", err)
	}

	return r.CreateFile(ctx, file, body)
}

func (r *FileAWSRepository) DeleteFile(ctx context.Context, email string, filename string) error {
	userPrefix := email + "/"
	fileKey := userPrefix + filename

	input := &s3.DeleteObjectInput{
		Bucket: aws.String(r.Bucket),
		Key:    aws.String(fileKey),
	}

	_, err := r.S3Client.DeleteObjectWithContext(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
	}

	return nil
}
