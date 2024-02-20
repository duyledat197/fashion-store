// Package service ...
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

type storageService struct {
	storage storage.Storage

	fileRepo interface {
		Create(ctx context.Context, db database.Executor, data *entity.File) error
	}

	db database.Database
	pb.UnimplementedUploadServiceServer
}

func (s *storageService) Upload(stream pb.UploadService_UploadServer) error {
	ctx := stream.Context()
	userCtx, err := xcontext.ExtractUserInfoFromContext(ctx)
	if err != nil {
		return status.Errorf(codes.PermissionDenied, "cant extract your information from request")
	}
	req, err := stream.Recv()
	if err != nil {
		return status.Errorf(codes.Unknown, "cannot receive image info")
	}

	var (
		fileName = req.GetInfo().GetFilename()
		mimeType = req.GetInfo().GetMimeType()
	)
	fileData := bytes.Buffer{}
	fileSize := 0

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			break
		}
		if err != nil {
			return status.Errorf(codes.Unknown, "unable to receive chunk data: %v", err)
		}
		chunk := req.GetChunkData()
		size := len(chunk)

		fileSize += size

		if _, err := fileData.Write(chunk); err != nil {
			return status.Errorf(codes.Internal, "unable to write chunk data: %v", err)
		}
	}

	url, err := s.storage.UploadObject(ctx, fileName, &fileData)
	if err != nil {
		return status.Errorf(codes.Internal, "unable to write chunk data: %v", err)
	}

	if err := s.fileRepo.Create(ctx, s.db, &entity.File{
		ID:        pg_util.NullString(uuid.NewString()),
		FileName:  pg_util.NullString(fileName),
		MimeType:  pg_util.NullString(mimeType),
		Size:      pg_util.NullInt64(int64(fileSize)),
		CreatedBy: pg_util.NullInt64(userCtx.UserID),
	}); err != nil {
		return status.Errorf(codes.Internal, "unable to create file: %v", err)
	}

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
