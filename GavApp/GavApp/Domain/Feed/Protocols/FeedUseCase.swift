import Domain
import Foundation

public struct FeedUseCase {
    private let repository: any FeedRepository

    public init(repository: any FeedRepository) {
        self.repository = repository
    }

    public func getFeed(userID: UUID, before: Date? = nil, limit: Int) async throws -> [Post] {
        let (posts, _) = try await repository.getFeed(userID: userID, before: before, limit: limit)
        return posts
    }
}