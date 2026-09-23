import SwiftUI

struct FeedView: View {
    @EnvironmentObject private var appViewModel: AppViewModel

    @State private var searchText = ""
    @State private var selectedCommentsPost: AppPost?
    @State private var accountSearchResults: [UserProfileModel] = []
    @State private var screenState: AppScreenState = .loading(
        message: "Загружаем ленту..."
    )
    @State private var searchState: AppScreenState = .content
    @State private var accountSearchTask: Task<Void, Never>?

    private var filteredFeed: [AppPost] {
        let query = searchText
            .trimmingCharacters(in: .whitespacesAndNewlines)
            .lowercased()

        guard !query.isEmpty else {
            return appViewModel.feed
        }

        return appViewModel.feed.filter {
            $0.content.lowercased().contains(query) ||
            $0.authorName.lowercased().contains(query) ||
            $0.authorHandle.lowercased().contains(query)
        }
    }

    var body: some View {
        NavigationStack {
            ZStack {
                feedBackground

                switch screenState {
                case .loading, .error, .offline:
                    AppStatusView(
                        state: screenState,
                        retryAction: {
                            Task {
                                await loadFeed()
                            }
                        }
                    )
                    .foregroundStyle(.white)

                case .content:
                    feedContent
                }
            }
            .navigationTitle("Лента")
            .searchable(
                text: $searchText,
                prompt: "Поиск постов и аккаунтов"
            )
            .onChange(of: searchText) { _, newValue in
                scheduleAccountSearch(newValue)
            }
            .onDisappear {
                accountSearchTask?.cancel()
            }
            .task {
                await loadFeed()
            }
            .sheet(item: $selectedCommentsPost) { post in
                CommentsView(post: post)
                    .environmentObject(appViewModel)
            }
        }
    }

    private var feedContent: some View {
        ScrollView {
            VStack(spacing: 16) {
                if !searchText
                    .trimmingCharacters(in: .whitespacesAndNewlines)
                    .isEmpty {
                    searchResultsSection
                } else if appViewModel.feed.isEmpty {
                    emptyFeed
                } else {
                    ForEach(appViewModel.feed) { post in
                        postView(post)
                    }
                }
            }
            .frame(maxWidth: 620)
            .frame(maxWidth: .infinity)
            .padding(.horizontal, 16)
            .padding(.vertical, 20)
        }
        .refreshable {
            await loadFeed(showLoading: false)
        }
    }

