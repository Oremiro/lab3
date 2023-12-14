package impl

import (
	"fmt"
	"io"
	"lab3/internal/infra/abs"
	"lab3/internal/model/entity"
	"lab3/pkg/guid"
	"os"
	"path/filepath"
	"time"
)

type FileService struct {
	repository   abs.Repository
	uploadFolder string
}

func NewFileService(repository abs.Repository, uploadFolder string) *FileService {
	return &FileService{repository: repository, uploadFolder: uploadFolder}
}

func (s *FileService) UploadFile(file io.Reader, originalFilename string) (string, error) {
	fileReferenceID := guid.NewGUID().String()

	filePath := filepath.Join(s.uploadFolder, fileReferenceID+"_"+originalFilename)

	out, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to save file: %v", err)
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return "", fmt.Errorf("failed to save file: %v", err)
	}

	metadata := entity.FileMetadata{
		FileReferenceID:  fileReferenceID,
		OriginalFilename: originalFilename,
		UploadTimestamp:  time.Now(),
	}

	err = s.repository.Insert(metadata)
	if err != nil {
		return "", fmt.Errorf("failed to insert file metadata: %v", err)
	}

	return fileReferenceID, nil
}

func (s *FileService) DownloadFile(fileReferenceID string) (string, error) {
	var metadata entity.FileMetadata
	err := s.repository.FindOneByField("file_reference_id", fileReferenceID, &metadata)
	if err != nil {
		return "", fmt.Errorf("failed to find file metadata: %v", err)
	}

	filePath := filepath.Join(s.uploadFolder, metadata.FileReferenceID+"_"+metadata.OriginalFilename)
	return filePath, nil
}
