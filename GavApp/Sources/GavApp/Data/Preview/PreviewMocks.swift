
import Foundation


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
            dogService: MockDogServiceAPI(),
            postService: MockPostServiceAPI(),
            feedService: MockFeedServiceAPI(),
            userService: MockUserServiceAPI(),
            followService: MockFollowServiceAPI(),
            statsService: MockStatsServiceAPI(),
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

struct MockDogServiceAPI: DogServiceAPIProtocol {
    func create(ownerID: UUID, input: CreateDogInput) async throws -> DogModel {
        DogModel(
            id: UUID(),
            ownerId: ownerID,
            name: input.name,
            breed: input.breed,
            photoUrl: input.photoUrl,
            status: input.status,
            age: input.age,
            gender: input.gender
        )
    }

    func getPrivate(dogID: UUID) async throws -> DogModel {
        try await getPublic(dogID: dogID)
    }

    func getPublic(dogID: UUID) async throws -> DogModel {
        DogModel(
            id: dogID,
            ownerId: UUID(),
            name: "Бимка",
            breed: "Лабрадор",
            photoUrl: "",
            status: "friendly",
            age: "adult",
            gender: "female"
        )
    }

    func update(dogID: UUID, input: UpdateDogInput) async throws {}
    func delete(dogID: UUID) async throws {}
    func listByOwnerID(ownerID: UUID) async throws -> [DogModel] { [] }
}

struct MockPostServiceAPI: PostServiceAPIProtocol {
    func create(userID: UUID, content: String, imageUrl: String?) async throws -> PostModel {
        PostModel(id: UUID(), userId: userID, content: content, imageUrl: imageUrl, createdAt: .now)
    }

    func getByID(id: UUID) async throws -> PostModel {
        PostModel(id: id, userId: UUID(), content: "", imageUrl: nil, createdAt: .now)
    }

    func listByUser(userID: UUID) async throws -> [PostModel] { [] }
    func delete(id: UUID) async throws {}
    func addLike(postID: UUID) async throws {}
    func removeLike(postID: UUID) async throws {}
    func createComment(postID: UUID, userID: UUID, content: String) async throws -> CommentModel {
        CommentModel(id: UUID(), postId: postID, userId: userID, content: content, createdAt: .now)
    }
    func listCommentsByPostID(postID: UUID) async throws -> [CommentModel] { [] }
    func deleteComment(id: UUID) async throws {}
}

struct MockFeedServiceAPI: FeedServiceAPIProtocol {
    func getFeed(userID: UUID, before: Date?, limit: Int) async throws -> [PostModel] { [] }
}

struct MockUserServiceAPI: UserServiceAPIProtocol {
    func getByID(id: UUID) async throws -> UserModel {
        UserModel(id: id, email: "preview@gav.app", role: "user", createdAt: .now, updatedAt: .now)
    }

    func update(id: UUID, input: UpdateUserInput) async throws {}
    func delete(id: UUID) async throws {}

    func getByEmail(email: String) async throws -> UserModel {
        UserModel(id: UUID(), email: email, role: "user", createdAt: .now, updatedAt: .now)
    }

    func updateLocation(id: UUID, input: UpdateLocationInput) async throws {}

    func findDogsNearby(
        id: UUID,
        centerLat: Double,
        centerLon: Double,
        radiusMeters: Double
    ) async throws -> [DogModel] { [] }
}

struct MockFollowServiceAPI: FollowServiceAPIProtocol {
    func follow(userID: UUID) async throws {}
    func unfollow(userID: UUID) async throws {}
    func getFollowers(userID: UUID) async throws -> [FollowModel] { [] }
    func getFollowing(userID: UUID) async throws -> [FollowModel] { [] }
}

struct MockStatsServiceAPI: StatsServiceAPIProtocol {
    func userStats(userID: UUID) async throws -> UserStatsModel {
        UserStatsModel(
            id: UUID(),
            userId: userID,
            postCount: 0,
            followersCount: 0,
            followingsCount: 0,
            dogsCount: 0,
            createdAt: .now,
            updatedAt: .now
        )
    }

    func profileStats(userID: UUID) async throws -> ProfileStatsModel {
        ProfileStatsModel(userId: userID, postCount: 0, followersCount: 0, followingsCount: 0)
    }

    func postStats(postID: UUID) async throws -> PostStatsModel {
        PostStatsModel(
            id: UUID(),
            postId: postID,
            likesCount: 0,
            commentsCount: 0,
            createdAt: .now,
            updatedAt: .now
        )
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
