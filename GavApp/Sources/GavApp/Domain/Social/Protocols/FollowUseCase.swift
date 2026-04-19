import Foundation

public struct FollowUseCase {
    private let repository: any FollowRepository

    public init(repository: any FollowRepository) {
        self.repository = repository
    }

    public func follow(follow: Follow) async throws {
        try await repository.follow(follow: follow)
    }

    public func unfollow(follow: Follow) async throws {
        try await repository.unfollow(follow: follow)
    }

    public func follows(follow: Follow) async throws -> Bool {
        return try await repository.follows(follow: follow)
    }

    public func getFollowers(userID: UUID) async throws -> [Follow] {
        return try await repository.getFollowers(userID: userID)
    }

    public func getFollowing(userID: UUID) async throws -> [Follow] {
        return try await repository.getFollowing(userID: userID)
    }
}