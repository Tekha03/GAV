import Domain
import Foundation

public struct CommentUseCase {
    private let repository: any CommentRepository

    public init(repository: any CommentRepository) {
        self.repository = repository
    }

    public func create(postID: UUID, userID: UUID, content: String) async throws -> Comment {
        return try await repository.create(postID: postID, userID: userID, content: content)
    }

    public func get(id: UUID) async throws -> Comment {
        return try await repository.get(id: id)
    }

    public func listByPostID(postID: UUID) async throws -> [Comment] {
        return try await repository.listByPostID(postID: postID)
    }

    public func delete(userID: UUID, id: UUID) async throws {
        try await repository.delete(userID: userID, id: id)
    }
}