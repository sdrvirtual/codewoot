package services

import (
	"context"
	"io"
	"log"
	"strings"

	"github.com/sdrvirtual/codewoot/internal/audio"
	"github.com/sdrvirtual/codewoot/internal/codechat"
	"github.com/sdrvirtual/codewoot/internal/config"
	"github.com/sdrvirtual/codewoot/internal/db"
	"github.com/sdrvirtual/codewoot/internal/domain"
	"github.com/sdrvirtual/codewoot/internal/dto"
)

type CodechatService struct {
	cfg    *config.Config
	client *codechat.Client
}

func NewCodechatService(cfg *config.Config, session db.CodechatSession) *CodechatService {
	instance := session.CodechatInstance
	instanceToken := session.CodechatInstcanceToken
	codechatClient, err := codechat.New(cfg.Codechat.URL, cfg.Codechat.GlobalToken, codechat.WithInstanceToken(instanceToken, instance))

	if err != nil {
		log.Fatal(err)
	}
	return &CodechatService{
		cfg:    cfg,
		client: codechatClient,
	}
}

type CodechatClientMessage struct {
	Text           string
	PhoneNumber    string
	AttachmentName *string
	MediaURL       *string
	MediaType      *string
	AudioFile      io.Reader
	AudioFileName  *string
	FileURL        *string
}

func NewCodechatClientMessage() CodechatClientMessage {
	return CodechatClientMessage{}
}

func (c *CodechatService) GetMediaContent(ctx context.Context, message *dto.CodechatData) (*dto.FileData, error) {
	data, err := c.client.GetMediaData(ctx, message)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *CodechatService) GetAudioContent(ctx context.Context, message *dto.CodechatData) (*dto.FileData, error) {
	data, err := c.client.GetMediaData(ctx, message)
	if err != nil {
		return nil, err
	}
	oggData, err := audio.TranscodeOggToMp3(data.File)
	if err != nil {
		return nil, err
	}
	data.File = oggData
	data.Mimetype = "audio/ogg"
	data.Name = strings.Split(data.Name, ".")[0] + ".ogg"
	return data, nil
}

func (c *CodechatService) TranscodeAudioFromURL(ctx context.Context, url string) (io.Reader, error) {
	return audio.DownloadAndTranscodeToOgg(ctx, url)
}

func (c *CodechatService) SendMessage(ctx context.Context, contact domain.ContactInfo, message CodechatClientMessage) error {
	if message.MediaURL != nil {
		mediaType := "image"
		if message.MediaType != nil && *message.MediaType != "" {
			mediaType = *message.MediaType
		}
		params := codechat.SendMediaParams{
			Number: contact.Phone,
			MediaMessage: codechat.CCMediaMessage{
				Media:     *message.MediaURL,
				Mediatype: mediaType,
				Caption:   message.Text,
			},
		}
		_, err := c.client.SendMedia(ctx, params)
		if err != nil {
			return err
		}
		return nil
	}
	if message.AudioFile != nil {
		fileName := "audio.ogg"
		if message.AudioFileName != nil {
			fileName = *message.AudioFileName
		}
		params := codechat.SendWhatsappAudioParams{
			Number:    contact.Phone,
			AudioFile: message.AudioFile,
			FileName:  fileName,
		}
		_, err := c.client.SendWhatsappAudio(ctx, params)
		if err != nil {
			return err
		}
		return nil
	}
	if message.FileURL != nil {
		params := codechat.SendMediaParams{
			Number: contact.Phone,
			MediaMessage: codechat.CCMediaMessage{
				Media:     *message.FileURL,
				FileName:  *message.AttachmentName,
				Mediatype: "document",
				Caption:   message.Text,
			},
		}
		_, err := c.client.SendMedia(ctx, params)
		if err != nil {
			return err
		}
		return nil
	}

	if message.Text != "" && message.MediaURL == nil {
		params := codechat.SendTextParams{
			Number:      contact.Phone,
			TextMessage: codechat.CCTextMessage{Text: message.Text},
		}
		_, err := c.client.SendText(ctx, params)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}
