import Foundation

protocol SettingsServiceAPIProtocol: Sendable {
    func getSettings() async throws -> UserSettingsModel
    func updateSettings(input: UpdateUserSettingsInput) async throws
}

@available(macOS 12.0, *)
final class SettingsServiceAPI: SettingsServiceAPIProtocol, @unchecked Sendable {
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

    func getSettings() async throws -> UserSettingsModel {
        let data = try await base.request("/api/v1/settings")
        return try JSONDecoder().decode(UserSettingsModel.self, from: data)
    }

    func updateSettings(input: UpdateUserSettingsInput) async throws {
        let body = try JSONEncoder().encode(input)
        _ = try await base.request(
            "/api/v1/settings",
            method: "PUT",
            body: body
        )
    }
}
