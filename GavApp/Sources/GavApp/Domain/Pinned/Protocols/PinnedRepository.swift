import Foundation

public protocol PinnedMessageRepository {
    func pinMessage(messageID: UUID) async throws
    func unpinMessage(messageID: UUID) async throws
    func getPinnedMessages(chatID: UUID) async throws -> [PinnedMessage]
}