import Foundation

final class FeedRepositoryImpl: FeedRepository {
    private let api: any FeedServiceAPIProtocol
    private let mapper: PostMapper

    init(
        api: any FeedServiceAPIProtocol,
        mapper: PostMapper = PostMapper()
    ) {
        self.api = api
        self.mapper = mapper
    }

    func getFeed(
        userID: UUID,
        before: Date?,
        limit: Int
    ) async throws -> ([Post], nextPageToken: Date?) {

        let models = try await api.getFeed(
            userID: userID,
            before: before,
            limit: limit
        )

        let posts = models.map { mapper.from(model: $0) }
        let nextPageToken = models.last?.createdAt

        return (posts, nextPageToken: nextPageToken)
    }
}