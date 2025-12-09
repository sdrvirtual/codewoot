package dto

type ChatwootWebhook struct {
	Account              CWAccount      `json:"account"`
	AdditionalAttributes map[string]any `json:"additional_attributes"`
	ContentAttributes    map[string]any `json:"content_attributes"`
	ContentType          string         `json:"content_type"`
	Content              string         `json:"content"`
	Conversation         CWConversation `json:"conversation"`
	CreatedAt            any            `json:"created_at"`
	ID                   int            `json:"id"`
	Inbox                CWInbox        `json:"inbox"`
	MessageType          CWMessageType  `json:"message_type"`
	Private              bool           `json:"private"`
	Sender               CWSimpleSender `json:"sender"`
	SourceID             *string        `json:"source_id"`
	Event                string         `json:"event"`
}

type CWMessageType string

const (
	Outgoing CWMessageType = "outgoing"
	Incoming CWMessageType = "incoming"
)

type CWAccount struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type CWInbox struct {
	ID          int     `json:"id"`
	AvatarUrl   string  `json:"avatar_url"`
	ChannelID   int     `json:"channel_id"`
	Name        string  `json:"name"`
	ChannelType string  `json:"channel_type"`
	Provider    *string `json:"provider"`
}

type CWSimpleSender struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Email *string `json:"email"`
	Type  string  `json:"type"`
}

type CWConversation struct {
	AdditionalAttributes map[string]any        `json:"additional_attributes"`
	CanReply             bool                  `json:"can_reply"`
	Channel              string                `json:"channel"`
	ContactInbox         CWWebhookContactInbox `json:"contact_inbox"`
	ID                   int                   `json:"id"`
	InboxID              int                   `json:"inbox_id"`
	Messages             []CWMessage           `json:"messages"`
	Labels               []string              `json:"labels"`
	Meta                 CWMeta                `json:"meta"`
	Status               string                `json:"status"`
	CustomAttributes     map[string]any        `json:"custom_attributes"`
	SnoozedUntil         any                   `json:"snoozed_until"`
	UnreadCount          int                   `json:"unread_count"`
	FirstReplyCreatedAt  any                   `json:"first_reply_created_at"`
	Priority             any                   `json:"priority"`
	WaitingSince         float64               `json:"waiting_since"`
	AgentLastSeenAt      float64               `json:"agent_last_seen_at"`
	ContactLastSeenAt    float64               `json:"contact_last_seen_at"`
	LastActivityAt       float64               `json:"last_activity_at"`
	Timestamp            float64               `json:"timestamp"`
	CreatedAt            float64               `json:"created_at"`
	UpdatedAt            float64               `json:"updated_at"`
}

type CWContactInbox struct {
	SourceID string  `json:"source_id"`
	Inbox    CWInbox `json:"inbox"`
}

type CWContact struct {
	ID                   int              `json:"id"`
	AdditionalAttributes map[string]any   `json:"additional_attributes"`
	AvailabilityStatus   *string          `json:"availability_status"`
	Email                string           `json:"email"`
	Name                 string           `json:"name"`
	PhoneNumber          string           `json:"phone_number"`
	Blocked              bool             `json:"blocked"`
	Identifier           *string          `json:"identifier"`
	Thumbnail            string           `json:"thumbnail"`
	CustomAttributes     map[string]any   `json:"custom_attributes"`
	CreatedAt            float64          `json:"created_at"`
	ContactInboxes       []CWContactInbox `json:"contact_inboxes"`
}

