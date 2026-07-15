import Foundation

public protocol PostRepository {
    func create(userID: UUID, content: String, imageUrl: String?) async throws -> Post
    func get(id: UUID) async throws -> Post
    func listByUser(userID: UUID) async throws -> [Post]
    func delete(userID: UUID, id: UUID) async throws
}