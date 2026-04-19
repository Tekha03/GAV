import Foundation

public struct CreatePostInput: Encodable {
    public let userId: UUID
    public let content: String
    public let imageUrl: String?

    public init(
        userId: UUID,
        content: String,
        imageUrl: String? = nil
    ) {
        self.userId = userId
        self.content = content
        self.imageUrl = imageUrl
    }
}

public struct UpdatePostInput: Encodable {
    public let content: String?
    public let imageUrl: String?

    public init(
        content: String? = nil,
        imageUrl: String? = nil
    ) {
        self.content = content
        self.imageUrl = imageUrl
    }
}