package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	entity "quiz-app/internal/domain/entities"
	"quiz-app/internal/domain/service"
	"quiz-app/internal/infrastructure/persistence/aws"
	"quiz-app/internal/pkg"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoutesFile struct {
	auth         *service.AuthHandler
	fileUseCase  *service.FileUseCase
	awsS3UseCase *aws.FileAWSRepository
}

func NewRoutesFile(usecase *service.FileUseCase, awsS3UseCase *aws.FileAWSRepository, auth *service.AuthHandler) *RoutesFile {
	return &RoutesFile{
		fileUseCase:  usecase,
		awsS3UseCase: awsS3UseCase,
		auth:         auth,
	}
}

func (rf *RoutesFile) GetRoutesFile(r *Router) {
	r.Router.Handle("/getallimagefile", rf.auth.AuthMiddleware(http.HandlerFunc(rf.getAllImageFile))).Methods("GET")
	r.Router.Handle("/getimagefile", rf.auth.AuthMiddleware(http.HandlerFunc(rf.getImageFile))).Methods("GET")
	r.Router.Handle("/upimagefile", rf.auth.AuthMiddleware(http.HandlerFunc(rf.uploadImageFile))).Methods("POST")

	r.Router.Handle("/upfile", rf.auth.AuthMiddleware(http.HandlerFunc(rf.uploadFile))).Methods("POST")
	r.Router.Handle("/getfile", rf.auth.AuthMiddleware(http.HandlerFunc(rf.userGetFile))).Methods("POST")
	r.Router.Handle("/file", rf.auth.AuthMiddleware(http.HandlerFunc(rf.deleteFile))).Methods("DELETE")
}

func (rf *RoutesFile) uploadImageHandler(w http.ResponseWriter, r *http.Request, emailID, email string) {

	// Parse the multipart form to handle the file upload (10MB limit)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	// Retrieve the file and file handler
	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Println("Error retrieving the file:", err)
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate if the file is an image
	if !isImageFile(handler.Filename) {
		http.Error(w, "Only image files are allowed", http.StatusBadRequest)
		return
	}

	// Create file entity for MongoDB
	fileData := entity.NewFile(handler.Header.Get("Content-Type"), handler.Size, email, handler.Filename)
	fileData.Metadata.EmailID = emailID

	// Check if the file with the same name already exists in MongoDB
	existingFile, err := rf.fileUseCase.FindByName(entity.File{
		Metadata: entity.FileMetadata{
			EmailID: emailID,
			Email:   email,
		},
		Filename: handler.Filename,
	})

	if err != nil {
		log.Println("Error checking existing file:", err)
		http.Error(w, "Failed to check existing file", http.StatusInternalServerError)
		return
	}
	if len(existingFile.([]any)) > 0 {
		http.Error(w, "File with the same name already exists", http.StatusConflict)
		return
	}

	// Create file metadata in MongoDB
	createdFile, err := rf.fileUseCase.CreateFile(fileData)
	if err != nil {
		log.Println("Error creating file metadata:", err)
		http.Error(w, "Failed to create file metadata", http.StatusInternalServerError)
		return
	}

	// Reset file pointer to the start for S3 upload
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		log.Println("Error seeking file:", err)
		http.Error(w, "Error processing file", http.StatusInternalServerError)
		return
	}

	// Create file entity for AWS S3
	awsFile := &entity.File{
		Filename: handler.Filename,
		Metadata: entity.FileMetadata{
			Email: email,
		},
	}

	// Upload file to AWS S3
	_, err = rf.awsS3UseCase.CreateFile(r.Context(), awsFile, file)
	if err != nil {
		log.Println("Failed to upload file to S3:", err)

		// Rollback MongoDB entry if S3 upload fails
		if err := rf.fileUseCase.DeleteFile(r.Context(), *fileData); err != nil {
			log.Println("Error rolling back MongoDB entry:", err)
		}
		http.Error(w, "Failed to upload file to S3", http.StatusInternalServerError)
		return
	}

	// Send success response with created file metadata
	pkg.SendResponse(w, http.StatusOK, createdFile)
}

func (rf *RoutesFile) uploadImageFile(w http.ResponseWriter, r *http.Request) {
	email := r.Context().Value("email").(string)
	emailID := r.Context().Value("email_id").(string)
	rf.uploadImageHandler(w, r, emailID, email)
}

