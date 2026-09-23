import Foundation

final class AuthManager: Sendable {
    private let keychain = KeychainWrapper()
    private let refreshBaseURL: URL
    private let refreshCoordinator = TokenRefreshCoordinator()
    private let accessTokenKey = "access_token"
    private let refreshTokenKey = "refresh_token"
    private let userIdKey = "user_id"

    init(refreshBaseURL: URL) {
        self.refreshBaseURL = refreshBaseURL
    }

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

    func refreshAccessToken(after failedToken: String?, using session: URLSession) async throws -> String {
        try await refreshCoordinator.refresh(
            after: failedToken,
            using: session,
            baseURL: refreshBaseURL,
            keychain: keychain
        )
    }
}

extension Notification.Name {
    static let gavSessionExpired = Notification.Name("gavSessionExpired")
}

@MainActor
private final class TokenRefreshCoordinator {
    private var inFlight: Task<AuthTokensModel, Error>?

    func refresh(
        after failedToken: String?,
        using session: URLSession,
        baseURL: URL,
        keychain: KeychainWrapper
    ) async throws -> String {
        if let current = keychain.get("access_token"), current != failedToken {
            return current
        }
        if let inFlight {
            return try await inFlight.value.accessToken
        }
        guard let refreshToken = keychain.get("refresh_token"),
              let url = URL(string: "/api/v1/auth/refresh", relativeTo: baseURL)?.absoluteURL else {
            expire(keychain)
            throw APIError.invalidResponse(statusCode: 401)
        }

        let task = Task<AuthTokensModel, Error> {
            var request = URLRequest(url: url)
            request.httpMethod = "POST"
            request.setValue("application/json", forHTTPHeaderField: "Content-Type")
            request.httpBody = try JSONEncoder().encode(["refresh_token": refreshToken])
            let (data, response) = try await session.data(for: request)
            guard let response = response as? HTTPURLResponse else {
                throw APIError.invalidResponse(statusCode: 0)
            }
            guard (200...299).contains(response.statusCode) else {
                throw APIError.invalidResponse(statusCode: response.statusCode)
            }
            if let tokens = try? JSONDecoder().decode(AuthTokensModel.self, from: data) {
                return tokens
            }
            return try JSONDecoder().decode(RefreshEnvelope.self, from: data).token
        }
        inFlight = task
        defer { inFlight = nil }
        do {
            let tokens = try await task.value
            keychain.set(tokens.accessToken, forKey: "access_token")
            keychain.set(tokens.refreshToken, forKey: "refresh_token")
            return tokens.accessToken
        } catch {
            if let apiError = error as? APIError,
               case .invalidResponse(let status) = apiError,
               status == 400 || status == 401 {
                expire(keychain)
            }
            throw error
        }
    }

    private func expire(_ keychain: KeychainWrapper) {
        keychain.delete("access_token")
        keychain.delete("refresh_token")
        keychain.delete("user_id")
        NotificationCenter.default.post(name: .gavSessionExpired, object: nil)
    }
}

private struct RefreshEnvelope: Decodable {
    let token: AuthTokensModel
}

struct Tokens {
    let accessToken: String
    let refreshToken: String
    let userId: UUID
}
