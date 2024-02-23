// Package service provides functionality related to storage services.
package service

import (
	"bytes"
	"context"
	"io"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "trintech/review/dto/storage-management/upload"
	"trintech/review/internal/storage-management/entity"
	"trintech/review/pkg/database"
	"trintech/review/pkg/http_server/xcontext"
	"trintech/review/pkg/pg_util"
	"trintech/review/pkg/storage"
)

// storageService handles file upload functionality using gRPC streaming.
type storageService struct {
	storage  storage.Storage
	fileRepo interface {
		Create(ctx context.Context, db database.Executor, data *entity.File) error
	}
	db database.Database
	pb.UnimplementedUploadServiceServer
}

// Upload handles the file upload gRPC streaming method.
func (s *storageService) Upload(stream pb.UploadService_UploadServer) error {
	ctx := stream.Context()

	// Extract user information from the context.
	userCtx, err := xcontext.ExtractUserInfoFromContext(ctx)
	if err != nil {
		return status.Errorf(codes.PermissionDenied, "cant extract your information from request")
	}

	// Receive the initial request to get file information.
	req, err := stream.Recv()
	if err != nil {
		return status.Errorf(codes.Unknown, "cannot receive image info")
	}

	// Extract file information from the request.
	var (
		fileName = req.GetInfo().GetFilename()
		mimeType = req.GetInfo().GetMimeType()
	)

	// Buffer to store file data.
	fileData := bytes.Buffer{}
	fileSize := 0

	// Loop through the stream to receive file chunks.
	for {
		req, err := stream.Recv()

		if err == io.EOF {
			break
		}
		if err != nil {
			return status.Errorf(codes.Unknown, "unable to receive chunk data: %v", err)
		}

		// Extract and process the received chunk data.
		chunk := req.GetChunkData()
		size := len(chunk)
		fileSize += size

		if _, err := fileData.Write(chunk); err != nil {
			return status.Errorf(codes.Internal, "unable to write chunk data: %v", err)
		}
	}

	// Upload the file to the storage provider.
	url, err := s.storage.UploadObject(ctx, fileName, &fileData)
	if err != nil {
		return status.Errorf(codes.Internal, "unable to write chunk data: %v", err)
	}

	// Create a file record in the database.
	if err := s.fileRepo.Create(ctx, s.db, &entity.File{
		ID:        pg_util.NullString(uuid.NewString()),
		FileName:  pg_util.NullString(fileName),
		MimeType:  pg_util.NullString(mimeType),
		Size:      pg_util.NullInt64(int64(fileSize)),
		CreatedBy: pg_util.NullInt64(userCtx.UserID),
	}); err != nil {
		return status.Errorf(codes.Internal, "unable to create file: %v", err)
	}

	// Send the response to the client with file information.
	if err := stream.SendAndClose(&pb.UploadResponse{
		Data: &pb.File{
			MimeType: mimeType,
			Size:     int64(fileSize),
			Url:      url,
		},
	}); err != nil {
		return status.Errorf(codes.Unknown, "unable to response: %v", err)
	}

	return nil
}
