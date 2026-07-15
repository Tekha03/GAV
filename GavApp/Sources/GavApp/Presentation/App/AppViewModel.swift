import Foundation
import SwiftUI
import Combine

enum DogMood: String, CaseIterable, Identifiable {
    case friendly
    case cautious
    case aggressive

    var id: String { rawValue }

    var title: String {
        switch self {
        case .friendly:
            return "Дружелюбный"
        case .cautious:
            return "Настороженный"
        case .aggressive:
            return "Агрессивный"
        }
    }

    var color: Color {
        switch self {
        case .friendly:
            return .green
        case .cautious:
            return .yellow
        case .aggressive:
            return .red
        }
    }
}

struct AppDog: Identifiable, Hashable {
    let id: UUID
    var name: String
    var breed: String
    var ageText: String
    var mood: DogMood
    var photoURL: URL?
    var notes: String
}

struct AppPost: Identifiable, Hashable {
    let id: UUID
    var authorName: String
    var authorHandle: String
    var authorPhotoURL: URL?
    var content: String
    var imageURL: URL?
    var likes: Int
    var comments: Int
    var createdAt: Date
    var isLiked: Bool = false
}

struct AppComment: Identifiable, Hashable {
    let id: UUID
    let postID: UUID
    let authorName: String
    let content: String
    let createdAt: Date
}

public struct AppChat: Identifiable, Hashable {
    public let id: UUID
    public var title: String
    public var lastMessage: String
    public var unreadCount: Int

    public init(id: UUID, title: String, lastMessage: String, unreadCount: Int) {
        self.id = id
        self.title = title
        self.lastMessage = lastMessage
        self.unreadCount = unreadCount
    }
}

extension AppChat {
    var domainChat: Chat {
        Chat(
            id: id,
            isGroup: false,
            title: title,
            photoUrl: "",
            createdAt: .now
        )
    }
}

struct AppVaccination: Identifiable, Hashable {
    let id: UUID
    let dogID: UUID
    var name: String
    var vaccinationDate: Date
    var reminderAfterDays: Int
    var nextDate: Date
    var notes: String
}

struct AppWalker: Identifiable, Hashable {
    let id: UUID
    var ownerName: String
    var dogName: String
    var mood: DogMood
    var latitude: Double
    var longitude: Double
}

struct AppProfile {
    var fullName: String
    var handle: String
    var bio: String
    var avatarURL: URL?
    var followers: Int
    var following: Int

    init(
        fullName: String,
        handle: String,
        bio: String,
        avatarURL: URL?,
        followers: Int,
        following: Int
    ) {
        self.fullName = fullName
        self.handle = handle
        self.bio = bio
        self.avatarURL = avatarURL
        self.followers = followers
        self.following = following
    }
}

@MainActor
final class AppViewModel: ObservableObject {
    @Published var profile: AppProfile
    @Published var dogs: [AppDog]
    @Published var posts: [AppPost]
    @Published var feed: [AppPost]
    @Published var chats: [AppChat]
    @Published var walkers: [AppWalker]
    @Published var vaccinations: [AppVaccination]
    @Published var accounts: [AppProfile] = []
    @Published var likedPostIDs: Set<UUID> = []

    @Published var currentUserId: UUID
    var chatUseCase: ChatUseCase
    var profileService: UserProfileServiceAPIProtocol
    var uploadService: UploadServiceAPIProtocol
    var dogService: DogServiceAPIProtocol
    var postService: PostServiceAPIProtocol
    var feedService: FeedServiceAPIProtocol
    var userService: UserServiceAPIProtocol
    var followService: FollowServiceAPIProtocol
    var statsService: StatsServiceAPIProtocol
    let canEditProfile: Bool

