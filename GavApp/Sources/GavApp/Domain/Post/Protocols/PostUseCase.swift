import Foundation

public struct PostUseCase {
    private let repository: any PostRepository

    public init(repository: any PostRepository) {
        self.repository = repository
    }

    public func create(userID: UUID, content: String, imageUrl: String?) async throws -> Post {
        return try await repository.create(
            userID: userID,
            content: content,
            imageUrl: imageUrl
        )
    }

    public func get(id: UUID) async throws -> Post {
        return try await repository.get(id: id)
    }

    public func listByUser(userID: UUID) async throws -> [Post] {
        return try await repository.listByUser(userID: userID)
    }

    public func delete(userID: UUID, id: UUID) async throws {
        try await repository.delete(userID: userID, id: id)
    }
}