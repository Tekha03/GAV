import Foundation

public struct AuthUseCase {
    private let repository: any AuthRepository

    public init(repository: any AuthRepository) {
        self.repository = repository
    }

    public func login(email: String, password: String) async throws -> AuthTokens {
        return try await repository.login(email: email, password: password)
    }

    public func register(email: String, password: String) async throws -> AuthTokens {
        return try await repository.register(email: email, password: password)
    }

    public func refreshTokens(refreshToken: String) async throws -> AuthTokens {
        return try await repository.refreshTokens(refreshToken: refreshToken)
    }

    public func logout() async throws {
        try await repository.logout()
    }
}