// Package dto
package dto

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type CodechatBufferString string

func (b *CodechatBufferString) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		*b = ""
		return nil
	}
	if data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		*b = CodechatBufferString(s)
		return nil
	}
	var obj struct {
		Type string `json:"type"`
		Data []byte `json:"data"`
	}
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	if obj.Type != "" && obj.Type != "Buffer" {
		return fmt.Errorf("unexpected buffer type: %s", obj.Type)
	}
	if len(obj.Data) == 0 {
		*b = ""
		return nil
	}
	*b = CodechatBufferString(base64.StdEncoding.EncodeToString(obj.Data))
	return nil
}

type CodechatLongString string

func (l *CodechatLongString) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		*l = ""
		return nil
	}
	if data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		*l = CodechatLongString(s)
		return nil
	}
	if (data[0] >= '0' && data[0] <= '9') || data[0] == '-' {
		var n json.Number
		if err := json.Unmarshal(data, &n); err != nil {
			return err
		}
		*l = CodechatLongString(n.String())
		return nil
	}
	var obj struct {
		Low      uint32 `json:"low"`
		High     uint32 `json:"high"`
		Unsigned bool   `json:"unsigned"`
	}
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	if obj.Unsigned {
		v := (uint64(obj.High) << 32) | uint64(obj.Low)
		*l = CodechatLongString(strconv.FormatUint(v, 10))
		return nil
	}
	v := (int64(int32(obj.High)) << 32) | int64(obj.Low)
	*l = CodechatLongString(strconv.FormatInt(v, 10))
	return nil
}

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
	Caption           string               `json:"caption"`
	DirectPath        string               `json:"directPath"`
	FileEncSha256     CodechatBufferString `json:"fileEncSha256"`
	FileLength        CodechatLongString   `json:"fileLength"`
	FileSha256        CodechatBufferString `json:"fileSha256"`
	Height            int                  `json:"height"`
	JpegThumbnail     CodechatBufferString `json:"jpegThumbnail"`
	MediaKey          CodechatBufferString `json:"mediaKey"`
	MediaKeyTimestamp CodechatLongString   `json:"mediaKeyTimestamp"`
	MimeType          string               `json:"mimetype"`
	URL               string               `json:"url"`
	ViewOnce          bool                 `json:"viewOnce"`
	Width             int                  `json:"width"`
}

func (CodechatImageContent) isCodechatMessageContent() {}

type CodechatAudioContent struct {
	DirectPath        string               `json:"directPath"`
	FileEncSha256     CodechatBufferString `json:"fileEncSha256"`
	FileLength        CodechatLongString   `json:"fileLength"`
	FileSha256        CodechatBufferString `json:"fileSha256"`
	MediaKey          CodechatBufferString `json:"mediaKey"`
	MediaKeyTimestamp CodechatLongString   `json:"mediaKeyTimestamp"`
	MimeType          string               `json:"mimetype"`
	Ptt               bool                 `json:"ptt"`
	Seconds           int                  `json:"seconds"`
	URL               string               `json:"url"`
	ViewOnce          bool                 `json:"viewOnce"`
	Waveform          CodechatBufferString `json:"waveform"`
}

func (CodechatAudioContent) isCodechatMessageContent() {}

type CodechatDocumentContent struct {
	Title             string               `json:"title"`
	Caption           string               `json:"caption"`
	DirectPath        string               `json:"directPath"`
	FileEncSha256     CodechatBufferString `json:"fileEncSha256"`
	FileLength        CodechatLongString   `json:"fileLength"`
	FileSha256        CodechatBufferString `json:"fileSha256"`
	Height            int                  `json:"height"`
	JpegThumbnail     CodechatBufferString `json:"jpegThumbnail"`
	MediaKey          CodechatBufferString `json:"mediaKey"`
	MediaKeyTimestamp CodechatLongString   `json:"mediaKeyTimestamp"`
	MimeType          string               `json:"mimetype"`
	URL               string               `json:"url"`
	ViewOnce          bool                 `json:"viewOnce"`
	Width             int                  `json:"width"`
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
	case "documentMessage", "documentWithCaptionMessage":
		var raw struct {
			Message struct {
				DocumentMessage CodechatDocumentContent `json:"documentMessage"`
			} `json:"message"`
		}

		if err := json.Unmarshal(aux.Content, &raw); err != nil {
			return err
		}
		c.Content = raw.Message.DocumentMessage
	case "conversation", "extendedTextMessage":
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