    init(
        profile: AppProfile,
        dogs: [AppDog],
        posts: [AppPost],
        feed: [AppPost],
        chats: [AppChat],
        walkers: [AppWalker],
        vaccinations: [AppVaccination],
        currentUserId: UUID,
        chatUseCase: ChatUseCase,
        profileService: UserProfileServiceAPIProtocol,
        uploadService: UploadServiceAPIProtocol,
        dogService: DogServiceAPIProtocol,
        postService: PostServiceAPIProtocol,
        feedService: FeedServiceAPIProtocol,
        userService: UserServiceAPIProtocol,
        followService: FollowServiceAPIProtocol,
        statsService: StatsServiceAPIProtocol,
        canEditProfile: Bool
    ) {
        self.profile = profile
        self.dogs = dogs
        self.posts = posts
        self.feed = feed
        self.chats = chats
        self.walkers = walkers
        self.vaccinations = vaccinations
        self.currentUserId = currentUserId
        self.chatUseCase = chatUseCase
        self.profileService = profileService
        self.uploadService = uploadService
        self.dogService = dogService
        self.postService = postService
        self.feedService = feedService
        self.userService = userService
        self.followService = followService
        self.statsService = statsService
        self.canEditProfile = canEditProfile
    }

    func toggleLike(postID: UUID) async {
        let wasLiked = likedPostIDs.contains(postID) ||
            posts.first(where: { $0.id == postID })?.isLiked == true ||
            feed.first(where: { $0.id == postID })?.isLiked == true
        let shouldLike = !wasLiked

        setLikeState(postID: postID, isLiked: shouldLike)

        do {
            if shouldLike {
                try await postService.addLike(postID: postID)
            } else {
                try await postService.removeLike(postID: postID)
            }
            _ = await refreshPostStats(postID: postID)
        } catch {
            setLikeState(postID: postID, isLiked: wasLiked)
        }
    }

    func vaccinations(for dogID: UUID) -> [AppVaccination] {
        vaccinations.filter { $0.dogID == dogID }
    }

    func loadChats() async {
        do {
            let userChats = try await chatUseCase.getUserChats(userID: currentUserId)
            var loadedChats: [AppChat] = []

            for chat in userChats {
                let title = await displayTitle(for: chat)
                loadedChats.append(
                    AppChat(
                        id: chat.id,
                        title: title,
                        lastMessage: "Откройте чат",
                        unreadCount: 0
                    )
                )
            }

            chats = loadedChats
        } catch {
            if chats.isEmpty {
                chats = []
            }
        }
    }

    func createPrivateChat(with participantID: UUID) async throws {
        let chat = try await chatUseCase.createPrivateChat(
            user1: currentUserId,
            user2: participantID
        )
        let title = await profileDisplayTitle(userID: participantID) ?? chat.title
        upsertChat(chat, title: title)
    }

    func isFollowing(_ userID: UUID) async -> Bool {
        do {
            let following = try await followService.getFollowing(userID: currentUserId)
            return following.contains { $0.followingId == userID }
        } catch {
            return false
        }
    }

    func follow(_ userID: UUID) async throws {
        profile.following += 1
        do {
            try await followService.follow(userID: userID)
            await loadAuthenticatedStats()
            await loadFeed()
        } catch {
            profile.following = max(0, profile.following - 1)
            throw error
        }
    }

    func unfollow(_ userID: UUID) async throws {
        profile.following = max(0, profile.following - 1)
        do {
            try await followService.unfollow(userID: userID)
            await loadAuthenticatedStats()
            await loadFeed()
        } catch {
            profile.following += 1
            throw error
        }
    }

    func searchProfiles(query: String) async throws -> [UserProfileModel] {
        try await profileService.search(query: query, limit: 10)
            .filter { $0.userId != currentUserId }
    }

    func loadAuthenticatedProfile() async {
        do {
            let profileModel = try await profileService.getByUserID(userID: currentUserId)
            applyProfile(profileModel)
        } catch {
            // Keep the authenticated user fallback if a profile has not been created yet.
        }

        await loadAuthenticatedStats()
    }

