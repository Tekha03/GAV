import Foundation

public struct UpdateUserSettingsInput: Encodable {
    public let profilePrivacy: Bool?
    public let showLocation: Bool?
    public let allowMessages: Bool?

    public init(
        profilePrivacy: Bool? = nil,
        showLocation: Bool? = nil,
        allowMessages: Bool? = nil
    ) {
        self.profilePrivacy = profilePrivacy
        self.showLocation = showLocation
        self.allowMessages = allowMessages
    }
}