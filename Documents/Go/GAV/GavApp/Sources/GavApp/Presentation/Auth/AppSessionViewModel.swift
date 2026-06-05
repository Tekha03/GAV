import Foundation
import Combine

@MainActor
final class AppSessionViewModel: ObservableObject {
    @Published private(set) var isAuthenticated: Bool
    @Published var isLoading = false
    @Published var errorMessage: String?

    private let authService: AuthServiceAPIProtocol
    private let authManager: AuthManager
    private let appViewModel: AppViewModel

    init(
        authService: AuthServiceAPIProtocol,
        authManager: AuthManager,
        appViewModel: AppViewModel
    ) {
        self.authService = authService
        self.authManager = authManager
        self.appViewModel = appViewModel
        self.isAuthenticated = authManager.currentToken() != nil && authManager.currentUserId() != nil
    }

    func login(email: String, password: String) async {
        await authenticate {
            try await authService.login(email: email, password: password)
        }
    }

    func register(email: String, password: String, username: String) async {
        let cleanUsername = Self.normalizedUsername(username)

        await authenticate {
            guard Self.isValidUsername(cleanUsername) else {
                throw APIError.userMessage("Никнейм должен быть от 3 до 30 символов: латинские буквы, цифры, _ или .")
            }

            let matches = try await appViewModel.profileService.search(query: cleanUsername, limit: 10)
            if matches.contains(where: { Self.normalizedUsername($0.username) == cleanUsername }) {
                throw APIError.userMessage("Этот никнейм уже занят")
            }

            try await authService.register(email: email, password: password)
        } afterUserLoaded: { [appViewModel] user in
            let profile = CreateProfileInput(
                name: Self.displayName(from: user.email),
                surname: "",
                username: cleanUsername,
                bio: ""
            ).toModel(userID: user.id)
            _ = try await appViewModel.profileService.create(userID: user.id, model: profile)
            appViewModel.profile.fullName = Self.displayName(from: user.email)
            appViewModel.profile.handle = "@\(cleanUsername)"
        }
    }

    func logout() {
        authManager.clearTokens()
        isAuthenticated = false
    }

    private func authenticate(
        _ tokenRequest: () async throws -> AuthTokensModel,
        afterUserLoaded: ((UserModel) async throws -> Void)? = nil
    ) async {
        isLoading = true
        errorMessage = nil

        do {
            let tokens = try await tokenRequest()
            authManager.saveTokenPair(
                accessToken: tokens.accessToken,
                refreshToken: tokens.refreshToken
            )

            let user = try await authService.getMe()
            authManager.saveSession(
                accessToken: tokens.accessToken,
                refreshToken: tokens.refreshToken,
                userId: user.id
            )
            appViewModel.applyAuthenticatedUser(user)
            try await afterUserLoaded?(user)
            await appViewModel.loadCurrentProfile()
            await appViewModel.loadChats()
            isAuthenticated = true
        } catch {
            authManager.clearTokens()
            errorMessage = error.localizedDescription
        }

        isLoading = false
    }

    private static func displayName(from email: String) -> String {
        email.split(separator: "@").first.map(String.init) ?? email
    }

    private static func normalizedUsername(_ username: String) -> String {
        username
            .trimmingCharacters(in: .whitespacesAndNewlines)
            .trimmingCharacters(in: CharacterSet(charactersIn: "@"))
            .lowercased()
    }

    private static func isValidUsername(_ username: String) -> Bool {
        guard (3...30).contains(username.count) else { return false }
        return username.unicodeScalars.allSatisfy { scalar in
            (CharacterSet.alphanumerics.contains(scalar) && scalar.isASCII)
                || scalar == "_"
                || scalar == "."
        }
    }
}
