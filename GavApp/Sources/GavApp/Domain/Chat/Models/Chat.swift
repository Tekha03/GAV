import Foundation

public struct Chat: Identifiable, Codable {
    public let id: UUID
    public let isGroup: Bool
    public let title: String
    public let photoUrl: String
    public let createdAt: Date

    public init(
        id: UUID,
        isGroup: Bool,
        title: String,
        photoUrl: String,
        createdAt: Date
    ) {
        self.id = id
        self.isGroup = isGroup
        self.title = title
        self.photoUrl = photoUrl
        self.createdAt = createdAt
    }
}