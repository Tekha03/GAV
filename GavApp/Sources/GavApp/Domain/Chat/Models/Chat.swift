import Foundation

public struct Chat: Identifiable, Codable {
    public let id: UUID
    public let isGroup: Bool
    public let title: String
    public let photoUrl: String
    public let createdAt: Date
}