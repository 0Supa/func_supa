package fun

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
)

type Upload struct {
	ID       string `json:"id"`
	Ext      string `json:"ext"`
	Type     string `json:"type"`
	Checksum string `json:"checksum"`
	Key      string `json:"key"`
	Link     string `json:"link"`
	Delete   string `json:"delete"`
}

func UploadFile(body io.ReadCloser, fileName string, contentType string) (upload Upload, err error) {
	fileBuf := &bytes.Buffer{}
	writer := multipart.NewWriter(fileBuf)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, QuoteEscaper.Replace(fileName)))
	h.Set("Content-Type", contentType)

	defer body.Close()
	part, err := writer.CreatePart(h)
	if _, err := io.Copy(part, body); err != nil {
		return upload, err
	}
	writer.Close()

	res, err := apiClient.Post("https://kappa.lol/api/upload?skip-cd=true", writer.FormDataContentType(), fileBuf)
	if err != nil {
		return
	}
	defer res.Body.Close()

	buf, _ := io.ReadAll(res.Body)
	err = json.Unmarshal(buf, &upload)

	if res.StatusCode != http.StatusOK {
		return upload, errors.New(string(buf))
	}

	return
}
