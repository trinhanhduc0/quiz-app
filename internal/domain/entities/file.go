package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// File đại diện cho thực thể tệp tin được lưu trữ trong MongoDB
type File struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FileType   string             `bson:"fileType" json:"fileType"`     // Kiểu file (image, pdf, ...)
	Size       int64              `bson:"size" json:"size"`             // Kích thước file (bytes)
	UploadDate time.Time          `bson:"uploadDate" json:"uploadDate"` // Ngày tải lên
	Metadata   FileMetadata       `bson:"metadata" json:"metadata"`     // Metadata liên quan đến file
	Filename   string             `bson:"filename" json:"filename"`     // Tên file gốc
	Url        string             `bson:"url" json:"url"`
}

// FileMetadata lưu trữ thông tin bổ sung về file
type FileMetadata struct {
	Email   string `bson:"email" json:"email"` // Email người tải lên
	EmailID string
}

// NewFile là factory function để tạo một file mới
func NewFile(fileType string, size int64, email, filename string) *File {
	return &File{
		FileType: fileType,
		Size:     size,
		Metadata: FileMetadata{
			Email: email,
		},
		Filename: filename,
	}
}
