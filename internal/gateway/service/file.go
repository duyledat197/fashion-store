// Package service provides the implementation of the UploadService interface
// for handling file uploads through HTTP.
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
	// CHUNK_SIZE defines the size of each chunk when reading the file.
	CHUNK_SIZE = 1024
	// DEFAULT_MAX_MEMORY is the default maximum memory allocated for file uploads.
	DEFAULT_MAX_MEMORY = 32 << 20 // 32 MB
)

// UploadService is the interface for handling file uploads.
type UploadService interface {
	HandleUploadFiles(w http.ResponseWriter, r *http.Request, params map[string]string)
}

// uploadService is the implementation of the UploadService interface.
type uploadService struct {
	fileUploadClient pb.UploadServiceClient
	authenticator    token_util.JWTAuthenticator
}

// NewUploadService creates a new instance of the uploadService.
func NewUploadService(fileUploadClient pb.UploadServiceClient) UploadService {
	return &uploadService{
		fileUploadClient: fileUploadClient,
	}
}

// uploadFileReq represents the request structure for uploading a file.
type uploadFileReq struct {
	file     multipart.File
	header   *multipart.FileHeader
	resultCh chan *pb.File
	fileName string
}

// HandleUploadFiles handles the HTTP request for file uploads.
func (s *uploadService) HandleUploadFiles(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	// Parse the multipart form to retrieve the uploaded files.
	if err := r.ParseMultipartForm(DEFAULT_MAX_MEMORY); err != nil {
		http_server.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("unable to parse form: %w", err))
		return
	}

	ctx := r.Context()
	// Extract JWT token from the request headers.
	md := http_server.MapMetaDataWithBearerToken(s.authenticator)(ctx, r)
	ctx = metadata.NewOutgoingContext(ctx, md)

	// Use errgroup to handle concurrent file uploads.
	eg, _ := errgroup.WithContext(ctx)

	// Create channels to collect results and signal completion.
	var result []*pb.File
	resultCh := make(chan *pb.File)
	done := make(chan bool)

	// Iterate through the files in the multipart form.
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

	// Collect results from the result channel concurrently.
	go func() {
		for v := range resultCh {
			result = append(result, v)
		}
		done <- true
	}()

	// Wait for all file uploads to complete.
	if err := eg.Wait(); err != nil {
		http_server.ErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("unable to upload files: %w", err))
	}

	// Close the result channel and wait for completion.
	close(resultCh)
	<-done

	// Send the collected file upload results in the HTTP response.
	http_server.DataResponse(w, map[string]any{
		"data": result,
	})
}

// uploadFile handles the upload of a single file to the server.
func (s *uploadService) uploadFile(ctx context.Context, req *uploadFileReq) error {
	defer req.file.Close()

	// Get the MIME type of the file.
	fileType, err := fileutil.GetMimeTypeFile(req.file)
	if err != nil {
		return fmt.Errorf("unable to get mime type: %w", err)
	}

	// Obtain the filename from the file header.
	fileName := req.header.Filename

	// Initialize the gRPC stream for file upload.
	stream, err := s.fileUploadClient.Upload(ctx)
	if err != nil {
		return fmt.Errorf("unable to transfer image: %w", err)
	}

	// Create the initial upload request with file information.
	uReq := &pb.UploadRequest{
		Data: &pb.UploadRequest_Info_{
			Info: &pb.UploadRequest_Info{
				MimeType: fileType,
				Filename: fileName,
			},
		},
	}

	// Send the initial file information to the server.
	if err := stream.Send(uReq); err != nil {
		return fmt.Errorf("unable transfer image info to upload server: %w, %w", err, stream.RecvMsg(nil))
	}

	// Create a buffered reader for reading chunks from the file.
	reader := bufio.NewReader(req.file)
	buffer := make([]byte, CHUNK_SIZE)

	// Read and send file chunks until the end of the file is reached.
	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("cannot read chunk to buffer: %w, %w", err, stream.RecvMsg(nil))
		}

		// Create and send an upload request for the file chunk.
		req := &pb.UploadRequest{
			Data: &pb.UploadRequest_ChunkData{
				ChunkData: buffer[:n],
			},
		}

		if err := stream.Send(req); err != nil {
			return fmt.Errorf("cannot send chunk to server: %w, %w", err, stream.RecvMsg(nil))
		}
	}

	// Close the gRPC stream and receive the final response from the server.
	resp, err := stream.CloseAndRecv()
	if err != nil {
		return fmt.Errorf("unable to upload file to server: %w", err)
	}

	// Send the file data to the result channel for further processing.
	req.resultCh <- resp.GetData()

	return nil
}
