package requestutil

import (
	"io"
	"mime/multipart"
	"net/http"
)

type MimeFile struct {
	Header  *multipart.FileHeader
	Content io.ReadCloser
}

func ExtractUpload(fileID string, r *http.Request) (*MimeFile, error) {

	// uploaded files should be stored to disk instead of heap
	// so we give the reader only one KB of buffer memory
	if err := r.ParseMultipartForm(1024); err != nil {
		return nil, err
	}

	file, header, err := r.FormFile(fileID)
	if err != nil {
		return nil, err
	}

	mimeFile := MimeFile{
		Header:  header,
		Content: file,
	}
	return &mimeFile, nil
}
