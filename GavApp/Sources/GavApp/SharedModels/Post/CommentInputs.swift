import Foundation

public struct CreateCommentInput: Encodable {
    public let postId: UUID
    public let userId: UUID
    public let content: String

    public init(postId: UUID, userId: UUID, content: String) {
        self.postId = postId
        self.userId = userId
        self.content = content
    }
}