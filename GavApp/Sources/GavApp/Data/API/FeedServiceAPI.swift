import Foundation

protocol FeedServiceAPIProtocol {
    func getFeed(
        userID: UUID,
        before: Date?,
        limit: Int
    ) async throws -> [PostModel]
}

@available(macOS 12.0, *)
final class FeedServiceAPI: FeedServiceAPIProtocol {
    private let base: BaseAPI

    init(
        baseURL: URL,
        session: URLSession = .shared,
        authManager: AuthManager
    ) {
        self.base = BaseAPI(
            baseURL: baseURL,
            session: session,
            authManager: authManager
        )
    }

    func getFeed(
        userID: UUID,
        before: Date?,
        limit: Int
    ) async throws -> [PostModel] {

        var components = URLComponents()
        components.queryItems = []

        components.queryItems?.append(
            URLQueryItem(name: "user_id", value: userID.uuidString)
        )

        if let before = before {
            components.queryItems?.append(
                URLQueryItem(
                    name: "before",
                    value: String(Int(before.timeIntervalSince1970))
                )
            )
        }

        if limit > 0 {
            components.queryItems?.append(
                URLQueryItem(name: "limit", value: "\(limit)")
            )
        }

        let query = components.percentEncodedQuery ?? ""
        let path = "/api/v1/feed" + (query.isEmpty ? "" : "?\(query)")

        let data = try await base.request(path)

        return try JSONDecoder().decode([PostModel].self, from: data)
    }
}