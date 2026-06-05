import Foundation

final class SettingsRepositoryImpl: SettingsRepository {
    private let api: any SettingsServiceAPIProtocol
    private let mapper: UserSettingsMapper

    init(api: any SettingsServiceAPIProtocol, mapper: UserSettingsMapper = UserSettingsMapper()) {
        self.api = api
        self.mapper = mapper
    }

    func getSettings() async throws -> UserSettings {
        let model = try await api.getSettings()
        return UserSettingsMapper.from(model: model)
    }

    func updateSettings(input: UpdateUserSettingsInput) async throws {
        try await api.updateSettings(input: input)
    }
}