import SwiftUI

struct UserProfileDetailView: View {
    @EnvironmentObject private var appViewModel: AppViewModel
    let profile: UserProfileModel

    @State private var isFollowing = false
    @State private var isBusy = false
    @State private var screenState: AppScreenState = .loading(
        message: "Загружаем профиль..."
    )
    @State private var actionErrorMessage: String?
    @State private var dogs: [AppDog] = []
    @State private var posts: [AppPost] = []
    @State private var stats: ProfileStatsModel?
    @State private var selectedCommentsPost: AppPost?

    var body: some View {
        ZStack {
            profileBackground

            switch screenState {
            case .loading, .error, .offline:
                AppStatusView(
                    state: screenState,
                    retryAction: {
                        Task {
                            await loadInitialState()
                        }
                    }
                )
                .foregroundStyle(.white)

            case .content:
                profileContent
            }
        }
        .navigationTitle(
            profile.username.isEmpty
                ? "Профиль"
                : "@\(profile.username)"
        )
        .navigationBarTitleDisplayMode(.inline)
        .preferredColorScheme(.dark)
        .task {
            await loadInitialState()
        }
        .sheet(item: $selectedCommentsPost) { post in
            CommentsView(post: post) { count in
                if let index = posts.firstIndex(
                    where: {
                        $0.id == post.id
                    }
                ) {
                    posts[index].comments = count
                }
            }
            .environmentObject(appViewModel)
        }
    }

    private var profileBackground: some View {
        LinearGradient(
            colors: [
                Color(red: 0.42, green: 0.22, blue: 0.72),
                .black
            ],
            startPoint: .top,
            endPoint: .bottom
        )
        .ignoresSafeArea()
    }

    private var profileContent: some View {
        ScrollView {
            VStack(spacing: 18) {
                header
                dogsSection
                postsSection

                if let actionErrorMessage {
                    HStack(spacing: 10) {
                        Image(
                            systemName: "exclamationmark.circle.fill"
                        )
                        .foregroundStyle(.orange)

                        Text(actionErrorMessage)
                            .font(.footnote)
                            .foregroundStyle(.white)
                            .frame(
                                maxWidth: .infinity,
                                alignment: .leading
                            )

                        Button {
                            self.actionErrorMessage = nil
                        } label: {
                            Image(systemName: "xmark")
                                .foregroundStyle(
                                    .white.opacity(0.7)
                                )
                        }
                    }
                    .padding(12)
                    .background(
                        .black.opacity(0.35),
                        in: RoundedRectangle(cornerRadius: 14)
                    )
                }
            }
            .padding(.horizontal, 16)
            .padding(.vertical, 20)
        }
        .refreshable {
            await loadInitialState(showLoading: false)
        }
    }

    private var header: some View {
        VStack(alignment: .leading, spacing: 16) {
            HStack(alignment: .center, spacing: 16) {
                avatar
                    .frame(width: 88, height: 88)

                VStack(alignment: .leading, spacing: 5) {
                    Text(displayName)
                        .font(.title3.bold())
                        .foregroundStyle(.white)

                    Text("@\(profile.username)")
                        .font(.subheadline)
                        .foregroundStyle(.white.opacity(0.74))
                }

                Spacer()
            }

            if !profile.bio.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty {
                Text(profile.bio)
                    .font(.footnote)
                    .foregroundStyle(.white.opacity(0.9))
                    .fixedSize(horizontal: false, vertical: true)
            }

            HStack(spacing: 10) {
                statCard(title: "Подписчики", value: Int(stats?.followersCount ?? 0))
                statCard(title: "Подписки", value: Int(stats?.followingsCount ?? 0))
                statCard(title: "Посты", value: Int(stats?.postCount ?? UInt(posts.count)))
            }

            HStack(spacing: 10) {
                Button {
                    Task { await toggleFollow() }
                } label: {
                    Label(isFollowing ? "Отписаться" : "Подписаться", systemImage: isFollowing ? "person.fill.xmark" : "person.fill.badge.plus")
                        .frame(maxWidth: .infinity)
                }
                .buttonStyle(.borderedProminent)
                .tint(isFollowing ? .white.opacity(0.18) : .orange)
                .disabled(isBusy)

                Button {
                    Task { await createChat() }
                } label: {
                    Label("Написать", systemImage: "bubble.left.and.bubble.right.fill")
                        .frame(maxWidth: .infinity)
                }
                .buttonStyle(.bordered)
                .tint(.white)
                .disabled(isBusy)
            }
        }
        .padding(18)
        .background(
            RoundedRectangle(cornerRadius: 24, style: .continuous)
                .fill(.white.opacity(0.08))
        )
        .overlay(
            RoundedRectangle(cornerRadius: 24, style: .continuous)
                .stroke(.white.opacity(0.08), lineWidth: 1)
        )
    }

