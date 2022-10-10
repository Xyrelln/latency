// resource update

package upload

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

type content struct {
	fname string
	ftype string
	fdata []byte
}

func SendPostRequest(url string, files ...content) ([]byte, error) {
	var (
		buf = new(bytes.Buffer)
		w   = multipart.NewWriter(buf)
	)
	for _, f := range files {
		part, err := w.CreateFormFile(f.ftype, filepath.Base(f.fname))
		if err != nil {
			return []byte{}, err
		}
		part.Write(f.fdata)
	}

	w.Close()

	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return []byte{}, err
	}
	req.Header.Add("Content-Type", w.FormDataContentType())

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer res.Body.Close()

	cnt, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}
	return cnt, nil
}
