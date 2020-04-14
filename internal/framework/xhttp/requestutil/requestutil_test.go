package requestutil_test

import (
	"bytes"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ExperienceOne/apikit/internal/framework/xhttp/requestutil"
)

func TestExtractUpload(t *testing.T) {

	buff := new(bytes.Buffer)
	writer := multipart.NewWriter(buff)
	filewriter, err := writer.CreateFormFile("test", "test")
	if err != nil {
		t.Fatal(err)
	}

	_, err = filewriter.Write([]byte("Hello world"))
	if err != nil {
		t.Fatal(err)
	}

	err = writer.Close()
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "http://www.example.de", buff)
	req.Header.Set("content-type", writer.FormDataContentType())

	file, err := requestutil.ExtractUpload("test", req)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Content.Close()

	data, err := ioutil.ReadAll(file.Content)
	if err != nil {
		t.Fatal(err)
	}

	result := bytes.Compare(data, []byte("Hello world"))
	if result != 0 {
		t.Fatalf(`content of file buffer isn't equal (%d)`, result)
	}
}
