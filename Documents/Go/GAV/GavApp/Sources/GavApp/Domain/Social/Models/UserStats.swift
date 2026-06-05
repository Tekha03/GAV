import Foundation

public struct UserStats: Equatable, Sendable {
    public let id: UUID
    public let userId: UUID
    public let postCount: UInt
    public let followersCount: UInt
    public let followingsCount: UInt
    public let dogsCount: UInt
    public let createdAt: Date
    public let updatedAt: Date

    public init(
        id: UUID,
        userId: UUID,
        postCount: UInt,
        followersCount: UInt,
        followingsCount: UInt,
        dogsCount: UInt,
        createdAt: Date,
        updatedAt: Date
    ) {
        self.id = id
        self.userId = userId
        self.postCount = postCount
        self.followersCount = followersCount
        self.followingsCount = followingsCount
        self.dogsCount = dogsCount
        self.createdAt = createdAt
        self.updatedAt = updatedAt
    }
}