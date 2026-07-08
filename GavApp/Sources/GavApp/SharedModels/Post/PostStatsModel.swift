import Foundation

public struct PostStatsModel: Codable, Equatable {
    public let id: UUID
    public let postId: UUID
    public let likesCount: UInt
    public let commentsCount: UInt
    public let createdAt: Date
    public let updatedAt: Date

    public init(
        id: UUID,
        postId: UUID,
        likesCount: UInt,
        commentsCount: UInt,
        createdAt: Date,
        updatedAt: Date
    ) {
        self.id = id
        self.postId = postId
        self.likesCount = likesCount
        self.commentsCount = commentsCount
        self.createdAt = createdAt
        self.updatedAt = updatedAt
    }

    public init(from decoder: Decoder) throws {
        let container = try decoder.container(keyedBy: CodingKeys.self)
        postId = try container.decode(UUID.self, forKey: .postId)
        id = try container.decodeIfPresent(UUID.self, forKey: .id) ?? postId
        likesCount = try container.decodeIfPresent(UInt.self, forKey: .likesCount) ?? 0
        commentsCount = try container.decodeIfPresent(UInt.self, forKey: .commentsCount) ?? 0
        createdAt = Self.decodeDateIfPresent(container, forKey: .createdAt) ?? .distantPast
        updatedAt = Self.decodeDateIfPresent(container, forKey: .updatedAt) ?? createdAt
    }

    private static func decodeDateIfPresent(
        _ container: KeyedDecodingContainer<CodingKeys>,
        forKey key: CodingKeys
    ) -> Date? {
        if let date = try? container.decodeIfPresent(Date.self, forKey: key) {
            return date
        }

        guard let value = try? container.decodeIfPresent(String.self, forKey: key) else {
            return nil
        }

        let withFraction = ISO8601DateFormatter()
        withFraction.formatOptions = [.withInternetDateTime, .withFractionalSeconds]
        if let date = withFraction.date(from: value) {
            return date
        }

        return ISO8601DateFormatter().date(from: value)
    }
}

extension PostStatsModel {
    private enum CodingKeys: String, CodingKey {
        case id
        case postId = "post_id"
        case likesCount = "likes_count"
        case commentsCount = "comments_count"
        case createdAt = "created_at"
        case updatedAt = "updated_at"
    }
}
