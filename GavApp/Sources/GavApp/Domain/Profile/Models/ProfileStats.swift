import Foundation

public struct ProfileStats: Equatable, Sendable {
    public let userId: UUID
    public let postCount: UInt
    public let followersCount: UInt
    public let followingsCount: UInt

    public init(
        userId: UUID,
        postCount: UInt,
        followersCount: UInt,
        followingsCount: UInt
    ) {
        self.userId = userId
        self.postCount = postCount
        self.followersCount = followersCount
        self.followingsCount = followingsCount
    }
}