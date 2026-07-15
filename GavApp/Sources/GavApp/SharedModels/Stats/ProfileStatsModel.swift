import Foundation

public struct ProfileStatsModel: Codable, Equatable {
    public let userId: UUID
    public let postCount: UInt
    public let followersCount: UInt
    public let followingsCount: UInt
}

extension ProfileStatsModel {
    private enum CodingKeys: String, CodingKey {
        case userId = "user_id"
        case postCount = "post_count"
        case followersCount = "followers_count"
        case followingsCount = "followings_count"
    }
}