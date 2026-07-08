import Foundation

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

@available(macOS 12.0, *)
final class PostServiceAPI: PostServiceAPIProtocol {
    private let base: BaseAPI
    private let decoder: JSONDecoder

    init(baseURL: URL, session: URLSession = .shared, authManager: AuthManager) {
        self.base = BaseAPI(baseURL: baseURL, session: session, authManager: authManager)
        let decoder = JSONDecoder()
        decoder.dateDecodingStrategy = .custom { decoder in
            let container = try decoder.singleValueContainer()
            let value = try container.decode(String.self)
            let isoWithFraction = ISO8601DateFormatter()
            isoWithFraction.formatOptions = [.withInternetDateTime, .withFractionalSeconds]
            if let date = isoWithFraction.date(from: value) {
                return date
            }
            let iso = ISO8601DateFormatter()
            if let date = iso.date(from: value) {
                return date
            }
            throw DecodingError.dataCorruptedError(
                in: container,
                debugDescription: "Invalid ISO8601 date: \(value)"
            )
        }
        self.decoder = decoder
    }

    func create(userID: UUID, content: String, imageUrl: String?) async throws -> PostModel {
        let path = "/api/v1/posts"
        let input = CreatePostInput(userId: userID, content: content, imageUrl: imageUrl)
        let body = try JSONEncoder().encode(input)
        let data = try await base.request(path, method: "POST", body: body)
        return try decoder.decode(PostModel.self, from: data)
    }

    func getByID(id: UUID) async throws -> PostModel {
        let path = "/api/v1/posts/\(id.uuidString)"
        let data = try await base.request(path)
        return try decoder.decode(PostModel.self, from: data)
    }

    func listByUser(userID: UUID) async throws -> [PostModel] {
        var components = URLComponents()
        components.queryItems = [
            URLQueryItem(name: "user_id", value: userID.uuidString)
        ]
        let query = components.percentEncodedQuery ?? ""
        let path = "/api/v1/posts" + (query.isEmpty ? "" : "?\(query)")
        let data = try await base.request(path)
        return try decoder.decode([PostModel].self, from: data)
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
        let input = CreateCommentInput(postId: postID, userId: userID, content: content)
        let body = try JSONEncoder().encode(input)
        _ = try await base.request(path, method: "POST", body: body)
        return CommentModel(id: UUID(), postId: postID, userId: userID, content: content, createdAt: .now)
    }

    func listCommentsByPostID(postID: UUID) async throws -> [CommentModel] {
        let path = "/api/v1/posts/\(postID.uuidString)/comments"
        let data = try await base.request(path)
        return try decoder.decode([CommentModel].self, from: data)
    }

    func deleteComment(id: UUID) async throws {
        let path = "/api/v1/comments/\(id.uuidString)"
        _ = try await base.request(path, method: "DELETE")
    }
}
