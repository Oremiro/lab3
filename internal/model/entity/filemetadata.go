package entity

import "time"

type FileMetadata struct {
	ID               string    `bson:"_id"`
	FileReferenceID  string    `bson:"file_reference_id"`
	OriginalFilename string    `bson:"original_filename"`
	UploadTimestamp  time.Time `bson:"upload_timestamp"`
}
