import Foundation

public protocol FollowRepository {
    func follow(follow: Follow) async throws
    func unfollow(follow: Follow) async throws
    func follows(follow: Follow) async throws -> Bool
    func getFollowers(userID: UUID) async throws -> [Follow]
    func getFollowing(userID: UUID) async throws -> [Follow]
}