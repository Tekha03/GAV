import Foundation

class AuthManager {
    private let keychain = KeychainWrapper()
    private let accessTokenKey = "access_token"
    private let refreshTokenKey = "refresh_token"

    func currentToken() -> String? {
        return keychain.get(accessTokenKey)
    }

    func saveTokens(tokens: Tokens) {
        keychain.set(tokens.accessToken, forKey: accessTokenKey)
        keychain.set(tokens.refreshToken, forKey: refreshTokenKey)
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
    }
}

struct Tokens {
    let accessToken: String
    let refreshToken: String
    let userId: UUID
}