import Foundation

public struct UserStatsModel: Codable, Equatable {
    public let id: UUID
    public let userId: UUID
    public let postCount: UInt
    public let followersCount: UInt
    public let followingsCount: UInt
    public let dogsCount: UInt
    public let createdAt: Date
    public let updatedAt: Date
}

extension UserStatsModel {
    private enum CodingKeys: String, CodingKey {
        case id
        case userId = "user_id"
        case postCount = "post_count"
        case followersCount = "followers_count"
        case followingsCount = "followings_count"
        case dogsCount = "dogs_count"
        case createdAt = "created_at"
        case updatedAt = "updated_at"
    }
}