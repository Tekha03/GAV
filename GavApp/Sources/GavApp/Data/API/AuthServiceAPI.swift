import Foundation

protocol AuthServiceAPIProtocol {
    func register(email: String, password: String) async throws -> AuthModel
    func login(email: String, password: String) async throws -> AuthTokensModel
    func refreshToken(refreshToken: String) async throws -> AuthTokensModel
    func logout() async throws
    func getMe() async throws -> UserModel
}

@available(macOS 12.0, *)
final class AuthServiceAPI: AuthServiceAPIProtocol {
    private let base: BaseAPI

    init(baseURL: URL, session: URLSession = .shared, authManager: AuthManager) {
        self.base = BaseAPI(baseURL: baseURL, session: session, authManager: authManager)
    }

    func register(email: String, password: String) async throws -> AuthModel {
        let path = "/api/v1/auth/register"
        let body = try JSONEncoder().encode(["email": email, "password": password])
        let data = try await base.request(path, method: "POST", body: body)
        return try JSONDecoder().decode(AuthModel.self, from: data)
    }

    func login(email: String, password: String) async throws -> AuthTokensModel {
        let path = "/api/v1/auth/login"
        let body = try JSONEncoder().encode(["email": email, "password": password])
        let data = try await base.request(path, method: "POST", body: body)
        return try JSONDecoder().decode(AuthTokensModel.self, from: data)
    }

    func refreshToken(refreshToken: String) async throws -> AuthTokensModel {
        let path = "/api/v1/auth/refresh"
        let body = try JSONEncoder().encode(["refresh_token": refreshToken])
        let data = try await base.request(path, method: "POST", body: body, requiresAuth: false)
        return try JSONDecoder().decode(AuthTokensModel.self, from: data)
    }

    func logout() async throws {
        let path = "/api/v1/auth/logout"
        _ = try await base.request(path, method: "POST")
    }

    func getMe() async throws -> UserModel {
        let path = "/api/v1/auth/me"
        let data = try await base.request(path, method: "GET")
        return try JSONDecoder().decode(UserModel.self, from: data)
    }
}
