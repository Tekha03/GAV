import Domain
import Foundation

public struct StatsUseCase {
    private let repository: any StatsRepository

    public init(repository: any StatsRepository) {
        self.repository = repository
    }

    public func userStats(userID: UUID) async throws -> UserStats {
        return try await repository.userStats(userID: userID)
    }

    public func profileStats(userID: UUID) async throws -> ProfileStats {
        return try await repository.profileStats(userID: userID)
    }

    public func postStats(postID: UUID) async throws -> PostStats {
        return try await repository.postStats(postID: postID)
    }

    public func incrementPosts(userID: UUID) async throws {
        try await repository.incrementPosts(userID: userID)
    }

    public func incrementFollowers(userID: UUID) async throws {
        try await repository.incrementFollowers(userID: userID)
    }

    public func incrementFollowings(userID: UUID) async throws {
        try await repository.incrementFollowings(userID: userID)
    }

    public func incrementDogs(userID: UUID) async throws {
        try await repository.incrementDogs(userID: userID)
    }

    public func decrementPosts(userID: UUID) async throws {
        try await repository.decrementPosts(userID: userID)
    }

    public func decrementFollowers(userID: UUID) async throws {
        try await repository.decrementFollowers(userID: userID)
    }

    public func decrementFollowings(userID: UUID) async throws {
        try await repository.decrementFollowings(userID: userID)
    }

    public func decrementDogs(userID: UUID) async throws {
        try await repository.decrementDogs(userID: userID)
    }

    public func incrementPostLikes(postID: UUID) async throws {
        try await repository.incrementPostLikes(postID: postID)
    }

    public func incrementPostComments(postID: UUID) async throws {
        try await repository.incrementPostComments(postID: postID)
    }

    public func decrementPostLikes(postID: UUID) async throws {
        try await repository.decrementPostLikes(postID: postID)
    }

    public func decrementPostComments(postID: UUID) async throws {
        try await repository.decrementPostComments(postID: postID)
    }
}