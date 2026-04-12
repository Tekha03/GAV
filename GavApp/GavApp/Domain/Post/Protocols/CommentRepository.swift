import Domain
import Foundation

public protocol CommentRepository {
    func create(postID: UUID, userID: UUID, content: String) async throws -> Comment
    func get(id: UUID) async throws -> Comment
    func listByPostID(postID: UUID) async throws -> [Comment]
    func delete(userID: UUID, id: UUID) async throws
}