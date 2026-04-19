import Foundation

public protocol ChatUseCase {
    func createPrivateChat(user1: UUID, user2: UUID) async throws -> Chat
    func createGroupChat(
        title: String,
        creator: UUID,
        members: [UUID]
    ) async throws -> Chat

    func getChatByID(id: UUID) async throws -> Chat
    func getUserChats(userID: UUID) async throws -> [Chat]

    func updateChatTitle(chatID: UUID, title: String) async throws
    func updateChatPhoto(chatID: UUID, photoUrl: String) async throws

    func leaveChat(userID: UUID, chatID: UUID) async throws
    func deleteChat(chatID: UUID) async throws

    func getChatMembers(chatID: UUID) async throws -> [ChatMember]
    func addMember(userID: UUID, chatID: UUID) async throws
    func removeMember(userID: UUID, chatID: UUID) async throws

    func getMessages(chatID: UUID, limit: Int, before: UUID?) async throws -> [Message]

    func sendMessage(
        chatID: UUID,
        text: String?,
        attachments: [AttachmentInput]?,
        replyToId: UUID?
    ) async throws -> Message

    func markAsRead(chatID: UUID, userID: UUID) async throws

    func sendTyping(chatID: UUID) async throws
}