    private var feedBackground: some View {
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

    private var emptyFeed: some View {
        VStack(spacing: 14) {
            Image(systemName: "photo.on.rectangle.angled")
                .font(.system(size: 42, weight: .semibold))
                .foregroundStyle(.orange)

            Text("Лента пока пустая")
                .font(.headline)
                .foregroundStyle(.white)

            Text("Когда появятся посты, они будут здесь.")
                .font(.subheadline)
                .multilineTextAlignment(.center)
                .foregroundStyle(.white.opacity(0.68))
                .frame(maxWidth: 260)
        }
        .frame(maxWidth: .infinity)
        .frame(minHeight: 420)
    }

    @ViewBuilder
    private var searchResultsSection: some View {
        VStack(alignment: .leading, spacing: 18) {
            switch searchState {
            case .loading:
                HStack {
                    Spacer()

                    ProgressView()
                        .tint(.white)

                    Spacer()
                }
                .padding(.vertical, 8)

            case .error, .offline:
                AppStatusView(
                    state: searchState,
                    retryAction: retrySearch
                )
                .frame(minHeight: 220)
                .foregroundStyle(.white)

            case .content:
                searchContent
            }
        }
    }

    @ViewBuilder
    private var searchContent: some View {
        if !accountSearchResults.isEmpty {
            VStack(alignment: .leading, spacing: 10) {
                Text("Аккаунты")
                    .font(.headline)
                    .foregroundStyle(.white)

                ForEach(
                    accountSearchResults,
                    id: \.userId
                ) { account in
                    accountRow(account)
                }
            }
        }

        if !filteredFeed.isEmpty {
            VStack(alignment: .leading, spacing: 10) {
                Text("Посты")
                    .font(.headline)
                    .foregroundStyle(.white)

                ForEach(filteredFeed) { post in
                    postView(post)
                }
            }
        }

        if accountSearchResults.isEmpty &&
            filteredFeed.isEmpty {
            ContentUnavailableView(
                "Ничего не найдено",
                systemImage: "magnifyingglass",
                description: Text("Попробуй другой запрос")
            )
            .foregroundStyle(.white)
        }
    }

    private func postView(_ post: AppPost) -> some View {
        PostView(post: post) {
            Task {
                await appViewModel.toggleLike(
                    postID: post.id
                )
            }
        } onComment: {
            selectedCommentsPost = post
        }
    }

    private func accountRow(
        _ account: UserProfileModel
    ) -> some View {
        NavigationLink {
            UserProfileDetailView(profile: account)
        } label: {
            HStack(spacing: 12) {
                AsyncImage(
                    url: account.profilePhotoUrl.flatMap {
                        MediaURLResolver.resolve($0)
                    }
                ) { phase in
                    switch phase {
                    case .success(let image):
                        image
                            .resizable()
                            .scaledToFill()

                    default:
                        Circle()
                            .fill(.white.opacity(0.12))
                            .overlay {
                                Image(systemName: "person.fill")
                            }
                    }
                }
                .frame(width: 44, height: 44)
                .clipShape(Circle())

                VStack(alignment: .leading, spacing: 2) {
                    Text(displayName(for: account))
                        .font(.subheadline.bold())
                        .foregroundStyle(.white)

                    Text("@\(account.username)")
                        .font(.caption)
                        .foregroundStyle(.white.opacity(0.75))

                    Text(account.bio)
                        .font(.caption2)
                        .foregroundStyle(.white.opacity(0.65))
                        .lineLimit(2)
                }

                Spacer()
            }
            .padding(12)
            .background(
                .white.opacity(0.08),
                in: RoundedRectangle(cornerRadius: 16)
            )
        }
        .buttonStyle(.plain)
    }

    private func loadFeed(
        showLoading: Bool = true
    ) async {
        if showLoading && appViewModel.feed.isEmpty {
            screenState = .loading(
                message: "Загружаем ленту..."
            )
        }

        do {
            try await appViewModel.reloadFeed()
            screenState = .content
        } catch {
            if appViewModel.feed.isEmpty {
                screenState = .from(error)
            } else {
                screenState = .content
            }
        }
    }

    private func scheduleAccountSearch(
        _ value: String
    ) {
        accountSearchTask?.cancel()

        let trimmed = value.trimmingCharacters(
            in: .whitespacesAndNewlines
        )

        guard trimmed.count >= 2 else {
            accountSearchResults = []
            searchState = .content
            return
        }

        searchState = .loading(
            message: "Ищем аккаунты..."
        )

        accountSearchTask = Task {
            do {
                try await Task.sleep(
                    nanoseconds: 300_000_000
                )

                guard !Task.isCancelled else {
                    return
                }

                let profiles = try await appViewModel.searchProfiles(
                    query: trimmed
                )

                guard !Task.isCancelled else {
                    return
                }

                accountSearchResults = profiles
                searchState = .content
            } catch is CancellationError {
                return
            } catch {
                guard !Task.isCancelled else {
                    return
                }

                accountSearchResults = []
                searchState = .from(error)
            }
        }
    }

    private func retrySearch() {
        scheduleAccountSearch(searchText)
    }

    private func displayName(
        for profile: UserProfileModel
    ) -> String {
        let fullName = [
            profile.name,
            profile.surname
        ]
        .map {
            $0.trimmingCharacters(
                in: .whitespacesAndNewlines
            )
        }
        .filter {
            !$0.isEmpty
        }
        .joined(separator: " ")

        return fullName.isEmpty
            ? "@\(profile.username)"
            : fullName
    }
}
