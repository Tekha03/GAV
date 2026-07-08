import Foundation

@MainActor
final class DependencyContainer {
    let appViewModel: AppViewModel
    let authManager: AuthManager
    let authService: AuthServiceAPIProtocol

    private let socialBaseURL = URL(string: "http://192.168.1.4:8080")!
    private let messengerBaseURL = URL(string: "http://192.168.1.4:8082")!

    init() {
        let authManager = AuthManager()
        MediaURLResolver.configure(socialBaseURL: socialBaseURL)
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
        appViewModel.dogService = DogServiceAPI(
            baseURL: socialBaseURL,
            authManager: authManager
        )
        appViewModel.postService = PostServiceAPI(
            baseURL: socialBaseURL,
            authManager: authManager
        )
        appViewModel.feedService = FeedServiceAPI(
            baseURL: socialBaseURL,
            authManager: authManager
        )
        appViewModel.userService = UserServiceAPI(
            baseURL: socialBaseURL,
            authManager: authManager
        )
        appViewModel.followService = FollowServiceAPI(
            baseURL: socialBaseURL,
            authManager: authManager
        )
        appViewModel.statsService = StatsServiceAPI(
            baseURL: socialBaseURL,
            authManager: authManager
        )
    }
}
