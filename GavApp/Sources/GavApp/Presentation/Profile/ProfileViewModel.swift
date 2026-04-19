import Foundation
import SwiftUI

@MainActor
@available(macOS 12.0, *)
final class ProfileViewModel: ObservableObject {

    @Published var profile: UserProfile?
    @Published var dogs: [Dog] = []
    @Published var posts: [Post] = []
    @Published var isLoading = false
    @Published private(set) var isOwner = false

    private let profileUseCase: ProfileUseCase
    private let dogUseCase: DogUseCase
    private let followUseCase: FollowUseCase
    private let feedUseCase: FeedUseCase
    private let statsUseCase: StatsUseCase
    private let userUseCase: UserUseCase
    private let postUseCase: PostUseCase
    private let likeUseCase: LikeUseCase

    private let currentUserId: UUID

    init(
        currentUserId: UUID,
        profileUseCase: ProfileUseCase,
        dogUseCase: DogUseCase,
        followUseCase: FollowUseCase,
        feedUseCase: FeedUseCase,
        statsUseCase: StatsUseCase,
        userUseCase: UserUseCase,
        postUseCase: PostUseCase,
        likeUseCase: LikeUseCase
    ) {
        self.currentUserId = currentUserId
        self.profileUseCase = profileUseCase
        self.dogUseCase = dogUseCase
        self.followUseCase = followUseCase
        self.feedUseCase = feedUseCase
        self.statsUseCase = statsUseCase
        self.userUseCase = userUseCase
        self.postUseCase = postUseCase
        self.likeUseCase = likeUseCase
    }

    // MARK: Load Profile

    func load(userId: UUID) async {
        isLoading = true
        defer { isLoading = false }

        do {
            let profile = try await profileUseCase.getByUserID(userID: userId)
            let dogs = try await dogUseCase.listByOwnerID(ownerID: userId)
            let posts = try await feedUseCase.getFeed(
                userID: userId,
                before: nil,
                limit: 30
            )

            _ = try await statsUseCase.profileStats(userID: userId)

            self.profile = profile
            self.dogs = dogs
            self.posts = posts
            self.isOwner = userId == currentUserId

        } catch {
            print("Load profile error: \(error)")
        }
    }

    // MARK: Follow

    func toggleFollow() async {
        guard let id = profile?.userId else { return }

        let follow = Follow(
            followerId: currentUserId,
            followingId: id
        )

        do {
            let exists = try await followUseCase.exists(follow: follow)

            if exists {
                try await followUseCase.unfollow(follow: follow)
                profile?.isFollowed = false
                profile?.followersCount -= 1
            } else {
                try await followUseCase.follow(follow: follow)
                profile?.isFollowed = true
                profile?.followersCount += 1
            }

        } catch {
            print("Follow error: \(error)")
        }
    }

    // MARK: Like Post

    func toggleLike(for post: Post) async {

        let like = Like(
            userId: currentUserId,
            postId: post.id
        )

        do {
            let exists = try await likeUseCase.exists(like: like)

            if let index = posts.firstIndex(where: { $0.id == post.id }) {

                if exists {
                    try await likeUseCase.remove(like: like)
                    posts[index].likesCount -= 1
                } else {
                    try await likeUseCase.add(like: like)
                    posts[index].likesCount += 1
                }
            }

        } catch {
            print("Like error: \(error)")
        }
    }

    // MARK: Create Post

    func createPost(
        content: String,
        imageUrl: String? = nil
    ) async {

        guard let userId = profile?.userId else { return }

        do {
            let post = try await postUseCase.create(
                userID: userId,
                content: content,
                imageUrl: imageUrl
            )

            posts.insert(post, at: 0)

        } catch {
            print("Create post error: \(error)")
        }
    }

    // MARK: Profile Update

    func updateProfile(input: UpdateProfileInput) async {
        do {
            try await profileUseCase.update(
                userID: currentUserId,
                input: input
            )

            profile?.name = input.name
            profile?.surname = input.surname
            profile?.bio = input.bio

        } catch {
            print("Update profile error: \(error)")
        }
    }

    // MARK: Dog

    func createDog(_ input: CreateDogInput) async {
        do {
            let dog = try await dogUseCase.create(
                userId: currentUserId,
                input: input
            )

            dogs.append(dog)

        } catch {
            print("Create dog error: \(error)")
        }
    }

    func updateDog(
        dogId: UUID,
        input: UpdateDogInput
    ) async {

        do {
            try await dogUseCase.update(
                dogId: dogId,
                input: input
            )
        } catch {
            print("Update dog error: \(error)")
        }
    }

    // MARK: Settings

    func updateUserSettings(
        settings: UserSettings
    ) async {

        do {
            try await userUseCase.updateUserSettings(
                id: currentUserId,
                input: settings
            )
        } catch {
            print("Settings update error: \(error)")
        }
    }
}