    func loadAuthenticatedContent() async {
        async let dogsTask: Void = loadDogs()
        async let postsTask: Void = loadPosts()
        async let feedTask: Void = loadFeed()
        _ = await (dogsTask, postsTask, feedTask)
    }

    func loadProfileContent(for profile: UserProfileModel) async throws -> (dogs: [AppDog], posts: [AppPost]) {
        let dogModels = try await dogService.listByOwnerID(ownerID: profile.userId)
        let postModels = try await postService.listByUser(userID: profile.userId)
        var profilePosts = postModels
            .map { appPost(from: $0, authorProfile: profile) }
            .sorted { $0.createdAt > $1.createdAt }
        await hydratePostStats(&profilePosts)

        return (
            dogs: dogModels.map { appDog(from: $0) },
            posts: profilePosts
        )
    }

    func loadProfileStats(for userID: UUID) async throws -> ProfileStatsModel {
        try await statsService.profileStats(userID: userID)
    }

    func refreshPostStats(postID: UUID) async -> PostStatsModel? {
        guard let stats = try? await statsService.postStats(postID: postID) else {
            return nil
        }

        setPostStats(postID: postID, stats: stats)
        return stats
    }

    func createDog(
        name: String,
        breed: String,
        ageText: String,
        mood: DogMood,
        photoUrl: String,
        notes: String
    ) async throws {
        let model = try await dogService.create(
            ownerID: currentUserId,
            input: CreateDogInput(
                name: name,
                breed: breed,
                age: dogAgeValue(from: ageText),
                status: mood.rawValue,
                gender: "female",
                photoUrl: photoUrl,
                notes: notes
            )
        )
        dogs.insert(appDog(from: model, notes: notes), at: 0)
    }

    func deleteDog(_ dog: AppDog) async throws {
        try await dogService.delete(dogID: dog.id)
        dogs.removeAll { $0.id == dog.id }
        vaccinations.removeAll { $0.dogID == dog.id }
    }

    func updateDog(_ dog: AppDog) async throws {
        try await dogService.update(
            dogID: dog.id,
            input: UpdateDogInput(
                name: dog.name,
                breed: dog.breed,
                age: dogAgeValue(from: dog.ageText),
                status: dog.mood.rawValue,
                gender: "female",
                notes: dog.notes
            )
        )

        if let index = dogs.firstIndex(where: { $0.id == dog.id }) {
            dogs[index] = dog
        }
    }

    func createPost(content: String, imageUrl: String?) async throws {
        let model = try await postService.create(
            userID: currentUserId,
            content: content,
            imageUrl: imageUrl
        )
        let post = appPost(from: model)
        posts.insert(post, at: 0)
        feed.insert(post, at: 0)
    }

    func loadComments(for postID: UUID) async throws -> [AppComment] {
        let models = try await postService.listCommentsByPostID(postID: postID)
        return models.map { appComment(from: $0) }
    }

    func addComment(to postID: UUID, content: String) async throws -> [AppComment] {
        _ = try await postService.createComment(
            postID: postID,
            userID: currentUserId,
            content: content
        )
        let comments = try await loadComments(for: postID)
        setCommentCount(postID: postID, count: comments.count)
        return comments
    }

    func setCommentCount(postID: UUID, count: Int) {
        let normalizedCount = max(0, count)
        if let index = posts.firstIndex(where: { $0.id == postID }) {
            posts[index].comments = normalizedCount
        }
        if let index = feed.firstIndex(where: { $0.id == postID }) {
            feed[index].comments = normalizedCount
        }
    }

