import Foundation

public struct CommentModel: Codable, Equatable {
    public let id: UUID
    public let postId: UUID
    public let userId: UUID
    public let content: String
    public let createdAt: Date
}

extension CommentModel {
    private enum CodingKeys: String, CodingKey {
        case id
        case postId = "post_id"
        case userId = "user_id"
        case content
        case createdAt = "created_at"
    }
}