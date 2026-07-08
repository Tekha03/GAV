public struct UserSettingsMapper {
    public static func from(model: UserSettingsModel) -> UserSettings {
        return UserSettings(
            userId: model.userId,
            profilePrivacy: model.profilePrivacy,
            showLocation: model.showLocation,
            allowMessages: model.allowMessages
        )
    }

    public static func to(model: UserSettings) -> UserSettingsModel {
        return UserSettingsModel(
            userId: model.userId,
            profilePrivacy: model.profilePrivacy,
            showLocation: model.showLocation,
            allowMessages: model.allowMessages
        )
    }
}
