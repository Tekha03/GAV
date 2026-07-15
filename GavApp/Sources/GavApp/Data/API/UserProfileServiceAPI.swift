import Foundation

protocol UserProfileServiceAPIProtocol: Sendable {
    func create(userID: UUID, model: UserProfileModel) async throws -> UserProfileModel
    func getByUserID(userID: UUID) async throws -> UserProfileModel
    func search(query: String, limit: Int) async throws -> [UserProfileModel]
    func update(userID: UUID, input: UpdateProfileInput) async throws
    func delete(userID: UUID) async throws
}

@available(macOS 12.0, *)
final class UserProfileServiceAPI: UserProfileServiceAPIProtocol, @unchecked Sendable {
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

    func create(userID: UUID, model: UserProfileModel) async throws -> UserProfileModel {
        let path = "/api/v1/users/\(userID.uuidString)/profile"
        let body = try JSONEncoder().encode(model)
        let data = try await base.request(path, method: "POST", body: body)
        return try JSONDecoder().decode(UserProfileModel.self, from: data)
    }

    func getByUserID(userID: UUID) async throws -> UserProfileModel {
        let path = "/api/v1/users/\(userID.uuidString)/profile"
        let data = try await base.request(path)
        return try JSONDecoder().decode(UserProfileModel.self, from: data)
    }

    func search(query: String, limit: Int = 10) async throws -> [UserProfileModel] {
        let encodedQuery = query.addingPercentEncoding(withAllowedCharacters: .urlQueryAllowed) ?? query
        let path = "/api/v1/profiles/search?q=\(encodedQuery)&limit=\(limit)"
        let data = try await base.request(path)
        return try JSONDecoder().decode([UserProfileModel].self, from: data)
    }

    func update(userID: UUID, input: UpdateProfileInput) async throws {
        let path = "/api/v1/users/\(userID.uuidString)/profile"
        let body = try JSONEncoder().encode(input)
        _ = try await base.request(path, method: "PUT", body: body)
    }

    func delete(userID: UUID) async throws {
        let path = "/api/v1/users/\(userID.uuidString)/profile"
        _ = try await base.request(path, method: "DELETE")
    }
}
