import Foundation

public protocol MessageRepository {
    func sendMessage(
        chatId: UUID,
        text: String?,
        attachments: [AttachmentInput]?,
        replyToId: UUID?
    ) async throws -> Message

    func editMessage(
        messageID: UUID,
        text: String
    ) async throws

    func deleteMessage(messageID: UUID) async throws

    func getMessages(
        chatID: UUID,
        limit: Int,
        before: UUID?
    ) async throws -> [Message]

    func markAsRead(chatID: UUID, userID: UUID) async throws

    func getUnreadCount(userID: UUID) async throws -> Int
    func getChatUnreadCount(chatID: UUID, userID: UUID) async throws -> Int
}