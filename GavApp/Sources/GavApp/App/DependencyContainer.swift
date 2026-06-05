import Foundation

@MainActor
final class DependencyContainer {
    let appViewModel: AppViewModel
    let authManager: AuthManager
    let authService: AuthServiceAPIProtocol

    private let socialBaseURL = URL(string: "http://localhost:8080")!
    private let messengerBaseURL = URL(string: "http://localhost:8082")!

    init() {
        let authManager = AuthManager()
        self.authManager = authManager
        self.authService = AuthServiceAPI(baseURL: socialBaseURL, authManager: authManager)

        appViewModel = .preview
        if let currentUserId = authManager.currentUserId() {
            appViewModel.applySavedSession(userID: currentUserId)
        }
        appViewModel.chatUseCase = ChatServiceAPI(
            baseURL: messengerBaseURL,
            authManager: authManager,
            currentUserIdProvider: { authManager.currentUserId() }
        )
        appViewModel.profileService = UserProfileServiceAPI(
            baseURL: socialBaseURL,
            authManager: authManager
        )
        appViewModel.uploadService = UploadServiceAPI(
            baseURL: socialBaseURL,
            authManager: authManager
        )
    }
}
