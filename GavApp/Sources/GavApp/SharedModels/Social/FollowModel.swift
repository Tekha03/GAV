import Foundation

public struct FollowModel: Codable, Equatable, Hashable {
    public let followerId: UUID
    public let followingId: UUID
}

extension FollowModel {
    private enum CodingKeys: String, CodingKey {
        case followerId = "follower_id"
        case followingId = "following_id"
    }
}