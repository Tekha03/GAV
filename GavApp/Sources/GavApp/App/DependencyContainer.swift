import Foundation

@MainActor
final class DependencyContainer {
    let appViewModel: AppViewModel
    let authManager: AuthManager
    let authService: AuthServiceAPIProtocol
    let configuration: AppConfiguration

    init(configuration: AppConfiguration = .current) {
        self.configuration = configuration

        let authManager = AuthManager()
        MediaURLResolver.configure(socialBaseURL: configuration.socialBaseURL)
        self.authManager = authManager
        self.authService = AuthServiceAPI(
            baseURL: configuration.socialBaseURL,
            authManager: authManager,
        )

        appViewModel = .preview
        if let currentUserId = authManager.currentUserId() {
            appViewModel.applySavedSession(userID: currentUserId)
        }

        appViewModel.chatUseCase = ChatServiceAPI(
            baseURL: configuration.messengerBaseURL,
            authManager: authManager,
            currentUserIdProvider: { authManager.currentUserId() }
        )

        appViewModel.profileService = UserProfileServiceAPI(
            baseURL: configuration.socialBaseURL,
            authManager: authManager
        )

        appViewModel.uploadService = UploadServiceAPI(
            baseURL: configuration.socialBaseURL,
            authManager: authManager
        )

        appViewModel.dogService = DogServiceAPI(
            baseURL: configuration.socialBaseURL,
            authManager: authManager
        )

        appViewModel.postService = PostServiceAPI(
            baseURL: configuration.socialBaseURL,
            authManager: authManager
        )

        appViewModel.feedService = FeedServiceAPI(
            baseURL: configuration.socialBaseURL,
            authManager: authManager
        )

        appViewModel.userService = UserServiceAPI(
            baseURL: configuration.socialBaseURL,
            authManager: authManager
        )

        appViewModel.followService = FollowServiceAPI(
            baseURL: configuration.socialBaseURL,
            authManager: authManager
        )

        appViewModel.statsService = StatsServiceAPI(
            baseURL: configuration.socialBaseURL,
            authManager: authManager
        )
    }
}