    func shareLocationAndLoadNearby(
        latitude: Double,
        longitude: Double,
        radiusMeters: Double = 1_000
    ) async throws {
        try await userService.updateLocation(
            id: currentUserId,
            input: UpdateLocationInput(
                lat: latitude,
                lon: longitude,
                locationStatus: .walking,
                visibility: .everyone
            )
        )

        let nearbyDogs = try await userService.findDogsNearby(
            id: currentUserId,
            centerLat: latitude,
            centerLon: longitude,
            radiusMeters: radiusMeters
        )

        walkers = nearbyDogs.compactMap { model in
            guard let lat = model.lat, let lon = model.lon else { return nil }
            return AppWalker(
                id: model.id,
                ownerName: "Поблизости",
                dogName: model.name,
                mood: DogMood(rawValue: model.status) ?? .friendly,
                latitude: lat,
                longitude: lon
            )
        }
    }

    func createGroupChat(title: String, memberIDs: [UUID]) async throws {
        let chat = try await chatUseCase.createGroupChat(
            title: title,
            creator: currentUserId,
            members: memberIDs
        )
        upsertChat(chat)
    }

    func applyAuthenticatedUser(_ user: UserModel) {
        currentUserId = user.id
        let username = handle(from: user.email).trimmingCharacters(in: CharacterSet(charactersIn: "@"))
        profile = AppProfile(
            fullName: displayName(from: user.email),
            handle: "@\(username)",
            bio: "",
            avatarURL: nil,
            followers: 0,
            following: 0
        )
        dogs = []
        posts = []
        feed = []
        chats = []
        walkers = []
        vaccinations = []
        likedPostIDs = []
    }

    func applySavedSession(userID: UUID) {
        currentUserId = userID
        profile = AppProfile(
            fullName: "Мой профиль",
            handle: "@me",
            bio: "",
            avatarURL: nil,
            followers: 0,
            following: 0
        )
        dogs = []
        posts = []
        feed = []
        walkers = []
        vaccinations = []
        likedPostIDs = []
    }

    private func applyProfile(_ model: UserProfileModel) {
        let fullName = [model.name, model.surname]
            .map { $0.trimmingCharacters(in: .whitespacesAndNewlines) }
            .filter { !$0.isEmpty }
            .joined(separator: " ")
        let username = model.username.trimmingCharacters(in: .whitespacesAndNewlines)

        profile = AppProfile(
            fullName: fullName.isEmpty ? "Мой профиль" : fullName,
            handle: username.isEmpty ? "@me" : "@\(username)",
            bio: model.bio,
            avatarURL: model.profilePhotoUrl.flatMap { MediaURLResolver.resolve($0) },
            followers: profile.followers,
            following: profile.following
        )
    }

    private func loadDogs() async {
        do {
            let models = try await dogService.listByOwnerID(ownerID: currentUserId)
            dogs = models.map { appDog(from: $0) }
        } catch {
            if dogs.isEmpty {
                dogs = []
            }
        }
    }

    private func loadPosts() async {
        do {
            let models = try await postService.listByUser(userID: currentUserId)
            var loadedPosts = models
                .map { appPost(from: $0) }
                .sorted { $0.createdAt > $1.createdAt }
            await hydratePostStats(&loadedPosts)
            posts = loadedPosts
        } catch {
            if posts.isEmpty {
                posts = []
            }
        }
    }

    private func loadAuthenticatedStats() async {
        do {
            let stats = try await statsService.profileStats(userID: currentUserId)
            profile.followers = Int(stats.followersCount)
            profile.following = Int(stats.followingsCount)
        } catch {
            // Keep local counters if stats are not available yet.
        }
    }

    private func loadFeed() async {
        do {
            let models = try await feedService.getFeed(userID: currentUserId, before: nil, limit: 50)
            var loadedFeed = models
                .map { appPost(from: $0) }
                .sorted { $0.createdAt > $1.createdAt }
            await hydratePostAuthors(&loadedFeed)
            await hydratePostStats(&loadedFeed)
            feed = loadedFeed
        } catch {
            if feed.isEmpty {
                feed = posts
            }
        }
    }