    private func statCard(title: String, value: Int) -> some View {
        VStack(spacing: 3) {
            Text("\(value)")
                .font(.headline.bold())
                .foregroundStyle(.white)
            Text(title)
                .font(.caption2)
                .foregroundStyle(.white.opacity(0.68))
                .lineLimit(1)
                .minimumScaleFactor(0.78)
        }
        .frame(maxWidth: .infinity)
        .padding(.vertical, 10)
        .background(.white.opacity(0.08), in: RoundedRectangle(cornerRadius: 14, style: .continuous))
    }

    private var avatar: some View {
        AsyncImage(url: profile.profilePhotoUrl.flatMap { MediaURLResolver.resolve($0) }) { phase in
            switch phase {
            case .success(let image):
                image.resizable().scaledToFill()
            default:
                ZStack {
                    Circle().fill(.white.opacity(0.12))
                    Image(systemName: "person.fill")
                        .font(.title)
                        .foregroundStyle(.white.opacity(0.9))
                }
            }
        }
        .clipShape(Circle())
        .overlay(Circle().stroke(.white.opacity(0.22), lineWidth: 2))
    }

    private var dogsSection: some View {
        VStack(alignment: .leading, spacing: 12) {
            Text("Собаки")
                .font(.headline)
                .foregroundStyle(.white)

            if dogs.isEmpty {
                emptyState("Собак пока нет", systemImage: "pawprint")
            } else {
                ScrollView(.horizontal, showsIndicators: false) {
                    HStack(spacing: 12) {
                        ForEach(dogs) { dog in
                            dogCard(dog)
                        }
                    }
                    .padding(.vertical, 2)
                }
            }
        }
        .frame(maxWidth: .infinity, alignment: .leading)
    }

    private var postsSection: some View {
        VStack(alignment: .leading, spacing: 12) {
            Text("Посты")
                .font(.headline)
                .foregroundStyle(.white)

            if posts.isEmpty {
                emptyState("Постов пока нет", systemImage: "photo.on.rectangle")
            } else {
                ForEach(posts) { post in
                    PostView(post: post) {
                        Task { await toggleLike(postID: post.id) }
                    } onComment: {
                        selectedCommentsPost = post
                    }
                }
            }
        }
        .frame(maxWidth: .infinity, alignment: .leading)
    }

    private func dogCard(_ dog: AppDog) -> some View {
        VStack(alignment: .leading, spacing: 8) {
            AsyncImage(url: dog.photoURL) { phase in
                switch phase {
                case .success(let image):
                    image.resizable().scaledToFill()
                default:
                    ZStack {
                        RoundedRectangle(cornerRadius: 14, style: .continuous)
                            .fill(.white.opacity(0.1))
                        Image(systemName: "pawprint.fill")
                            .font(.title2)
                            .foregroundStyle(.white.opacity(0.85))
                    }
                }
            }
            .frame(width: 132, height: 108)
            .clipShape(RoundedRectangle(cornerRadius: 14, style: .continuous))

            Text(dog.name)
                .font(.subheadline.bold())
                .foregroundStyle(.white)
                .lineLimit(1)

            Text(dog.breed)
                .font(.caption)
                .foregroundStyle(.white.opacity(0.7))
                .lineLimit(1)

            if !dog.notes.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty {
                Text(dog.notes)
                    .font(.caption2)
                    .foregroundStyle(.white.opacity(0.62))
                    .lineLimit(2)
            }
        }
        .frame(width: 132, alignment: .leading)
        .padding(10)
        .background(.white.opacity(0.08), in: RoundedRectangle(cornerRadius: 18, style: .continuous))
    }

