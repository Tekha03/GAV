import Foundation
import SwiftUI
import Combine

@MainActor
final class ProfileViewModel: ObservableObject {
    @Published var profile: AppProfile
    @Published var currentUser: UserModel?
    @Published var currentUserId: UUID?
    @Published var stats: ProfileStatsModel?
    @Published var isLoading = false
    @Published var errorMessage: String?

    private let authManager: AuthManager
    private let userProfileService: UserProfileServiceAPIProtocol
    private let settingsService: SettingsServiceAPIProtocol
    private let dogService: DogServiceAPIProtocol
    private let statsService: StatsServiceAPIProtocol

    init(
        authManager: AuthManager,
        userProfileService: UserProfileServiceAPIProtocol,
        settingsService: SettingsServiceAPIProtocol,
        dogService: DogServiceAPIProtocol,
        statsService: StatsServiceAPIProtocol,
        previewProfile: AppProfile? = nil
    ) {
        self.authManager = authManager
        self.userProfileService = userProfileService
        self.settingsService = settingsService
        self.dogService = dogService
        self.statsService = statsService
        self.profile = previewProfile ?? AppProfile(
            fullName: "Viktoria Kashurkina",
            handle: "vickdogmom",
            bio: "dog-friendly",
            avatarURL: nil,
            followers: 0,
            following: 0
        )
        self.currentUserId = authManager.currentUserId()
    }

    func loadProfile() async {
        guard let userID = authManager.currentUserId() else {
            errorMessage = "No user ID"
            return
        }

        currentUserId = userID
        isLoading = true
        defer { isLoading = false }

        let profileService = userProfileService
        let statsService = statsService

        do {
            let profileModel = try await profileService.getByUserID(userID: userID)
            let statsModel = try await statsService.profileStats(userID: userID)

            self.stats = statsModel
            self.profile = AppProfile(
                fullName: "\(profileModel.name) \(profileModel.surname)",
                handle: profileModel.username,
                bio: profileModel.bio,
                avatarURL: profileModel.profilePhotoUrl.flatMap(URL.init(string:)),
                followers: Int(statsModel.followersCount),
                following: Int(statsModel.followingsCount)
            )
        } catch {
            errorMessage = error.localizedDescription
        }
    }

    func updateProfile(input: UpdateProfileInput) async {
        guard let userID = currentUserId ?? authManager.currentUserId() else { return }

        let profileService = userProfileService

        do {
            try await profileService.update(userID: userID, input: input)
            await loadProfile()
        } catch {
            errorMessage = error.localizedDescription
        }
    }

    func updateUserSettings(settings: UserSettingsModel) async {
        let settingsService = settingsService

        do {
            let input = UpdateUserSettingsInput(
                profilePrivacy: settings.profilePrivacy,
                showLocation: settings.showLocation,
                allowMessages: settings.allowMessages
            )
            try await settingsService.updateSettings(input: input)
        } catch {
            errorMessage = error.localizedDescription
        }
    }

    func addDog(input: CreateDogInput) async {
        guard let userID = currentUserId ?? authManager.currentUserId() else { return }

        let dogService = dogService

        do {
            _ = try await dogService.create(ownerID: userID, input: input)
            await loadProfile()
        } catch {
            errorMessage = error.localizedDescription
        }
    }
}
