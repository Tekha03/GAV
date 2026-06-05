import Foundation

public struct Reaction: Identifiable, Codable {
    public let id: UUID
    public let messageID: UUID
    public let userID: UUID
    public let emoji: String
}