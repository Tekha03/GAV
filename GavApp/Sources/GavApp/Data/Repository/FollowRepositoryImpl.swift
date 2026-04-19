import Foundation

final class FollowRepositoryImpl: FollowRepository {
    private let api: any FollowServiceAPIProtocol
    private let mapper: FollowMapper

    init(api: any FollowServiceAPIProtocol, mapper: FollowMapper = FollowMapper()) {
        self.api = api
        self.mapper = mapper
    }

    func follow(follow: Follow) async throws {
        let followingID = follow.followingId
        try await api.follow(userID: followingID)
    }

    func unfollow(follow: Follow) async throws {
        let followingID = follow.followingId
        try await api.unfollow(userID: followingID)
    }

    func follows(follow: Follow) async throws -> Bool {
        fatalError("API не проверяет наличие follow напрямую")
    }

    func getFollowers(userID: UUID) async throws -> [Follow] {
        let models = try await api.getFollowers(userID: userID)
        return models.map { mapper.from(model: $0) }
    }

    func getFollowing(userID: UUID) async throws -> [Follow] {
        let models = try await api.getFollowing(userID: userID)
        return models.map { mapper.from(model: $0) }
    }
}