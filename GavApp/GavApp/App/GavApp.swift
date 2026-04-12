// App/GavApp.swift
@main
struct GavApp: App {
    private let container = DependencyContainer()

    var body: some Scene {
        WindowGroup {
            HomeView()
                .environmentObject(container.viewModel)
        }
    }

    final class DependencyContainer {
        let baseURL: URL = {
            #if DEBUG
                return URL(string: "http://localhost:8080")!
            #else
                return URL(string: "https://api.gavapp.example")!
            #endif
        }()

        let session = URLSession.shared
        let authManager: AuthManager

        // MARK: - API

        let authServiceAPI: AuthServiceAPIProtocol
        let userServiceAPI: UserServiceAPIProtocol
        let profileServiceAPI: UserProfileServiceAPIProtocol
        let postServiceAPI: PostServiceAPIProtocol
        let feedServiceAPI: FeedServiceAPIProtocol
        let dogServiceAPI: DogServiceAPIProtocol
        let vaccinationServiceAPI: VaccinationServiceAPIProtocol
        let followServiceAPI: FollowServiceAPIProtocol
        let statsServiceAPI: StatsServiceAPIProtocol
        let settingsServiceAPI: SettingsServiceAPIProtocol
        let uploadServiceAPI: UploadServiceAPIProtocol

        // MARK: - Repositories

        let authRepository: any AuthRepository
        let userRepository: any UserRepository
        let profileRepository: any UserProfileRepository
        let postRepository: any PostRepository
        let feedRepository: any FeedRepository
        let dogRepository: any DogRepository
        let vaccinationRepository: any VaccinationRepository
        let followRepository: any FollowRepository
        let statsRepository: any StatsRepository
        let settingsRepository: any SettingsRepository
        let uploadRepository: any UploadRepository

        // MARK: - UseCases

        let authUseCase: AuthUseCase
        let userUseCase: UserUseCase
        let profileUseCase: ProfileUseCase
        let postUseCase: PostUseCase
        let feedUseCase: FeedUseCase
        let dogUseCase: DogUseCase
        let vaccinationUseCase: VaccinationUseCase
        let followUseCase: FollowUseCase
        let statsUseCase: StatsUseCase
        let settingsUseCase: SettingsUseCase
        let uploadUseCase: UploadUseCase

        // MARK: - ViewModels

        let viewModel: ProfileViewModel

        init() {
            authManager = AuthManager()

            // API‑клиенты
            authServiceAPI = AuthServiceAPI(
                baseURL: baseURL,
                session: session,
                authManager: authManager
            )
            dogServiceAPI = DogServiceAPI(
                baseURL: baseURL,
                session: session,
                authManager: authManager
            )
            // и т.д. для всех ...ServiceAPI

            // Репозитории
            authRepository = AuthRepositoryImpl(api: authServiceAPI)
            dogRepository = DogRepositoryImpl(api: dogServiceAPI)
            // и т.д.

            // UseCase
            authUseCase = AuthUseCase(repository: authRepository)
            dogUseCase = DogUseCase(repository: dogRepository)
            // и т.д.

            // ViewModels
            viewModel = ProfileViewModel(
                profileUseCase: profileUseCase,
                dogUseCase: dogUseCase,
                postUseCase: feedUseCase,
                statsUseCase: statsUseCase
            )
        }
    }
}