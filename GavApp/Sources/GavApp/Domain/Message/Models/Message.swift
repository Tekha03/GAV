import Foundation

public struct Message: Identifiable, Codable {
    public let id: UUID
    public let chatId: UUID
    public let senderId: UUID
    public let text: String? // может быть nil, если attachment
    public let replyToId: UUID?
    public let createdAt: Date
    public var editedAt: Date?

    public var attachments: [Attachment] = []
    public var reactions: [Reaction] = []
}