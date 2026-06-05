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
        self.canEditProfile = canEditProfile
    }

    func addDog() {
        dogs.insert(
            AppDog(
                id: UUID(),
                name: "Новая собака",
                breed: "Уточнить породу",
                ageText: "1 год",
                mood: .friendly,
                photoURL: nil,
                notes: "Карточка создана. Следующий шаг: привязать к `POST /api/v1/dogs/`."
            ),
            at: 0
        )
    }

    func addPost() {
        let post = AppPost(
            id: UUID(),
            authorName: profile.fullName,
            authorHandle: profile.handle,
            authorPhotoURL: profile.avatarURL,
            content: "Новый пост хозяина. Следом можно связать с `POST /api/v1/posts/` и загрузкой `POST /api/v1/upload/post-image`.",
            imageURL: nil,
            likes: 0,
            comments: 0,
            createdAt: .now
        )
        posts.insert(post, at: 0)
        feed.insert(post, at: 0)
    }

    func toggleLike(postID: UUID) {
        if let index = posts.firstIndex(where: { $0.id == postID }) {
            posts[index].isLiked.toggle()
            posts[index].likes += posts[index].isLiked ? 1 : -1
            posts[index].likes = max(0, posts[index].likes)
        }

        if let index = feed.firstIndex(where: { $0.id == postID }) {
            feed[index].isLiked.toggle()
            feed[index].likes += feed[index].isLiked ? 1 : -1
            feed[index].likes = max(0, feed[index].likes)
        }
    }

    func vaccinations(for dogID: UUID) -> [AppVaccination] {
        vaccinations.filter { $0.dogID == dogID }
    }

    func loadChats() async {
        do {
            let userChats = try await chatUseCase.getUserChats(userID: currentUserId)
            chats = userChats.map {
                AppChat(
                    id: $0.id,
                    title: $0.title,
                    lastMessage: "Откройте чат",
                    unreadCount: 0
                )
            }
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
        upsertChat(chat)
    }

    func searchProfiles(query: String) async throws -> [UserProfileModel] {
        try await profileService.search(query: query, limit: 10)
            .filter { $0.userId != currentUserId }
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
        let item = AppChat(
            id: chat.id,
            title: chat.title,
            lastMessage: "Откройте чат",
            unreadCount: 0
        )

        if let index = chats.firstIndex(where: { $0.id == chat.id }) {
            chats[index] = item
        } else {
            chats.insert(item, at: 0)
        }
    }
}

extension AppViewModel {
    static let preview: AppViewModel = {
        let dog1 = AppDog(
            id: UUID(),
            name: "Бимка",
            breed: "Лабрадор",
            ageText: "3 года",
            mood: .friendly,
            photoURL: URL(string: "https://images.unsplash.com/photo-1518717758536-85ae29035b6d"),
            notes: "Любит утренние прогулки и не конфликтует с другими собаками."
        )
        let dog2 = AppDog(
            id: UUID(),
            name: "Рокси",
            breed: "Корги",
            ageText: "8 месяцев",
            mood: .cautious,
            photoURL: URL(string: "https://images.unsplash.com/photo-1548199973-03cce0bbc87b"),
            notes: "Пока настороженно знакомится с новыми маршрутами."
        )

        let ownPost = AppPost(
            id: UUID(),
            authorName: "Вика Кашуркина",
            authorHandle: "@vickdogmom",
            authorPhotoURL: URL(string: ""),
            content: "Профиль должен быть как в инсте: фото, био, сторисы собак и посты хозяина внизу.",
            imageURL: URL(string: ""),
            likes: 128,
            comments: 14,
            createdAt: .now.addingTimeInterval(-3600)
        )

        return AppViewModel(
            profile: AppProfile(
                fullName: "Вика Кашуркина",
                handle: "@vickdogmom",
                bio: "Собираю dog-friendly приложение: профиль, карта прогулок, мессенджер и трекер прививок.",
                avatarURL: URL(string: "https://images.unsplash.com/photo-1494790108377-be9c29b29330"),
                followers: 842,
                following: 301
            ),
            dogs: [dog1, dog2],
            posts: [ownPost],
            feed: [
                ownPost,
                AppPost(
                    id: UUID(),
                    authorName: "Анна Белова",
                    authorHandle: "@anna_walks",
                    authorPhotoURL: URL(string: "https://images.unsplash.com/photo-1438761681033-6461ffad8d80"),
                    content: "На карте хочется видеть не просто людей, а связку хозяин + собака + характер.",
                    imageURL: URL(string: "https://images.unsplash.com/photo-1507146426996-ef05306b995a"),
                    likes: 67,
                    comments: 9,
                    createdAt: .now.addingTimeInterval(-7200)
                )
            ],
            chats: [
                AppChat(id: UUID(), title: "Анна", lastMessage: "В парке в 18:30?", unreadCount: 2),
                AppChat(id: UUID(), title: "Dog walkers", lastMessage: "Рокси сегодня супер", unreadCount: 0)
            ],
            walkers: [
                AppWalker(id: UUID(), ownerName: "Анна", dogName: "Бимка", mood: .friendly, latitude: 55.751244, longitude: 37.618423),
                AppWalker(id: UUID(), ownerName: "Макс", dogName: "Рекс", mood: .aggressive, latitude: 55.753500, longitude: 37.610500),
                AppWalker(id: UUID(), ownerName: "Софи", dogName: "Булка", mood: .cautious, latitude: 55.748500, longitude: 37.624000)
            ],
            vaccinations: [
                AppVaccination(
                    id: UUID(),
                    dogID: dog1.id,
                    name: "Бешенство",
                    vaccinationDate: .now.addingTimeInterval(-60 * 60 * 24 * 25),
                    reminderAfterDays: 330,
                    nextDate: .now.addingTimeInterval(2_000_000),
                    notes: "Следующая ревакцинация через 11 месяцев."
                ),
                AppVaccination(
                    id: UUID(),
                    dogID: dog2.id,
                    name: "Щенячья первая",
                    vaccinationDate: .now.addingTimeInterval(-60 * 60 * 24 * 10),
                    reminderAfterDays: 20,
                    nextDate: .now.addingTimeInterval(800_000),
                    notes: "Нужно поставить напоминание."
                )
            ],
            currentUserId: UUID(),
            chatUseCase: MockChatUseCase(),
            profileService: MockUserProfileServiceAPI(),
            uploadService: MockUploadServiceAPI(),
            canEditProfile: true
        )
    }()
}

