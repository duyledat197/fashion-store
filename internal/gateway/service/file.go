// Package service ...
package service

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/metadata"

	pb "trintech/review/dto/storage-management/upload"
	fileutil "trintech/review/pkg/file_util"
	"trintech/review/pkg/http_server"
	"trintech/review/pkg/token_util"
)

const (
	// CHUNK_SIZE ...
	CHUNK_SIZE = 1024
	// DEFAULT_MAX_MEMORY ...
	DEFAULT_MAX_MEMORY = 32 << 20 // 32 MB
)

// UploadService ...
type UploadService interface {
	HandleUploadFiles(w http.ResponseWriter, r *http.Request, params map[string]string)
}

type uploadService struct {
	fileUploadClient pb.UploadServiceClient

	authenticator token_util.JWTAuthenticator
}

// NewUploadService ...
func NewUploadService(fileUploadClient pb.UploadServiceClient) UploadService {
	return &uploadService{
		fileUploadClient: fileUploadClient,
	}
}

type uploadFileReq struct {
	file     multipart.File
	header   *multipart.FileHeader
	resultCh chan *pb.File

	fileName string
}

func (s *uploadService) HandleUploadFiles(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	if err := r.ParseMultipartForm(DEFAULT_MAX_MEMORY); err != nil {
		http_server.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("unable to parse form: %w", err))
		return
	}

	ctx := r.Context()
	md := http_server.MapMetaDataWithBearerToken(s.authenticator)(ctx, r)
	ctx = metadata.NewOutgoingContext(ctx, md)

	eg, _ := errgroup.WithContext(ctx)

	var result []*pb.File
	resultCh := make(chan *pb.File)
	done := make(chan bool)

	for _, files := range r.MultipartForm.File {
		for _, h := range files {
			header := *h

			eg.Go(func() error {
				file, err := header.Open()
				if err != nil {
					return fmt.Errorf("unable to get file: %w", err)

				}
				return s.uploadFile(ctx, &uploadFileReq{
					file:     file,
					header:   &header,
					resultCh: resultCh,
				})
			})
		}
	}

	go func() {
		for v := range resultCh {
			result = append(result, v)
		}
		done <- true
	}()

	if err := eg.Wait(); err != nil {
		http_server.ErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("unable to upload files: %w", err))
	}

	close(resultCh)
	<-done

	http_server.DataResponse(w, map[string]any{
		"data": result,
	})
}

func (s *uploadService) uploadFile(ctx context.Context, req *uploadFileReq) error {
	defer req.file.Close()

	fileType, err := fileutil.GetMimeTypeFile(req.file)
	if err != nil {
		return fmt.Errorf("unable to get mime type: %w", err)
	}

	fileName := req.header.Filename

	stream, err := s.fileUploadClient.Upload(ctx)
	if err != nil {
		return fmt.Errorf("unable to transfer image: %w", err)
	}
	uReq := &pb.UploadRequest{
		Data: &pb.UploadRequest_Info_{
			Info: &pb.UploadRequest_Info{
				MimeType: fileType,
				Filename: fileName,
			},
		},
	}

	if err := stream.Send(uReq); err != nil {
		return fmt.Errorf("unable transfer image info to upload server: %w, %w", err, stream.RecvMsg(nil))
	}

	reader := bufio.NewReader(req.file)
	buffer := make([]byte, CHUNK_SIZE)

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("cannot read chunk to buffer: %w, %w", err, stream.RecvMsg(nil))
		}

		req := &pb.UploadRequest{
			Data: &pb.UploadRequest_ChunkData{
				ChunkData: buffer[:n],
			},
		}

		if err := stream.Send(req); err != nil {
			return fmt.Errorf("cannot send chunk to server: %w, %w", err, stream.RecvMsg(nil))
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		return fmt.Errorf("unable to upload file to server: %w", err)
	}

	req.resultCh <- resp.GetData()

	return nil
}
