import Foundation

protocol FeedServiceAPIProtocol: Sendable {
    func getFeed(
        userID: UUID,
        before: Date?,
        limit: Int
    ) async throws -> [PostModel]
}

@available(macOS 12.0, *)
final class FeedServiceAPI: FeedServiceAPIProtocol, @unchecked Sendable {
    private let base: BaseAPI
    private let decoder: JSONDecoder

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
        let decoder = JSONDecoder()
        decoder.dateDecodingStrategy = .custom { decoder in
            let container = try decoder.singleValueContainer()
            let value = try container.decode(String.self)
            let isoWithFraction = ISO8601DateFormatter()
            isoWithFraction.formatOptions = [.withInternetDateTime, .withFractionalSeconds]
            if let date = isoWithFraction.date(from: value) {
                return date
            }
            let iso = ISO8601DateFormatter()
            if let date = iso.date(from: value) {
                return date
            }
            throw DecodingError.dataCorruptedError(
                in: container,
                debugDescription: "Invalid ISO8601 date: \(value)"
            )
        }
        self.decoder = decoder
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
            let formatter = ISO8601DateFormatter()
            formatter.formatOptions = [.withInternetDateTime, .withFractionalSeconds]
            components.queryItems?.append(
                URLQueryItem(
                    name: "cursor",
                    value: formatter.string(from: before)
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

        return try decoder.decode(FeedResponseModel.self, from: data).posts
    }
}

private struct FeedResponseModel: Decodable {
    let posts: [PostModel]
}
