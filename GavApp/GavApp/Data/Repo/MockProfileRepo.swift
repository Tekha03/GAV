import Foundation

final class MockProfileRepository: ProfileRepository {

    func fetchProfile(userId: UUID) async -> UserProfile {
        try? await Task.sleep(nanoseconds: 300_000_000)

        return UserProfile(
            userId: userId,
            name: "Виктория",
            surname: "Кашуркина",
            username: "@duduka",
            bio: "Собаки > люди",
            profilePhotoUrl: "vick",
            followersCount: 100,
            followingCount: 50,
            isFollowed: true
        )
    }

    func followUser(userId: UUID) async -> Bool {
        true
    }
}
