import Foundation

public struct AuthModel: Codable, Equatable {
    public let id: UUID
    public let email: String
}

public struct AuthTokensModel: Codable, Equatable {
    public let accessToken: String
    public let refreshToken: String
}

extension AuthTokensModel {
    private enum CodingKeys: String, CodingKey {
        case accessToken = "access_token"
        case refreshToken = "refresh_token"
    }
}