    private func appDog(from model: DogModel, notes: String = "") -> AppDog {
        AppDog(
            id: model.id,
            name: model.name,
            breed: model.breed,
            ageText: dogAgeText(from: model.age),
            mood: DogMood(rawValue: model.status) ?? .friendly,
            photoURL: MediaURLResolver.resolve(model.photoUrl),
            notes: notes.isEmpty ? model.notes : notes
        )
    }

    private func appPost(from model: PostModel, authorProfile: UserProfileModel? = nil) -> AppPost {
        let authorName: String
        let authorHandle: String
        let authorPhotoURL: URL?

        if let authorProfile {
            let fullName = [authorProfile.name, authorProfile.surname]
                .map { $0.trimmingCharacters(in: .whitespacesAndNewlines) }
                .filter { !$0.isEmpty }
                .joined(separator: " ")
            authorName = fullName.isEmpty ? "@\(authorProfile.username)" : fullName
            authorHandle = "@\(authorProfile.username)"
            authorPhotoURL = authorProfile.profilePhotoUrl.flatMap { MediaURLResolver.resolve($0) }
        } else if model.userId == currentUserId {
            authorName = profile.fullName
            authorHandle = profile.handle
            authorPhotoURL = profile.avatarURL
        } else {
            authorName = "Пользователь"
            authorHandle = "@user"
            authorPhotoURL = nil
        }

        return AppPost(
            id: model.id,
            authorName: authorName,
            authorHandle: authorHandle,
            authorPhotoURL: authorPhotoURL,
            content: model.content,
            imageURL: model.imageUrl.flatMap { MediaURLResolver.resolve($0) },
            likes: 0,
            comments: 0,
            createdAt: model.createdAt
        )
    }

    private func hydratePostStats(_ posts: inout [AppPost]) async {
        for index in posts.indices {
            if let stats = try? await statsService.postStats(postID: posts[index].id) {
                posts[index].likes = Int(stats.likesCount)
                posts[index].comments = Int(stats.commentsCount)
            }
        }
    }

    private func hydratePostAuthors(_ posts: inout [AppPost]) async {
        var profilesByUserID: [UUID: UserProfileModel] = [:]

        for post in posts where post.authorHandle == "@user" {
            guard let model = try? await postService.getByID(id: post.id) else { continue }
            if profilesByUserID[model.userId] == nil,
               let profile = try? await profileService.getByUserID(userID: model.userId) {
                profilesByUserID[model.userId] = profile
            }
        }

        for index in posts.indices where posts[index].authorHandle == "@user" {
            guard let model = try? await postService.getByID(id: posts[index].id),
                  let authorProfile = profilesByUserID[model.userId] else { continue }
            posts[index] = appPost(from: model, authorProfile: authorProfile)
        }
    }

    private func appComment(from model: CommentModel) -> AppComment {
        AppComment(
            id: model.id,
            postID: model.postId,
            authorName: model.userId == currentUserId ? profile.fullName : "Пользователь",
            content: model.content,
            createdAt: model.createdAt
        )
    }

    private func setLikeState(postID: UUID, isLiked: Bool) {
        if isLiked {
            likedPostIDs.insert(postID)
        } else {
            likedPostIDs.remove(postID)
        }

        if let index = posts.firstIndex(where: { $0.id == postID }) {
            let delta = likeDelta(current: posts[index].isLiked, next: isLiked)
            posts[index].isLiked = isLiked
            posts[index].likes = max(0, posts[index].likes + delta)
        }

        if let index = feed.firstIndex(where: { $0.id == postID }) {
            let delta = likeDelta(current: feed[index].isLiked, next: isLiked)
            feed[index].isLiked = isLiked
            feed[index].likes = max(0, feed[index].likes + delta)
        }
    }

    private func setPostStats(postID: UUID, stats: PostStatsModel) {
        if let index = posts.firstIndex(where: { $0.id == postID }) {
            posts[index].likes = Int(stats.likesCount)
            posts[index].comments = Int(stats.commentsCount)
        }

        if let index = feed.firstIndex(where: { $0.id == postID }) {
            feed[index].likes = Int(stats.likesCount)
            feed[index].comments = Int(stats.commentsCount)
        }
    }

