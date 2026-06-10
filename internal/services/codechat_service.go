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
	"github.com/sdrvirtual/codewoot/internal/media"
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
	MediaFile      io.Reader
	MediaFileName  *string
	MediaMimeType  *string
	AudioFile      io.Reader
	AudioFileName  *string
	FileURL        *string
}

type CodechatRoute struct {
	Number string
	JID    string
	LID    string
	Exists bool
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

func (c *CodechatService) PrepareImageFromURL(ctx context.Context, url string) (*media.PreparedMedia, error) {
	return media.PrepareImageFromURL(ctx, url)
}

func (c *CodechatService) PrepareVideoFromURL(ctx context.Context, url string) (*media.PreparedMedia, error) {
	return media.PrepareVideoFromURL(ctx, url)
}

func (c *CodechatService) ResolveRoutingNumber(ctx context.Context, phone string) (CodechatRoute, error) {
	route := CodechatRoute{Number: phone}
	numbers, err := c.client.WhatsappNumbers(ctx, []string{phone})
	if err != nil {
		return route, err
	}
	if len(numbers) == 0 {
		return route, nil
	}

	resolved := numbers[0]
	route.JID = strings.TrimSpace(resolved.JID)
	route.LID = strings.TrimSpace(resolved.LID)
	route.Exists = resolved.Exists
	if isLidJID(route.LID) {
		route.Number = route.LID
	}
	return route, nil
}

func (c *CodechatService) SendMessage(ctx context.Context, contact domain.ContactInfo, message CodechatClientMessage) error {
	number := contact.Phone
	if contact.RoutingJID != "" {
		number = contact.RoutingJID
	}

	if message.MediaFile != nil {
		if closer, ok := message.MediaFile.(io.Closer); ok {
			defer closer.Close()
		}

		mediaType := "image"
		if message.MediaType != nil && *message.MediaType != "" {
			mediaType = *message.MediaType
		}
		fileName := "media"
		if message.MediaFileName != nil && *message.MediaFileName != "" {
			fileName = *message.MediaFileName
		}
		mimeType := "application/octet-stream"
		if message.MediaMimeType != nil && *message.MediaMimeType != "" {
			mimeType = *message.MediaMimeType
		}
		params := codechat.SendMediaFileParams{
			Number:    number,
			MediaType: mediaType,
			Caption:   message.Text,
			File:      message.MediaFile,
			FileName:  fileName,
			MimeType:  mimeType,
		}
		_, err := c.client.SendMediaFile(ctx, params)
		if err != nil {
			return err
		}
		return nil
	}

	if message.MediaURL != nil {
		mediaType := "image"
		if message.MediaType != nil && *message.MediaType != "" {
			mediaType = *message.MediaType
		}
		params := codechat.SendMediaParams{
			Number: number,
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
			Number:    number,
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
			Number: number,
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
			Number:      number,
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
