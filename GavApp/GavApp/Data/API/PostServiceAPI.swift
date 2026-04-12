import Foundation
import SharedModels

protocol PostServiceAPIProtocol {
    func create(userID: UUID, content: String, imageUrl: String?) async throws -> PostModel
    func getByID(id: UUID) async throws -> PostModel
    func listByUser(userID: UUID) async throws -> [PostModel]
    func delete(id: UUID) async throws

    func addLike(postID: UUID) async throws
    func removeLike(postID: UUID) async throws

    func createComment(postID: UUID, userID: UUID, content: String) async throws -> CommentModel
    func listCommentsByPostID(postID: UUID) async throws -> [CommentModel]
    func deleteComment(id: UUID) async throws
}

final class PostServiceAPI: PostServiceAPIProtocol {
    private let base: BaseAPI

    init(baseURL: URL, session: URLSession = .shared, authManager: AuthManager) {
        self.base = BaseAPI(baseURL: baseURL, session: session, authManager: authManager)
    }

    func create(userID: UUID, content: String, imageUrl: String?) async throws -> PostModel {
        let path = "/api/v1/posts"
        let input = CreatePostInput(userId: userID, content: content, imageUrl: imageUrl)
        let body = try JSONEncoder().encode(input.toModel())
        let data = try await base.request(path, method: "POST", body: body)
        return try JSONDecoder().decode(PostModel.self, from: data)
    }

    func getByID(id: UUID) async throws -> PostModel {
        let path = "/api/v1/posts/\(id.uuidString)"
        let data = try await base.request(path)
        return try JSONDecoder().decode(PostModel.self, from: data)
    }

    func listByUser(userID: UUID) async throws -> [PostModel] {
        let path = "/api/v1/posts"
        let data = try await base.request(path)
        return try JSONDecoder().decode([PostModel].self, from: data)
    }

    func delete(id: UUID) async throws {
        let path = "/api/v1/posts/\(id.uuidString)"
        _ = try await base.request(path, method: "DELETE")
    }

    func addLike(postID: UUID) async throws {
        let path = "/api/v1/posts/\(postID.uuidString)/likes"
        _ = try await base.request(path, method: "POST")
    }

    func removeLike(postID: UUID) async throws {
        let path = "/api/v1/posts/\(postID.uuidString)/likes"
        _ = try await base.request(path, method: "DELETE")
    }

    func createComment(postID: UUID, userID: UUID, content: String) async throws -> CommentModel {
        let path = "/api/v1/posts/\(postID.uuidString)/comments"
        let input = CreateCommentInput(userId: userID, content: content)
        let body = try JSONEncoder().encode(input.toModel())
        let data = try await base.request(path, method: "POST", body: body)
        return try JSONDecoder().decode(CommentModel.self, from: data)
    }

    func listCommentsByPostID(postID: UUID) async throws -> [CommentModel] {
        let path = "/api/v1/posts/\(postID.uuidString)/comments"
        let data = try await base.request(path)
        return try JSONDecoder().decode([CommentModel].self, from: data)
    }

    func deleteComment(id: UUID) async throws {
        let path = "/api/v1/comments/\(id.uuidString)"
        _ = try await base.request(path, method: "DELETE")
    }
}