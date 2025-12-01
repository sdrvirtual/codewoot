package dto

import (
	"time"
)

type CodechatInstance struct {
	Id                 int       `json:"id"`
	Name               string    `json:"name"`
	Description        string    `json:"description"`
	ConnectionStatus   string    `json:"connectionStatus"`
	OwnerJid           string    `json:"ownerJid"`
	ProfilePicUrl      string    `json:"profilePicUrl"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
	ExternalAttributes string    `json:"externalAttributes"`
}

type CodechatTextMessage struct {
	Text string `json:"text"`
}

type CodechatData struct {
	Id               int                 `json:"id"`
	KeyId            string              `json:"keyId"`
	KeyRemoteJid     string              `json:"KeyRemoteJid"`
	KeyFromMe        bool                `json:"keyFromMe"`
	PushName         string              `json:"pushName"`
	MessageType      string              `json:"messageType"`
	Content          CodechatTextMessage `json:"content"`
	MessageTimestamp int                 `json:"messageTimestamp"`
	InstanceId       int                 `json:"instanceId"`
	Device           string              `json:"device"`
	IsGroup          bool                `json:"isGroup"`
}

type CodechatWebhook struct {
	Event    string           `json:"event"`
	Instance CodechatInstance `json:"instance"`
	Data     CodechatData     `json:"data"`
}
