import Foundation

public protocol AuthRepository {
    func login(email: String, password: String) async throws -> AuthTokens
    func register(email: String, password: String) async throws -> AuthModel
    func refreshTokens(refreshToken: String) async throws -> AuthTokens
    func logout() async throws
}