import Foundation

public protocol ChatRepository {
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
}