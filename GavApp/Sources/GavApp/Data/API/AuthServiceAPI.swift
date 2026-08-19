import Foundation

protocol AuthServiceAPIProtocol: Sendable {
    func register(email: String, password: String) async throws -> AuthTokensModel
    func login(email: String, password: String) async throws -> AuthTokensModel
    func refreshToken(refreshToken: String) async throws -> AuthTokensModel
    func logout() async throws
    func getMe() async throws -> UserModel
}

@available(macOS 12.0, *)
final class AuthServiceAPI: AuthServiceAPIProtocol, @unchecked Sendable {
    private let base: BaseAPI

    init(baseURL: URL, session: URLSession = .shared, authManager: AuthManager) {
        self.base = BaseAPI(baseURL: baseURL, session: session, authManager: authManager)
    }

    func register(email: String, password: String) async throws -> AuthTokensModel {
        let path = "/api/v1/auth/register"
        let body = try JSONEncoder().encode(["email": email, "password": password])
        let data = try await base.request(path, method: "POST", body: body, requiresAuth: false)
        return try decodeTokens(from: data)
    }

    func login(email: String, password: String) async throws -> AuthTokensModel {
        let path = "/api/v1/auth/login"
        let body = try JSONEncoder().encode(["email": email, "password": password])
        let data = try await base.request(path, method: "POST", body: body, requiresAuth: false)
        return try decodeTokens(from: data)
    }

    func refreshToken(refreshToken: String) async throws -> AuthTokensModel {
        let path = "/api/v1/auth/refresh"
        let body = try JSONEncoder().encode(["refresh_token": refreshToken])
        let data = try await base.request(path, method: "POST", body: body, requiresAuth: false)
        return try decodeTokens(from: data)
    }

    func logout() async throws {
        let path = "/api/v1/auth/logout"
        _ = try await base.request(path, method: "POST")
    }

    func getMe() async throws -> UserModel {
        let path = "/api/v1/auth/me"
        let data = try await base.request(path, method: "GET")
        let decoder = JSONDecoder()
        if let user = try? decoder.decode(UserModel.self, from: data) {
            return user
        }
        return try decoder.decode(CurrentUserResponse.self, from: data).model
    }

    private func decodeTokens(from data: Data) throws -> AuthTokensModel {
        let decoder = JSONDecoder()
        if let tokens = try? decoder.decode(AuthTokensModel.self, from: data) {
            return tokens
        }

        return try decoder.decode(AuthTokenEnvelope.self, from: data).token
    }
}

private struct AuthTokenEnvelope: Decodable {
    let token: AuthTokensModel
}

private struct CurrentUserResponse: Decodable {
    let id: UUID
    let email: String

    var model: UserModel {
        UserModel(
            id: id,
            email: email,
            role: "user",
            createdAt: Date(),
            updatedAt: Date()
        )
    }

    private enum CodingKeys: String, CodingKey {
        case id
        case email
        case legacyID = "ID"
        case legacyEmail = "Email"
    }

    init(from decoder: Decoder) throws {
        let container = try decoder.container(keyedBy: CodingKeys.self)
        id = try container.decodeIfPresent(UUID.self, forKey: .id)
            ?? container.decode(UUID.self, forKey: .legacyID)
        email = try container.decodeIfPresent(String.self, forKey: .email)
            ?? container.decode(String.self, forKey: .legacyEmail)
    }
}
