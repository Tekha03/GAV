import Foundation

public struct Like: Hashable, Sendable {
    public let userId: UUID
    public let postId: UUID

    public init(userId: UUID, postId: UUID) {
        self.userId = userId
        self.postId = postId
    }
}
