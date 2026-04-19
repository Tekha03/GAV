public struct AuthMapper {
    public init() {}

    // user profile
    public func from(model: AuthModel) -> UserInfo {
        UserInfo(
            id: model.id,
            email: model.email
        )
    }

    // tokens
    public func from(model: AuthTokensModel) -> AuthTokens {
        AuthTokens(
            accessToken: model.accessToken,
            refreshToken: model.refreshToken
        )
    }

    public func to(model: AuthTokens) -> AuthTokensModel {
        AuthTokensModel(
            accessToken: model.accessToken,
            refreshToken: model.refreshToken
        )
    }
}