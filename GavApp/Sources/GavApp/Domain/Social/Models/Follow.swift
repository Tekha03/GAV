import Foundation

public struct Follow: Hashable, Sendable {
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

extension Follow {
    func unpack() -> (UUID, UUID) {
        (followerId, followingId)
    }
}