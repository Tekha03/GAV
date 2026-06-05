import Foundation

public struct PostStats: Equatable, Sendable {
    public let id: UUID
    public let postId: UUID
    public let likesCount: UInt
    public let commentsCount: UInt
    public let createdAt: Date
    public let updatedAt: Date

    public init(
        id: UUID,
        postId: UUID,
        likesCount: UInt,
        commentsCount: UInt,
        createdAt: Date,
        updatedAt: Date
    ) {
        self.id = id
        self.postId = postId
        self.likesCount = likesCount
        self.commentsCount = commentsCount
        self.createdAt = createdAt
        self.updatedAt = updatedAt
    }
}