    private func likeDelta(current: Bool, next: Bool) -> Int {
        switch (current, next) {
        case (false, true):
            return 1
        case (true, false):
            return -1
        default:
            return 0
        }
    }

    private func dogAgeValue(from text: String) -> String {
        let clean = text.trimmingCharacters(in: .whitespacesAndNewlines).lowercased()
        if clean.contains("щен") || clean.contains("puppy") {
            return "puppy"
        }
        if clean.contains("стар") || clean.contains("elder") || clean.contains("senior") {
            return "elderly"
        }
        return "adult"
    }

    private func dogAgeText(from value: String) -> String {
        switch value {
        case "puppy":
            return "Щенок"
        case "elderly", "senior":
            return "Пожилая"
        default:
            return "Взрослая"
        }
    }

    private func displayName(from email: String) -> String {
        email.split(separator: "@").first.map(String.init) ?? email
    }

    private func handle(from email: String) -> String {
        let prefix = displayName(from: email)
            .lowercased()
            .filter { $0.isLetter || $0.isNumber || $0 == "_" }

        return prefix.isEmpty ? "@me" : "@\(prefix)"
    }

    private func upsertChat(_ chat: Chat) {
        upsertChat(chat, title: chat.title)
    }

    private func upsertChat(_ chat: Chat, title: String) {
        let item = AppChat(
            id: chat.id,
            title: title,
            lastMessage: "Откройте чат",
            unreadCount: 0
        )

        if let index = chats.firstIndex(where: { $0.id == chat.id }) {
            chats[index] = item
        } else {
            chats.insert(item, at: 0)
        }
    }

    private func displayTitle(for chat: Chat) async -> String {
        if chat.isGroup {
            return chat.title.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty ? "Группа" : chat.title
        }

        guard let members = try? await chatUseCase.getChatMembers(chatID: chat.id),
              let otherUserID = members.first(where: { $0.userId != currentUserId })?.userId,
              let title = await profileDisplayTitle(userID: otherUserID) else {
            return chat.title.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty ? "Чат" : chat.title
        }

        return title
    }

    private func profileDisplayTitle(userID: UUID) async -> String? {
        guard let profile = try? await profileService.getByUserID(userID: userID) else {
            return nil
        }

        let fullName = [profile.name, profile.surname]
            .map { $0.trimmingCharacters(in: .whitespacesAndNewlines) }
            .filter { !$0.isEmpty }
            .joined(separator: " ")

        if !fullName.isEmpty {
            return fullName
        }

        let username = profile.username.trimmingCharacters(in: .whitespacesAndNewlines)
        return username.isEmpty ? nil : "@\(username)"
    }
}

extension AppViewModel {
    static func runtime(
        currentUserId: UUID,
        chatUseCase: ChatUseCase,
        profileService: UserProfileServiceAPIProtocol,
        uploadService: UploadServiceAPIProtocol,
        dogService: DogServiceAPIProtocol,
        postService: PostServiceAPIProtocol,
        feedService: FeedServiceAPIProtocol,
        userService: UserServiceAPIProtocol,
        followService: FollowServiceAPIProtocol,
        statsService: StatsServiceAPIProtocol
    ) -> AppViewModel {
        AppViewModel(
            profile: AppProfile(
                fullName: "",
                handle: "",
                bio: "",
                avatarURL: nil,
                followers: 0,
                following: 0
            ),
            dogs: [],
            posts: [],
            feed: [],
            chats: [],
            walkers: [],
            vaccinations: [],
            currentUserId: currentUserId,
            chatUseCase: chatUseCase,
            profileService: profileService,
            uploadService: uploadService,
            dogService: dogService,
            postService: postService,
            feedService: feedService,
            userService: userService,
            followService: followService,
            statsService: statsService,
            canEditProfile: true
        )
    }
}