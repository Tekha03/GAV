import Domain
import SharedModels

public struct UserSettingsMapper {
    public static func from(model: UserSettingsModel) -> Domain.UserSettings {
        return Domain.UserSettings(
            userId: model.userId,
            profilePrivacy: model.profilePrivacy,
            showLocation: model.showLocation,
            allowMessages: model.allowMessages
        )
    }

    public static func to(model: Domain.UserSettings) -> UserSettingsModel {
        return UserSettingsModel(
            userId: model.userId,
            profilePrivacy: model.profilePrivacy,
            showLocation: model.showLocation,
            allowMessages: model.allowMessages
        )
    }
}