package chat

type SendMessageInput struct {
    ChatID      uint
    SenderID    uint
    Text        *string
    Attachments []AttachmentInput
}

type AttachmentInput struct {
    Type     AttachmentType
    URL      string
    FileName string
    FileSize int64
}
