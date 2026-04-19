import Foundation

public struct PostStatsModel: Codable, Equatable {
    public let id: UUID
    public let postId: UUID
    public let likesCount: UInt
    public let commentsCount: UInt
    public let createdAt: Date
    public let updatedAt: Date
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