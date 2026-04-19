import Foundation

final class AuthRepositoryImpl: AuthRepository {
    private let api: any AuthServiceAPIProtocol
    private let mapper: AuthMapper

    init(
        api: any AuthServiceAPIProtocol,
        mapper: AuthMapper = AuthMapper()
    ) {
        self.api = api
        self.mapper = mapper
    }

    func login(email: String, password: String) async throws -> AuthTokens {
        let model = try await api.login(email: email, password: password)
        return mapper.from(model: model)
    }

    func register(email: String, password: String) async throws -> AuthModel {
        let model = try await api.register(email: email, password: password)
        return model
    }

    func refreshTokens(refreshToken: String) async throws -> AuthTokens {
        let model = try await api.refreshToken(refreshToken: refreshToken)
        return mapper.from(model: model)
    }

    func logout() async throws {
        try await api.logout()
    }
}