struct MockUserProfileServiceAPI: UserProfileServiceAPIProtocol {
    func create(userID: UUID, model: UserProfileModel) async throws -> UserProfileModel {
        model
    }

    func getByUserID(userID: UUID) async throws -> UserProfileModel {
        UserProfileModel(
            userId: userID,
            name: "Анна",
            surname: "Белова",
            username: "anna_walks",
            profilePhotoUrl: nil,
            bio: "",
            address: nil,
            birthDate: nil,
            lat: nil,
            lon: nil,
            locationStatus: 0,
            locationVisibility: 0,
            showLocation: false,
            isProfilePublic: true
        )
    }

    func search(query: String, limit: Int) async throws -> [UserProfileModel] {
        [
            UserProfileModel(
                userId: UUID(),
                name: "Анна",
                surname: "Белова",
                username: "anna_walks",
                profilePhotoUrl: nil,
                bio: "Гуляем в парке каждый вечер.",
                address: nil,
                birthDate: nil,
                lat: nil,
                lon: nil,
                locationStatus: 0,
                locationVisibility: 0,
                showLocation: false,
                isProfilePublic: true
            )
        ]
    }

    func update(userID: UUID, input: UpdateProfileInput) async throws {}
    func delete(userID: UUID) async throws {}
}

struct MockUploadServiceAPI: UploadServiceAPIProtocol {
    func uploadAvatar(_ imageData: Data, mimeType: String?) async throws -> MediaInfoModel {
        MediaInfoModel(url: "", mimeType: mimeType ?? "image/jpeg")
    }

    func uploadPostImage(_ imageData: Data, mimeType: String?) async throws -> MediaInfoModel {
        MediaInfoModel(url: "", mimeType: mimeType ?? "image/jpeg")
    }

    func uploadDogImage(_ imageData: Data, mimeType: String?) async throws -> MediaInfoModel {
        MediaInfoModel(url: "", mimeType: mimeType ?? "image/jpeg")
    }
}

struct MockChatUseCase: ChatUseCase {
    func createPrivateChat(user1: UUID, user2: UUID) async throws -> Chat {
        Chat(id: UUID(), isGroup: false, title: "Чат", photoUrl: "", createdAt: .now)
    }

    func createGroupChat(title: String, creator: UUID, members: [UUID]) async throws -> Chat {
        Chat(id: UUID(), isGroup: true, title: title, photoUrl: "", createdAt: .now)
    }

    func getChatByID(id: UUID) async throws -> Chat {
        Chat(id: id, isGroup: false, title: "Чат", photoUrl: "", createdAt: .now)
    }

    func getUserChats(userID: UUID) async throws -> [Chat] { [] }

    func updateChatTitle(chatID: UUID, title: String) async throws {}
    func updateChatPhoto(chatID: UUID, photoUrl: String) async throws {}
    func leaveChat(userID: UUID, chatID: UUID) async throws {}
    func deleteChat(chatID: UUID) async throws {}
    func getChatMembers(chatID: UUID) async throws -> [ChatMember] { [] }
    func addMember(userID: UUID, chatID: UUID) async throws {}
    func removeMember(userID: UUID, chatID: UUID) async throws {}

    func getMessages(chatID: UUID, limit: Int, before: UUID?) async throws -> [Message] { [] }

    func sendMessage(
        chatID: UUID,
        text: String?,
        attachments: [AttachmentInput]?,
        replyToId: UUID?
    ) async throws -> Message {
        Message(
            id: UUID(),
            chatId: chatID,
            senderId: UUID(),
            text: text,
            replyToId: replyToId,
            createdAt: .now,
            editedAt: nil,
            attachments: [],
            reactions: []
        )
    }

    func markAsRead(chatID: UUID, userID: UUID) async throws {}
    func sendTyping(chatID: UUID) async throws {}
}
