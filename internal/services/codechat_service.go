package services

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/sdrvirtual/codewoot/internal/codechat"
	"github.com/sdrvirtual/codewoot/internal/config"
	"github.com/sdrvirtual/codewoot/internal/dto"
)

type CodechatService struct {
	cfg           *config.Config
	client        *codechat.Client
}

func NewCodechatService(cfg *config.Config) *CodechatService {
	// TODO: Pegar da db
	instance := "codechat_v1"
	instanceToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpbnN0YW5jZU5hbWUiOiJjb2RlY2hhdF92MSIsImFwaU5hbWUiOiJ3aGF0c2FwcC1hcGkiLCJ0b2tlbklkIjoiMDFLQjNZMFlYVDY4QkRCOFo4UEtIRDY1WTkiLCJpYXQiOjE3NjQyODk5NjksImV4cCI6MTc2NDI4OTk2OSwic3ViIjoiZy10In0.HSVHOU_kCwOguJv-bLN23hpxibYuveXfPylq9DxITI4"
	codechatClient, err := codechat.New(cfg.Codechat.URL, cfg.Codechat.GlobalToken, codechat.WithInstanceToken(instanceToken, instance))

	if err != nil {
		log.Fatal(err)
	}
	return &CodechatService{
		cfg:           cfg,
		client:        codechatClient,
	}
}

func (c* CodechatService) SendTextMessage(toNumber string, text string) error {
	payload := map[string]any{
		"number": toNumber,
		"textMessage": map[string]any{
			"text": text,
		},
	}
	_, err := c.client.SendText(context.TODO(), payload)
	return err
}

func (c* CodechatService) GetAudioContent(ctx context.Context, message *dto.CodechatData) (*dto.FileData, error) {
	data, err := c.client.GetMediaData(ctx, message)
	if err != nil {
		return nil, err
	}

	file ,err := os.Create("./"+data.Name)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(file, data.File)

	if err != nil {
		return nil, err
	}
	return data, nil
}
