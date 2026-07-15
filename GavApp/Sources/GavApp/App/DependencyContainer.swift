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
        MediaURLResolver.configure(socialBaseURL: configuration.socialBaseUrl)
        self.authManager = authManager
        self.authService = AuthServiceAPI(
            baseURL: configuration.socialBaseUrl,
            authManager: authManager
        )

        let chatUseCase = ChatServiceAPI(
            baseURL: configuration.messengerBaseUrl,
            authManager: authManager,
            currentUserIdProvider: { authManager.currentUserId() }
        )
        let profileService = UserProfileServiceAPI(
            baseURL: configuration.socialBaseUrl,
            authManager: authManager
        )
        let uploadService = UploadServiceAPI(
            baseURL: configuration.socialBaseUrl,
            authManager: authManager
        )
        let dogService = DogServiceAPI(
            baseURL: configuration.socialBaseUrl,
            authManager: authManager
        )
        let postService = PostServiceAPI(
            baseURL: configuration.socialBaseUrl,
            authManager: authManager
        )
        let feedService = FeedServiceAPI(
            baseURL: configuration.socialBaseUrl,
            authManager: authManager
        )
        let userService = UserServiceAPI(
            baseURL: configuration.socialBaseUrl,
            authManager: authManager
        )
        let followService = FollowServiceAPI(
            baseURL: configuration.socialBaseUrl,
            authManager: authManager
        )
        let statsService = StatsServiceAPI(
            baseURL: configuration.socialBaseUrl,
            authManager: authManager
        )

        appViewModel = AppViewModel.runtime(
            currentUserId: authManager.currentUserId() ?? UUID(),
            chatUseCase: chatUseCase,
            profileService: profileService,
            uploadService: uploadService,
            dogService: dogService,
            postService: postService,
            feedService: feedService,
            userService: userService,
            followService: followService,
            statsService: statsService
        )
    }
}
