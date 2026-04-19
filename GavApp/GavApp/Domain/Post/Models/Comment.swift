public struct Comment: Identifiable, Equatable, Sendable {
    public let id: UUID
    public let postId: UUID
    public let userId: UUID
    public let content: String
    public let createdAt: Date

    public init(
        id: UUID,
        postId: UUID,
        userId: UUID,
        content: String,
        createdAt: Date
    ) {
        self.id = id
        self.postId = postId
        self.userId = userId
        self.content = content
        self.createdAt = createdAt
    }
}