import Foundation

final class AuthManager: Sendable {
    private let keychain = KeychainWrapper()
    private let accessTokenKey = "access_token"
    private let refreshTokenKey = "refresh_token"
    private let userIdKey = "user_id"

    func currentToken() -> String? {
        keychain.get(accessTokenKey)
    }

    func currentUserId() -> UUID? {
        guard let raw = keychain.get(userIdKey) else { return nil }
        return UUID(uuidString: raw)
    }

    func saveTokens(tokens: Tokens) {
        keychain.set(tokens.accessToken, forKey: accessTokenKey)
        keychain.set(tokens.refreshToken, forKey: refreshTokenKey)
        keychain.set(tokens.userId.uuidString, forKey: userIdKey)
    }

    func saveSession(accessToken: String, refreshToken: String, userId: UUID) {
        keychain.set(accessToken, forKey: accessTokenKey)
        keychain.set(refreshToken, forKey: refreshTokenKey)
        keychain.set(userId.uuidString, forKey: userIdKey)
    }

    func saveTokenPair(accessToken: String, refreshToken: String) {
        keychain.set(accessToken, forKey: accessTokenKey)
        keychain.set(refreshToken, forKey: refreshTokenKey)
    }

    func getAccessToken() -> String? {
        keychain.get(accessTokenKey)
    }

    func getRefreshToken() -> String? {
        keychain.get(refreshTokenKey)
    }

    func clearTokens() {
        keychain.delete(accessTokenKey)
        keychain.delete(refreshTokenKey)
        keychain.delete(userIdKey)
    }
}

struct Tokens {
    let accessToken: String
    let refreshToken: String
    let userId: UUID
}