func (rf *RoutesFile) uploadFile(w http.ResponseWriter, r *http.Request) {
	email := r.Context().Value("email").(string)
	emailID := r.Context().Value("email_id").(string)

	// Parse the multipart form to handle the file upload (10MB limit)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	// Retrieve the file and file handler
	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Println("Error retrieving the file:", err)
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create file entity for MongoDB
	fileData := entity.NewFile(handler.Header.Get("Content-Type"), handler.Size, email, handler.Filename)
	fileData.Metadata.EmailID = emailID

	// Check if the file with the same name already exists in MongoDB
	existingFile, err := rf.fileUseCase.FindByName(entity.File{
		Metadata: entity.FileMetadata{
			EmailID: emailID,
			Email:   email,
		},
		Filename: handler.Filename,
	})

	if err != nil {
		log.Println("Error checking existing file:", err)
		http.Error(w, "Failed to check existing file", http.StatusInternalServerError)
		return
	}
	if len(existingFile.([]any)) > 0 {
		http.Error(w, "File with the same name already exists", http.StatusConflict)
		return
	}

	// Create file metadata in MongoDB
	createdFile, err := rf.fileUseCase.CreateFile(fileData)
	if err != nil {
		log.Println("Error creating file metadata:", err)
		http.Error(w, "Failed to create file metadata", http.StatusInternalServerError)
		return
	}

	// Reset file pointer to the start for S3 upload
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		log.Println("Error seeking file:", err)
		http.Error(w, "Error processing file", http.StatusInternalServerError)
		return
	}

	// Create file entity for AWS S3
	awsFile := &entity.File{
		Filename: handler.Filename,
		Metadata: entity.FileMetadata{
			Email: email,
		},
	}

	// Upload file to AWS S3
	_, err = rf.awsS3UseCase.CreateFile(r.Context(), awsFile, file)
	if err != nil {
		log.Println("Failed to upload file to S3:", err)

		// Rollback MongoDB entry if S3 upload fails
		if err := rf.fileUseCase.DeleteFile(r.Context(), *fileData); err != nil {
			log.Println("Error rolling back MongoDB entry:", err)
		}
		http.Error(w, "Failed to upload file to S3", http.StatusInternalServerError)
		return
	}

	// Send success response with created file metadata
	pkg.SendResponse(w, http.StatusOK, createdFile)
}

func (rf *RoutesFile) deleteFile(w http.ResponseWriter, r *http.Request) {
	email_id := r.Context().Value("email_id").(string)
	email := r.Context().Value("email").(string)

	var deleteFile struct {
		ID       primitive.ObjectID `bson:"_id" json:"_id"`
		Filename string             `bson:"filename" json:"filename"`
	}

	if err := json.NewDecoder(r.Body).Decode(&deleteFile); err != nil {
		pkg.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := rf.fileUseCase.DeleteFile(r.Context(), entity.File{
		Metadata: entity.FileMetadata{
			EmailID: email_id,
			Email:   email,
		},
		ID: deleteFile.ID,
	})

	if err != nil {
		http.Error(w, "Error delete file from db", http.StatusInternalServerError)
		return
	}

	err = rf.awsS3UseCase.DeleteFile(r.Context(), email, deleteFile.Filename)
	if err != nil {
		http.Error(w, "Error delete file from S3", http.StatusInternalServerError)
		return
	}

	pkg.SendResponse(w, http.StatusOK, true)

}

func (rf *RoutesFile) getAllImageFile(w http.ResponseWriter, r *http.Request) {
	emailID := r.Context().Value("email_id").(string)
	email := r.Context().Value("email").(string)

	// Lấy danh sách file metadata từ MongoDB
	rawFiles, err := rf.fileUseCase.GetAllImageFile(r.Context(), emailID)
	if err != nil {
		pkg.SendError(w, "Error getting all image files", http.StatusInternalServerError)
		return
	}

	// Tạo slice để lưu trữ thông tin từng file ảnh
	var fileContents []map[string]interface{}

	// Lặp qua các file và lấy nội dung từ S3
	for _, fileData := range rawFiles {
		fmt.Println(fileData)

		filename, ok := fileData.(primitive.M)["filename"].(string)
		if !ok {
			log.Println("Invalid filename format")
			continue
		}

		// Lấy nội dung file từ S3
		fileContent, err := rf.awsS3UseCase.GetFile(r.Context(), email, filename)
		if err != nil {
			log.Printf("Error retrieving file from S3: %v", err)
			continue
		}
		defer fileContent.Close()

		// Đọc toàn bộ nội dung file
		content, err := io.ReadAll(fileContent)
		if err != nil {
			log.Printf("Error reading file content: %v", err)
			continue
		}

		// Thêm thông tin file vào slice
		fileContents = append(fileContents, map[string]interface{}{
			"filename": filename,
			"content":  content, // hoặc encode thành base64 nếu cần thiết
		})
	}

	// Gửi response với danh sách file và nội dung
	pkg.SendResponse(w, http.StatusOK, fileContents)
}

func (rf *RoutesFile) getImageFile(w http.ResponseWriter, r *http.Request) {
	email := r.Context().Value("email").(string)
	filename := r.URL.Query().Get("filename")

	if filename == "" {
		http.Error(w, "Filename is required", http.StatusBadRequest)
		return
	}

	fileContent, err := rf.awsS3UseCase.GetFile(r.Context(), email, filename)
	if err != nil {
		http.Error(w, "Error retrieving file from S3", http.StatusInternalServerError)
		return
	}
	defer fileContent.Close()

	//w.Header().Set("Content-Type", fileInfo.Metadata.ContentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	if _, err := io.Copy(w, fileContent); err != nil {
		log.Printf("Error sending file content: %v", err)
	}
}

func (rf *RoutesFile) userGetFile(w http.ResponseWriter, r *http.Request) {
	var getFile struct {
		Email    string `json:"email"`
		Filename string `json:"filename"`
	}
	if err := json.NewDecoder(r.Body).Decode(&getFile); err != nil {
		pkg.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fileContent, err := rf.awsS3UseCase.GetFile(r.Context(), getFile.Email, getFile.Filename)
	if err != nil {
		http.Error(w, "Error retrieving file from S3", http.StatusInternalServerError)
		return
	}
	defer fileContent.Close()

	//w.Header().Set("Content-Type", fileInfo.Metadata.ContentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", getFile.Filename))
	if _, err := io.Copy(w, fileContent); err != nil {
		log.Printf("Error sending file content: %v", err)
	}
}

func isImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp":
		return true
	}
	return false
}
