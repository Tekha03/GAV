import Foundation
import Domain
import Data
import SharedModels

final class CommentRepositoryImpl: CommentRepository {
    private let api: any PostServiceAPIProtocol
    private let mapper: CommentMapper

    init(api: any PostServiceAPIProtocol, mapper: CommentMapper = CommentMapper()) {
        self.api = api
        self.mapper = mapper
    }

    func create(postID: UUID, userID: UUID, content: String) async throws -> Comment {
        let model = try await api.createComment(postID: postID, userID: userID, content: content)
        return CommentMapper.from(model: model)
    }

    func get(id: UUID) async throws -> Comment {
        fatalError("GET /comments/{id} не описан в роутере")
    }

    func listByPostID(postID: UUID) async throws -> [Comment] {
        let models = try await api.listCommentsByPostID(postID: postID)
        return models.map { CommentMapper.from(model: $0) }
    }

    func delete(userID: UUID, id: UUID) async throws {
        try await api.deleteComment(id: id)
    }
}