package service

import (
	"context"
	entity "quiz-app/internal/domain/entities"
	"quiz-app/internal/domain/repository"
)

type FileUseCase struct {
	repo repository.FileRepository
}

func NewFileUseCase(repo repository.FileRepository) *FileUseCase {
	return &FileUseCase{repo: repo}
}

func (cuc *FileUseCase) CreateFile(file *entity.File) (any, error) {
	newID, err := cuc.repo.CreateFile(context.TODO(), file)
	file.ID = newID
	return file, err
}

func (cuc *FileUseCase) GetFile(file entity.File) (any, error) {
	return cuc.repo.GetFile(context.TODO(), &file)
}

func (cuc *FileUseCase) FindByName(file entity.File) (any, error) {
	return cuc.repo.FindByName(context.TODO(), &file)
}

func (cuc *FileUseCase) UpdateFile(file *entity.File) (any, error) {
	return cuc.repo.UpdateFile(context.TODO(), file)
}

func (cuc *FileUseCase) DeleteFile(ctx context.Context, file entity.File) error {
	return cuc.repo.DeleteFile(context.TODO(), file)
}

func (cuc *FileUseCase) GetAllImageFile(ctx context.Context, email string) ([]any, error) {
	return cuc.repo.GetAllImageFile(context.TODO(), email)
}

func (cuc *FileUseCase) GetAllFile(ctx context.Context, email string) (any, error) {
	return cuc.repo.GetAllFile(context.TODO(), email)
}