type CWWebhookContactInbox struct {
	ID           int    `json:"id"`
	ContactID    int    `json:"contact_id"`
	InboxID      int    `json:"inbox_id"`
	SourceID     string `json:"source_id"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	HmacVerified bool   `json:"hmac_verified"`
	PubsubToken  string `json:"pubsub_token"`
}

// TEXT MESSAGE
//     "messages": [
//       {
//         "id": 98,
//         "content": "opa\n",
//         "account_id": 2,
//         "inbox_id": 1,
//         "conversation_id": 12,
//         "message_type": 1,
//         "created_at": 1764803328,
//         "updated_at": "2025-12-03T23:08:48.069Z",
//         "private": false,
//         "status": "sent",
//         "source_id": null,
//         "content_type": "text",
//         "content_attributes": {},
//         "sender_type": "User",
//         "sender_id": 2,
//         "external_source_ids": {},
//         "additional_attributes": {},
//         "processed_message_content": "opa\n",
//         "sentiment": {},
//         "conversation": {
//           "assignee_id": 2,
//           "unread_count": 0,
//           "last_activity_at": 1764803328,
//           "contact_inbox": {
//             "source_id": "5a2148ac-6ede-4d2e-9f33-8d12b6e8ce1c"
//           }
//         },
//         "sender": {
//           "id": 2,
//           "name": "Cássio",
//           "available_name": "Cássio",
//           "avatar_url": "",
//           "type": "user",
//           "availability_status": null,
//           "thumbnail": ""
//         }
//       }
//     ],

// AUDIO MESSAGE
//     "messages": [
//       {
//         "id": 99,
//         "content": null,
//         "account_id": 2,
//         "inbox_id": 1,
//         "conversation_id": 12,
//         "message_type": 1,
//         "created_at": 1764803563,
//         "updated_at": "2025-12-03T23:12:43.839Z",
//         "private": false,
//         "status": "sent",
//         "source_id": null,
//         "content_type": "text",
//         "content_attributes": {},
//         "sender_type": "User",
//         "sender_id": 2,
//         "external_source_ids": {},
//         "additional_attributes": {},
//         "processed_message_content": null,
//         "sentiment": {},
//         "conversation": {
//           "assignee_id": 2,
//           "unread_count": 0,
//           "last_activity_at": 1764803563,
//           "contact_inbox": {
//             "source_id": "5a2148ac-6ede-4d2e-9f33-8d12b6e8ce1c"
//           }
//         },
//         "attachments": [
//           {
//             "id": 15,
//             "message_id": 99,
//             "file_type": "audio",
//             "account_id": 2,
//             "extension": null,
//             "data_url": "http://localhost:3000/rails/active_storage/blobs/redirect/eyJfcmFpbHMiOnsibWVzc2FnZSI6IkJBaHBGQT09IiwiZXhwIjpudWxsLCJwdXIiOiJibG9iX2lkIn19--af4235df8dd0a4a4bbce663ede4661bd17d7f969/68ebb6434341.mp3",
//             "thumb_url": "",
//             "file_size": 52608,
//             "width": null,
//             "height": null,
//             "transcribed_text": ""
//           }
//         ],
//         "sender": {
//           "id": 2,
//           "name": "Cássio",
//           "available_name": "Cássio",
//           "avatar_url": "",
//           "type": "user",
//           "availability_status": null,
//           "thumbnail": ""
//         }
//       }
//     ],

// MEDIA MESSAGE
//     "messages": [
//       {
//         "id": 100,
//         "content": "oiee",
//         "account_id": 2,
//         "inbox_id": 1,
//         "conversation_id": 12,
//         "message_type": 1,
//         "created_at": 1764804115,
//         "updated_at": "2025-12-03T23:21:55.704Z",
//         "private": false,
//         "status": "sent",
//         "source_id": null,
//         "content_type": "text",
//         "content_attributes": {},
//         "sender_type": "User",
//         "sender_id": 2,
//         "external_source_ids": {},
//         "additional_attributes": {},
//         "processed_message_content": "oiee",
//         "sentiment": {},
//         "conversation": {
//           "assignee_id": 2,
//           "unread_count": 0,
//           "last_activity_at": 1764804115,
//           "contact_inbox": {
//             "source_id": "5a2148ac-6ede-4d2e-9f33-8d12b6e8ce1c"
//           }
//         },
//         "attachments": [
//           {
//             "id": 16,
//             "message_id": 100,
//             "file_type": "image",
//             "account_id": 2,
//             "extension": null,
//             "data_url": "http://localhost:3000/rails/active_storage/blobs/redirect/eyJfcmFpbHMiOnsibWVzc2FnZSI6IkJBaHBGUT09IiwiZXhwIjpudWxsLCJwdXIiOiJibG9iX2lkIn19--b44af82d04b087b545436b1c92e3fff2bbcf5841/05.png",
//             "thumb_url": "http://localhost:3000/rails/active_storage/representations/redirect/eyJfcmFpbHMiOnsibWVzc2FnZSI6IkJBaHBGUT09IiwiZXhwIjpudWxsLCJwdXIiOiJibG9iX2lkIn19--b44af82d04b087b545436b1c92e3fff2bbcf5841/eyJfcmFpbHMiOnsibWVzc2FnZSI6IkJBaDdCem9MWm05eWJXRjBTU0lJY0c1bkJqb0dSVlE2RTNKbGMybDZaVjkwYjE5bWFXeHNXd2RwQWZvdyIsImV4cCI6bnVsbCwicHVyIjoidmFyaWF0aW9uIn19--f1b26a23780bc9bc51fbe09fba8e4c2bd19fa942/05.png",
//             "file_size": 3206366,
//             "width": null,
//             "height": null
//           }
//         ],
//         "sender": {
//           "id": 2,
//           "name": "Cássio",
//           "available_name": "Cássio",
//           "avatar_url": "",
//           "type": "user",
//           "availability_status": null,
//           "thumbnail": ""
//         }
//       }
//     ],

type CWAttachment struct {
	ID        int     `json:"id"`
	MessageID int     `json:"message_id"`
	FileType  string  `json:"file_type"`
	AccountID int     `json:"account_id"`
	Extension *string `json:"extension"`
	DataURL   *string `json:"data_url"`
	ThumbURL  *string `json:"thumb_url"`
	FileSize  int     `json:"file_size"`
	Width     int     `json:"width"`
	Height    int     `json:"height"`
}

type CWMessage struct {
	ID                   int               `json:"id"`
	Content              *string           `json:"content"`
	AccountID            int               `json:"account_id"`
	InboxID              int               `json:"inbox_id"`
	ConversationID       int               `json:"conversation_id"`
	MessageType          int               `json:"message_type"`
	CreatedAt            float64           `json:"created_at"`
	UpdatedAt            string            `json:"updated_at"`
	Private              bool              `json:"private"`
	Status               string            `json:"status"`
	SourceID             *string           `json:"source_id"`
	ContentType          string            `json:"content_type"`
	ContentAttributes    map[string]any    `json:"content_attributes"`
	SenderType           string            `json:"sender_type"`
	SenderID             int               `json:"sender_id"`
	ExternalSourceIDs    map[string]any    `json:"external_source_ids"`
	AdditionalAttributes map[string]any    `json:"additional_attributes"`
	ProcessedMessage     string            `json:"processed_message_content"`
	Sentiment            map[string]any    `json:"sentiment"`
	Conversation         CWMsgConversation `json:"conversation"`
	Sender               CWMessageSender   `json:"sender"`
	Attachments          []CWAttachment    `json:"attachments"`
}

type CWMsgConversation struct {
	AssigneeID     int     `json:"assignee_id"`
	UnreadCount    int     `json:"unread_count"`
	LastActivityAt float64 `json:"last_activity_at"`
	ContactInbox   struct {
		SourceID string `json:"source_id"`
	} `json:"contact_inbox"`
}

type CWMessageSender struct {
	ID                 int     `json:"id"`
	Name               string  `json:"name"`
	AvailableName      string  `json:"available_name"`
	AvatarURL          string  `json:"avatar_url"`
	Type               string  `json:"type"`
	AvailabilityStatus *string `json:"availability_status"`
	Thumbnail          string  `json:"thumbnail"`
}

type CWMeta struct {
	Sender       CWSenderMeta `json:"sender"`
	Assignee     CWSenderMeta `json:"assignee"`
	Team         any          `json:"team"`
	HmacVerified bool         `json:"hmac_verified"`
}

type CWSenderMeta struct {
	AdditionalAttributes map[string]any `json:"additional_attributes"`
	CustomAttributes     map[string]any `json:"custom_attributes"`
	Email                *string        `json:"email"`
	ID                   int            `json:"id"`
	Identifier           *string        `json:"identifier"`
	Name                 string         `json:"name"`
	PhoneNumber          string         `json:"phone_number"`
	Thumbnail            string         `json:"thumbnail"`
	Blocked              bool           `json:"blocked"`
	Type                 string         `json:"type"`
}
