import Foundation

public struct FollowerStats: Equatable, Sendable {
    public let followerId: UUID
    public let followingId: UUID

    public init(
        followerId: UUID,
        followingId: UUID
    ) {
        self.followerId = followerId
        self.followingId = followingId
    }
}