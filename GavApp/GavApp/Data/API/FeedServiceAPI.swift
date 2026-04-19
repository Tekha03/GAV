import Foundation
import SharedModels

protocol FeedServiceAPIProtocol {
    func getFeed(before: Date?, limit: Int) async throws -> [PostModel]
}

final class FeedServiceAPI: FeedServiceAPIProtocol {
    private let base: BaseAPI

    init(baseURL: URL, session: URLSession = .shared, authManager: AuthManager) {
        self.base = BaseAPI(baseURL: baseURL, session: session, authManager: authManager)
    }

    func getFeed(before: Date?, limit: Int) async throws -> [PostModel] {
        var path = "/api/v1/feed"
        if let before = before {
            path += "?before=\(before.timeIntervalSince1970)"
        }
        if limit > 0 {
            path += (path.contains("?") ? "&" : "?") + "limit=\(limit)"
        }
        let data = try await base.request(path)
        return try JSONDecoder().decode([PostModel].self, from: data)
    }
}