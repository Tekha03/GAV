import Foundation

public protocol FeedRepository {
    func getFeed(
        userID: UUID,
        before: Date?,
        limit: Int
    ) async throws -> ([Post], nextPageToken: Date?)
}