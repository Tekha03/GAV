import Foundation

public struct PostModel: Codable, Equatable {
    public let id: UUID
    public let userId: UUID
    public let content: String
    public let imageUrl: String?
    public let createdAt: Date
}

extension PostModel {
    private enum CodingKeys: String, CodingKey {
        case id
        case userId = "user_id"
        case content
        case imageUrl = "image_url"
        case createdAt = "created_at"
    }
}