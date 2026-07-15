import Foundation

public protocol SettingsRepository {
    func getSettings() async throws -> UserSettings
    func updateSettings(input: UpdateUserSettingsInput) async throws
}