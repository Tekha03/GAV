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
        let hasSavedSession = authManager.currentToken() != nil && authManager.currentUserId() != nil
        self.isAuthenticated = hasSavedSession
        self.isLoading = hasSavedSession
    }

    func login(email: String, password: String) async {
        await authenticate {
            try await authService.login(email: email, password: password)
        }
    }

    func register(email: String, password: String, firstName: String, lastName: String, username: String) async {
        await authenticate {
            try await authService.register(email: email, password: password)
        } afterUserLoaded: { [appViewModel] user in
            let cleanUsername = Self.normalizedUsername(username)
            let profile = CreateProfileInput(
                name: firstName.trimmingCharacters(in: .whitespacesAndNewlines),
                surname: lastName.trimmingCharacters(in: .whitespacesAndNewlines),
                username: cleanUsername,
                bio: ""
            ).toModel(userID: user.id)
            _ = try await appViewModel.profileService.create(userID: user.id, model: profile)
            appViewModel.profile.fullName = [profile.name, profile.surname]
                .map { $0.trimmingCharacters(in: .whitespacesAndNewlines) }
                .filter { !$0.isEmpty }
                .joined(separator: " ")
            appViewModel.profile.handle = "@\(cleanUsername)"
        }
    }

    func logout() {
        authManager.clearTokens()
        isAuthenticated = false
    }

    func restoreSavedSessionIfNeeded() async {
        guard authManager.currentToken() != nil, authManager.currentUserId() != nil else {
            isLoading = false
            isAuthenticated = false
            return
        }

        isLoading = true
        errorMessage = nil

        do {
            let user = try await authService.getMe()
            authManager.saveSession(
                accessToken: authManager.getAccessToken() ?? "",
                refreshToken: authManager.getRefreshToken() ?? "",
                userId: user.id
            )
            appViewModel.applyAuthenticatedUser(user)
            await appViewModel.loadAuthenticatedProfile()
            await appViewModel.loadAuthenticatedContent()
            await appViewModel.loadChats()
            isAuthenticated = true
        } catch {
            authManager.clearTokens()
            isAuthenticated = false
        }

        isLoading = false
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
            await appViewModel.loadAuthenticatedProfile()
            await appViewModel.loadAuthenticatedContent()
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
}
