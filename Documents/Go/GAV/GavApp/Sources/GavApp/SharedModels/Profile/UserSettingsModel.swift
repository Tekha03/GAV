import Foundation

public struct UserSettingsModel: Codable, Equatable {
    public let userId: UUID
    public let profilePrivacy: Bool
    public let showLocation: Bool
    public let allowMessages: Bool
}

extension UserSettingsModel {
    private enum CodingKeys: String, CodingKey {
        case userId = "user_id"
        case profilePrivacy = "profile_privacy"
        case showLocation = "show_location"
        case allowMessages = "allow_messages"
    }
}