import Domain
import Foundation

public protocol StatsRepository {
    func userStats(userID: UUID) async throws -> UserStats
    func profileStats(userID: UUID) async throws -> ProfileStats
    func postStats(postID: UUID) async throws -> PostStats

    func incrementPosts(userID: UUID) async throws
    func incrementFollowers(userID: UUID) async throws
    func incrementFollowings(userID: UUID) async throws
    func incrementDogs(userID: UUID) async throws

    func decrementPosts(userID: UUID) async throws
    func decrementFollowers(userID: UUID) async throws
    func decrementFollowings(userID: UUID) async throws
    func decrementDogs(userID: UUID) async throws

    func incrementPostLikes(postID: UUID) async throws
    func incrementPostComments(postID: UUID) async throws
    func decrementPostLikes(postID: UUID) async throws
    func decrementPostComments(postID: UUID) async throws
}