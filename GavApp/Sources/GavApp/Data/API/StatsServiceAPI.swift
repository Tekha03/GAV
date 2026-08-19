import Foundation

protocol StatsServiceAPIProtocol: Sendable {
    func userStats(userID: UUID) async throws -> UserStatsModel
    func profileStats(userID: UUID) async throws -> ProfileStatsModel
    func postStats(postID: UUID) async throws -> PostStatsModel
}

@available(macOS 12.0, *)
final class StatsServiceAPI: StatsServiceAPIProtocol, @unchecked Sendable {
    private let base: BaseAPI

    init(
        baseURL: URL,
        session: URLSession = .shared,
        authManager: AuthManager
    ) {
        self.base = BaseAPI(
            baseURL: baseURL,
            session: session,
            authManager: authManager
        )
    }

    func userStats(userID: UUID) async throws -> UserStatsModel {
        let path = "/api/v1/stats/user/\(userID.uuidString)"
        let data = try await base.request(path)
        return try JSONDecoder().decode(UserStatsModel.self, from: data)
    }

    func profileStats(userID: UUID) async throws -> ProfileStatsModel {
        let path = "/api/v1/stats/profile/\(userID.uuidString)"
        let data = try await base.request(path)
        return try JSONDecoder().decode(ProfileStatsModel.self, from: data)
    }

    func postStats(postID: UUID) async throws -> PostStatsModel {
        let path = "/api/v1/stats/post/\(postID.uuidString)"
        let data = try await base.request(path)
        return try JSONDecoder().decode(PostStatsModel.self, from: data)
    }
}
