import Foundation

public struct PinnedMessage: Codable {
    public let chatID: UUID
    public let messageID: UUID
    public let pinnedAt: Date
}