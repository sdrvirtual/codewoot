// Package dto
package dto

import (
	"encoding/json"
	"fmt"
	"time"
)

type CodechatInstance struct {
	ID                 int       `json:"id"`
	Name               string    `json:"name"`
	Description        string    `json:"description"`
	ConnectionStatus   string    `json:"connectionStatus"`
	OwnerJid           string    `json:"ownerJid"`
	ProfilePicURL      string    `json:"profilePicUrl"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
	ExternalAttributes string    `json:"externalAttributes"`
}

type CodechatMessageContent interface {
	isCodechatMessageContent()
}

type CodechatTextContent struct {
	Text string `json:"text"`
}

func (CodechatTextContent) isCodechatMessageContent() {}

type CodechatImageContent struct {
	Caption           string `json:"caption"`
	DirectPath        string `json:"directPath"`
	FileEncSha256     string `json:"fileEncSha256"`
	FileLength        string `json:"fileLength"`
	FileSha256        string `json:"fileSha256"`
	Height            int    `json:"height"`
	JpegThumbnail     string `json:"jpegThumbnail"`
	MediaKey          string `json:"mediaKey"`
	MediaKeyTimestamp string `json:"mediaKeyTimestamp"`
	MimeType          string `json:"mimetype"`
	URL               string `json:"url"`
	ViewOnce          bool   `json:"viewOnce"`
	Width             int    `json:"width"`
}

func (CodechatImageContent) isCodechatMessageContent() {}

type CodechatAudioContent struct {
	DirectPath        string `json:"directPath"`
	FileEncSha256     string `json:"fileEncSha256"`
	FileLength        string `json:"fileLength"`
	FileSha256        string `json:"fileSha256"`
	MediaKey          string `json:"mediaKey"`
	MediaKeyTimestamp string `json:"mediaKeyTimestamp"`
	MimeType          string `json:"mimetype"`
	Ptt               bool   `json:"ptt"`
	Seconds           int    `json:"seconds"`
	URL               string `json:"url"`
	ViewOnce          bool   `json:"viewOnce"`
	Waveform          string `json:"waveform"`
}

func (CodechatAudioContent) isCodechatMessageContent() {}

type CodechatDocumentContent struct {
	// TODO
}

func (CodechatDocumentContent) isCodechatMessageContent() {}

type CodechatData struct {
	ID               int                    `json:"id"`
	KeyID            string                 `json:"keyId"`
	KeyRemoteJid     string                 `json:"KeyRemoteJid"`
	KeyFromMe        bool                   `json:"keyFromMe"`
	PushName         string                 `json:"pushName"`
	MessageType      string                 `json:"messageType"`
	Content          CodechatMessageContent `json:"content"`
	MessageTimestamp int                    `json:"messageTimestamp"`
	InstanceID       int                    `json:"instanceId"`
	Device           string                 `json:"device"`
	IsGroup          bool                   `json:"isGroup"`
}

type CodechatWebhook struct {
	Event    string           `json:"event"`
	Instance CodechatInstance `json:"instance"`
	Data     CodechatData     `json:"data"`
}

func (c *CodechatData) UnmarshalJSON(data []byte) error {
	type Alias CodechatData
	aux := &struct {
		Content json.RawMessage `json:"content"`
		*Alias
	}{
		Alias: (*Alias)(c),
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	switch c.MessageType {
	case "protocolMessage":
		// TODO: handle this
	case "documentMessage":
		var msg CodechatDocumentContent
		if err := json.Unmarshal(aux.Content, &msg); err != nil {
			return err
		}
		c.Content = msg
	case "conversation":
		var msg CodechatTextContent
		if err := json.Unmarshal(aux.Content, &msg); err != nil {
			return err
		}
		c.Content = msg
	case "extendedTextMessage":
		var msg CodechatTextContent
		if err := json.Unmarshal(aux.Content, &msg); err != nil {
			return err
		}
		c.Content = msg

	case "audioMessage":
		var msg CodechatAudioContent
		if err := json.Unmarshal(aux.Content, &msg); err != nil {
			return err
		}
		c.Content = msg

	case "imageMessage":
		var msg CodechatImageContent
		if err := json.Unmarshal(aux.Content, &msg); err != nil {
			return err
		}

		c.Content = msg

	default:
		return fmt.Errorf("unknown message type: %s", c.MessageType)
	}

	return nil
}
