import Foundation

@MainActor
@available(macOS 12.0, *)
final class AppViewModel: ObservableObject {
    @Published private(set) var currentUser: User?

    let profileUseCase: ProfileUseCase
    let feedUseCase: FeedUseCase
    let followUseCase: FollowUseCase
    let dogUseCase: DogUseCase
    let vaccinationUseCase: VaccinationUseCase
    let chatUseCase: ChatUseCase
    let mapUseCase: MapUseCase
    let settingsUseCase: SettingsUseCase
    let authUseCase: AuthUseCase

    init(
        currentUser: User,
        profileUseCase: ProfileUseCase,
        feedUseCase: FeedUseCase,
        followUseCase: FollowUseCase,
        dogUseCase: DogUseCase,
        vaccinationUseCase: VaccinationUseCase,
        chatUseCase: ChatUseCase,
        mapUseCase: MapUseCase,
        settingsUseCase: SettingsUseCase,
        authUseCase: AuthUseCase
    ) {
        self.currentUser = currentUser
        self.profileUseCase = profileUseCase
        self.feedUseCase = feedUseCase
        self.followUseCase = followUseCase
        self.dogUseCase = dogUseCase
        self.vaccinationUseCase = vaccinationUseCase
        self.chatUseCase = chatUseCase
        self.mapUseCase = mapUseCase
        self.settingsUseCase = settingsUseCase
        self.authUseCase = authUseCase
    }

    func login(email: String, password: String) async throws {
        let tokens = try await authUseCase.login(email: email, password: password)
        authManager.saveTokens(tokens)
        let user = try await profileUseCase.getByUserID(userID: tokens.userId)
        currentUser = user
    }

    func logout() async throws {
        try await authUseCase.logout()
        currentUser = nil
    }
}