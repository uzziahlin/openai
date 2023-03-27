package openai

import "context"

type FileService interface {
	List(ctx context.Context) (*FileListResponse, error)
	Upload(ctx context.Context, req *FileUploadRequest) (*File, error)
	Delete(ctx context.Context, fileId string) (*FileDeleteResponse, error)
	Retrieve(ctx context.Context, fileId string) (*File, error)
	RetrieveContent(ctx context.Context, fileId string) (string, error)
}

type FileListResponse struct {
	Data []*File `json:"data"`
}

type FileUploadRequest struct {
	File    string `json:"file"`
	Purpose string `json:"purpose"`
}

type FileDeleteResponse struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Deleted bool   `json:"deleted"`
}

type File struct {
	Id        string `json:"id"`
	Object    string `json:"object"`
	Bytes     int64  `json:"bytes"`
	CreatedAt int64  `json:"created_at"`
	Filename  string `json:"filename"`
	Purpose   string `json:"purpose"`
}

type FileServiceOp struct {
	client *Client
}

func (f FileServiceOp) List(ctx context.Context) (*FileListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (f FileServiceOp) Upload(ctx context.Context, req *FileUploadRequest) (*File, error) {
	//TODO implement me
	panic("implement me")
}

func (f FileServiceOp) Delete(ctx context.Context, fileId string) (*FileDeleteResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (f FileServiceOp) Retrieve(ctx context.Context, fileId string) (*File, error) {
	//TODO implement me
	panic("implement me")
}

func (f FileServiceOp) RetrieveContent(ctx context.Context, fileId string) (string, error) {
	//TODO implement me
	panic("implement me")
}
