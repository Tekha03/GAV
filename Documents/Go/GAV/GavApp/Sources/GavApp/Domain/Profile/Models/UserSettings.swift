import Foundation

public struct UserSettings: Equatable, Sendable {
    public let userId: UUID
    public let profilePrivacy: Bool
    public let showLocation: Bool
    public let allowMessages: Bool

    public init(
        userId: UUID,
        profilePrivacy: Bool = true,
        showLocation: Bool = false,
        allowMessages: Bool = true
    ) {
        self.userId = userId
        self.profilePrivacy = profilePrivacy
        self.showLocation = showLocation
        self.allowMessages = allowMessages
    }
}