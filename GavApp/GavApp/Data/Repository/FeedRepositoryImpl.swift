import Foundation
import Domain
import Data
import SharedModels

final class FeedRepositoryImpl: FeedRepository {
    private let api: any FeedServiceAPIProtocol
    private let mapper: PostMapper

    init(api: any FeedServiceAPIProtocol, mapper: PostMapper = PostMapper()) {
        self.api = api
        self.mapper = mapper
    }

    func getFeed(userID: UUID, before: Date?, limit: Int) async throws -> [Post] {
        let models = try await api.getFeed(before: before, limit: limit)
        return models.map { PostMapper.from(model: $0) }
    }
}