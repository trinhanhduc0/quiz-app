package persistence

import (
	"context"
	"fmt"

	entity "quiz-app/internal/domain/entities"
	"quiz-app/internal/domain/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FileMongoRepository implements the repository.FileRepository interface
type FileMongoRepository struct {
	CollRepo repository.CRUDMongoDB
}

// NewFileMongoRepository creates a new instance of FileMongoRepository
func NewFileMongoRepository() repository.FileRepository {
	collRepo := NewCollRepository("dbapp", "files")
	return &FileMongoRepository{
		CollRepo: collRepo,
	}
}

// CreateFile implements repository.FileRepository.CreateFile
func (r *FileMongoRepository) CreateFile(ctx context.Context, file *entity.File) (primitive.ObjectID, error) {
	insertedID, err := r.CollRepo.Create(ctx, file)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("failed to create file: %w", err)
	}
	return insertedID.(primitive.ObjectID), nil
}

// GetAllFiles implements repository.FileRepository.GetAllFiles
func (r *FileMongoRepository) GetFile(ctx context.Context, file *entity.File) (any, error) {
	filter := bson.M{"metadata.email": file.Metadata.Email} // Assuming files are filtered by userID
	allFiles, err := r.CollRepo.GetAll(ctx, filter)         // Fetch all matching files
	if err != nil {
		return nil, fmt.Errorf("failed to get files by userID: %w", err)
	}
	return allFiles, nil
}

// GetAllFiles implements repository.FileRepository.GetAllFiles
func (r *FileMongoRepository) FindByName(ctx context.Context, file *entity.File) (any, error) {
	filter := bson.M{"metadata.email": file.Metadata.Email, "filename": file.Filename} // Assuming files are filtered by userID
	allFiles, err := r.CollRepo.GetAll(ctx, filter)                                    // Fetch all matching files
	if err != nil {
		return nil, fmt.Errorf("failed to get files by userID: %w", err)
	}
	return allFiles, nil
}

// GetAllFiles implements repository.FileRepository.GetAllFiles
func (r *FileMongoRepository) GetAllFile(ctx context.Context, email_id string) (any, error) {
	filter := bson.M{"metadata.emailid": email_id}  // Assuming files are filtered by userID
	allFiles, err := r.CollRepo.GetAll(ctx, filter) // Fetch all matching files
	if err != nil {
		return nil, fmt.Errorf("failed to get files by userID: %w", err)
	}

	return allFiles, nil
}

// GetAllImageFile implements repository.FileRepository.GetAllImageFile
func (r *FileMongoRepository) GetAllImageFile(ctx context.Context, email_id string) ([]any, error) {
	filter := bson.M{
		"metadata.emailid": email_id,
		"fileType": bson.M{
			"$regex": "^image/", // Lọc fileType bắt đầu bằng "image/"
		},
	}
	projection := bson.M{"filename": 1, "_id": 0}
	allFiles, err := r.CollRepo.GetWithProjection(ctx, filter, projection) // Fetch all matching files
	if err != nil {
		return nil, fmt.Errorf("failed to get image files for user %s: %w", email_id, err)
	}

	return allFiles, nil
}

// UpdateFile implements repository.FileRepository.UpdateFile
func (r *FileMongoRepository) UpdateFile(ctx context.Context, file *entity.File) (any, error) {
	filter := bson.M{"metadata.email": file.Metadata.Email, "_id": file.ID}

	_, err := r.CollRepo.Update(ctx, filter, bson.M{"$set": file})
	if err != nil {
		return &entity.File{}, fmt.Errorf("failed to update file: %w", err)
	}
	return file, nil
}

// DeleteFile implements repository.FileRepository.DeleteFile
func (r *FileMongoRepository) DeleteFile(ctx context.Context, file entity.File) error {
	filter := bson.M{"metadata.email": file.Metadata.Email, "_id": file.ID}
	result, err := r.CollRepo.Delete(ctx, filter)
	if err != nil && result.DeletedCount == 0 {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}