    private func emptyState(_ title: String, systemImage: String) -> some View {
        HStack(spacing: 10) {
            Image(systemName: systemImage)
                .font(.headline)
                .foregroundStyle(.orange)
            Text(title)
                .font(.subheadline)
                .foregroundStyle(.white.opacity(0.74))
            Spacer()
        }
        .padding(14)
        .background(.white.opacity(0.08), in: RoundedRectangle(cornerRadius: 18, style: .continuous))
    }

    private var displayName: String {
        let fullName = [profile.name, profile.surname]
            .map { $0.trimmingCharacters(in: .whitespacesAndNewlines) }
            .filter { !$0.isEmpty }
            .joined(separator: " ")
        return fullName.isEmpty ? "@\(profile.username)" : fullName
    }

    private func loadInitialState(
        showLoading: Bool = true
    ) async {
        if showLoading && dogs.isEmpty && posts.isEmpty {
            screenState = .loading(
                message: "Загружаем профиль..."
            )
        }

        actionErrorMessage = nil

        do {
            async let followingTask = appViewModel.isFollowing(
                profile.userId
            )

            async let contentTask = appViewModel.loadProfileContent(
                for: profile
            )

            async let statsTask = appViewModel.loadProfileStats(
                for: profile.userId
            )

            let following = await followingTask
            let content = try await contentTask

            isFollowing = following
            dogs = content.dogs
            posts = content.posts

            do {
                stats = try await statsTask
            } catch {
                stats = nil
            }

            screenState = .content
        } catch {
            if dogs.isEmpty && posts.isEmpty {
                screenState = .from(error)
            } else {
                screenState = .content
                actionErrorMessage = error.localizedDescription
            }
        }
    }

    private func toggleFollow() async {
        isBusy = true
        actionErrorMessage = nil
        defer { isBusy = false }

        do {
            if isFollowing {
                adjustFollowerCount(by: -1)
                try await appViewModel.unfollow(profile.userId)
                isFollowing = false
            } else {
                adjustFollowerCount(by: 1)
                try await appViewModel.follow(profile.userId)
                isFollowing = true
            }
            if let freshStats = try? await appViewModel.loadProfileStats(for: profile.userId) {
                stats = freshStats
            }
        } catch {
            adjustFollowerCount(by: isFollowing ? 1 : -1)
            actionErrorMessage = "Не удалось обновить подписку"
        }
    }

    private func adjustFollowerCount(by delta: Int) {
        let current = stats ?? ProfileStatsModel(
            userId: profile.userId,
            postCount: UInt(posts.count),
            followersCount: 0,
            followingsCount: 0
        )
        let nextFollowers = max(0, Int(current.followersCount) + delta)

        stats = ProfileStatsModel(
            userId: current.userId,
            postCount: current.postCount,
            followersCount: UInt(nextFollowers),
            followingsCount: current.followingsCount
        )
    }

    private func toggleLike(postID: UUID) async {
        let wasLiked = posts.first(where: { $0.id == postID })?.isLiked == true ||
            appViewModel.likedPostIDs.contains(postID)

        await appViewModel.toggleLike(postID: postID)

        let isLiked = appViewModel.likedPostIDs.contains(postID)
        guard let index = posts.firstIndex(where: { $0.id == postID }) else { return }

        if let stats = await appViewModel.refreshPostStats(postID: postID) {
            posts[index].likes = Int(stats.likesCount)
            posts[index].comments = Int(stats.commentsCount)
        } else if wasLiked != isLiked {
            posts[index].likes = max(0, posts[index].likes + (isLiked ? 1 : -1))
        }
        posts[index].isLiked = isLiked
    }

    private func createChat() async {
        isBusy = true
        actionErrorMessage = nil
        defer { isBusy = false }

        do {
            try await appViewModel.createPrivateChat(with: profile.userId)
        } catch {
            actionErrorMessage = "Не удалось создать чат"
        }
    }
}
