package codechat

import (
	"context"
	"fmt"
	"mime"
	"net/http"
	"net/url"

	"github.com/sdrvirtual/codewoot/internal/dto"
)

func (c *Client) GetMediaData(ctx context.Context, message *dto.CodechatData) (*dto.FileData, error) {
	if c.instance == "" {
		return nil, fmt.Errorf("instanceName is required")
	}
	p := fmt.Sprintf("/chat/mediaData/%s", url.PathEscape(c.instance))
	req, err := c.newRequest(ctx, http.MethodPost, p, message)

	q := req.URL.Query()
	q.Set("binary", "true")
	req.URL.RawQuery = q.Encode()

	if err != nil {
		return nil, err
	}
	_, r, err := c.do(req)
	if err != nil {
		return nil, err
	}

	_, params, err := mime.ParseMediaType(r.Header.Get("Content-Disposition"))
	if err != nil {
		return nil, err
	}

	fileData := dto.FileData{
		Name: params["filename"],
		Mimetype: r.Header.Get("Content-Type"),
		File: r.Body,
	}

	return &fileData, nil
}
