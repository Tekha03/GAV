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

        let authService = AuthServiceAPI(
            baseURL: configuration.socialBaseURL,
            authManager: authManager,
        )

        let chatService = ChatServiceAPI(
            baseURL: configuration.messengerBaseURL,
            authManager: authManager,
            currentUserIdProvider: { authManager.currentUserID() }
        )

        let profileService = UserProfileServiceAPI(
            baseURL: configuration.socialBaseURL,
            authManager: authManager
        )

        let uploadService = UploadServiceAPI(
            baseURL: configuration.socialBaseURL,
            authManager: authManager
        )

        let dogService = DogServiceAPI(
            baseURL: configuration.socialBaseURL,
            authManager: authManager
        )

        let postService = PostServiceAPI(
            baseURL: configuration.socialBaseURL,
            authManager: authManager
        )

        let feedService = FeedServiceAPI(
            baseURL: configuration.socialBaseURL,
            authManager: authManager
        )

        let userService = UserServiceAPI(
            baseURL: configuration.socialBaseURL,
            authManager: authManager
        )

        let followService = FollowServiceAPI(
            baseURL: configuration.socialBaseURL,
            authManager: authManager
        )

        let statsService = StatsServiceAPI(
            baseURL: configuration.socialBaseURL,
            authManager: authManager
        )

        self.authManager = authManager
        self.authService = authService

        self.AppViewModel = AppViewModel.runtime(
            currentUserId: authManager.currentUserId(),
            chatUseCase: chatService,
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
