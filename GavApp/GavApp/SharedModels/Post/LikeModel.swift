import Foundation

public struct LikeModel: Codable, Equatable, Hashable {
    public let userId: UUID
    public let postId: UUID
}

extension LikeModel {
    private enum CodingKeys: String, CodingKey {
        case userId = "user_id"
        case postId = "post_id"
    }
}