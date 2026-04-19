import Foundation

protocol FollowServiceAPIProtocol {
    func follow(userID: UUID) async throws
    func unfollow(userID: UUID) async throws
    func getFollowers(userID: UUID) async throws -> [FollowModel]
    func getFollowing(userID: UUID) async throws -> [FollowModel]
}

@available(macOS 12.0, *)
final class FollowServiceAPI: FollowServiceAPIProtocol {
    private let base: BaseAPI

    init(baseURL: URL, session: URLSession = .shared, authManager: AuthManager) {
        self.base = BaseAPI(baseURL: baseURL, session: session, authManager: authManager)
    }

    func follow(userID: UUID) async throws {
        let path = "/api/v1/follows/\(userID.uuidString)"
        _ = try await base.request(path, method: "POST")
    }

    func unfollow(userID: UUID) async throws {
        let path = "/api/v1/follows/\(userID.uuidString)"
        _ = try await base.request(path, method: "DELETE")
    }

    func getFollowers(userID: UUID) async throws -> [FollowModel] {
        let path = "/api/v1/follows/followers/\(userID.uuidString)"
        let data = try await base.request(path)
        return try JSONDecoder().decode([FollowModel].self, from: data)
    }

    func getFollowing(userID: UUID) async throws -> [FollowModel] {
        let path = "/api/v1/follows/following/\(userID.uuidString)"
        let data = try await base.request(path)
        return try JSONDecoder().decode([FollowModel].self, from: data)
    }
}