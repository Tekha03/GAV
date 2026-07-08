import Foundation

public struct PostModel: Codable, Equatable {
    public let id: UUID
    public let userId: UUID
    public let content: String
    public let imageUrl: String?
    public let createdAt: Date

    public init(id: UUID, userId: UUID, content: String, imageUrl: String?, createdAt: Date) {
        self.id = id
        self.userId = userId
        self.content = content
        self.imageUrl = imageUrl
        self.createdAt = createdAt
    }

    public init(from decoder: Decoder) throws {
        let container = try decoder.container(keyedBy: CodingKeys.self)
        id = try container.decode(UUID.self, forKey: .id)
        userId = try container.decodeIfPresent(UUID.self, forKey: .userId)
            ?? container.decode(UUID.self, forKey: .authorId)
        content = try container.decode(String.self, forKey: .content)
        imageUrl = try container.decodeIfPresent(String.self, forKey: .imageUrl)
        createdAt = try container.decode(Date.self, forKey: .createdAt)
    }

    public func encode(to encoder: Encoder) throws {
        var container = encoder.container(keyedBy: CodingKeys.self)
        try container.encode(id, forKey: .id)
        try container.encode(userId, forKey: .userId)
        try container.encode(content, forKey: .content)
        try container.encodeIfPresent(imageUrl, forKey: .imageUrl)
        try container.encode(createdAt, forKey: .createdAt)
    }
}

extension PostModel {
    private enum CodingKeys: String, CodingKey {
        case id
        case userId = "user_id"
        case authorId = "author_id"
        case content
        case imageUrl = "image_url"
        case createdAt = "created_at"
    }
}
