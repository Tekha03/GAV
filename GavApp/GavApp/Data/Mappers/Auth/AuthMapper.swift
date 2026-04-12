import Domain
import SharedModels

public struct AuthMapper {
    public static func from(model: AuthModel) -> Domain.UserInfo {
        return Domain.UserInfo(id: model.id, email: model.email)
    }

    public static func to(model: Domain.UserInfo) -> AuthModel {
        return AuthModel(id: model.id, email: model.email)
    }

    public static func from(model: AuthTokensModel) -> Domain.AuthTokens {
        return Domain.AuthTokens(
            accessToken: model.accessToken,
            refreshToken: model.refreshToken
        )
    }

    public static func to(model: Domain.AuthTokens) -> AuthTokensModel {
        return AuthTokensModel(
            accessToken: model.accessToken,
            refreshToken: model.refreshToken
        )
    }
}