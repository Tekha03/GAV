package chat

type AttachmentType string

var (
	AttachmentImage AttachmentType = "image"
	AttachmentVideo AttachmentType = "video"
	AttachmentVoice AttachmentType = "voice"
	AttachmentFile  AttachmentType = "file"
)

type Attachment struct {
	ID        uint
	MessageID uint
	URL       string
	Type      AttachmentType
	FileName  string
	FileSize  string
	MimeSize  string
}
