import SwiftUI

struct FeedView: View {
    @EnvironmentObject private var appViewModel: AppViewModel
    @State private var searchText = ""
    @State private var selectedCommentsPost: AppPost?
    @State private var accountSearchResults: [UserProfileModel] = []
    @State private var isSearchingAccounts = false
    @State private var accountSearchTask: Task<Void, Never>?

    private var filteredFeed: [AppPost] {
        let query = searchText.trimmingCharacters(in: .whitespacesAndNewlines).lowercased()
        guard !query.isEmpty else { return appViewModel.feed }

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

                ScrollView {
                    VStack(spacing: 16) {
                        if !searchText.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty {
                            searchResultsSection
                        } else if filteredFeed.isEmpty {
                            emptyFeed
                        } else {
                            ForEach(filteredFeed) { post in
                                PostView(post: post) {
                                    Task { await appViewModel.toggleLike(postID: post.id) }
                                } onComment: {
                                    selectedCommentsPost = post
                                }
                            }
                        }
                    }
                    .frame(maxWidth: 620)
                    .frame(maxWidth: .infinity)
                    .padding(.horizontal, 16)
                    .padding(.vertical, 20)
                }
            }
            .navigationTitle("Лента")
            .searchable(text: $searchText, prompt: "Поиск постов и аккаунтов")
            .onChange(of: searchText) { _, newValue in
                scheduleAccountSearch(newValue)
            }
            .onDisappear {
                accountSearchTask?.cancel()
            }
            .task {
                await appViewModel.loadAuthenticatedContent()
            }
            .sheet(item: $selectedCommentsPost) { post in
                CommentsView(post: post)
                    .environmentObject(appViewModel)
            }
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

    private var searchResultsSection: some View {
        VStack(alignment: .leading, spacing: 18) {
            if isSearchingAccounts {
                HStack {
                    Spacer()
                    ProgressView()
                        .tint(.white)
                    Spacer()
                }
                .padding(.vertical, 8)
            } else if !accountSearchResults.isEmpty {
                VStack(alignment: .leading, spacing: 10) {
                    Text("Аккаунты")
                        .font(.headline)
                        .foregroundStyle(.white)

                    ForEach(accountSearchResults, id: \.userId) { account in
                        NavigationLink {
                            UserProfileDetailView(profile: account)
                        } label: {
                            HStack(spacing: 12) {
                                AsyncImage(url: account.profilePhotoUrl.flatMap { MediaURLResolver.resolve($0) }) { phase in
                                    switch phase {
                                    case .success(let image):
                                        image.resizable().scaledToFill()
                                    default:
                                        Circle()
                                            .fill(.white.opacity(0.12))
                                            .overlay(Image(systemName: "person.fill"))
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
                            .background(.white.opacity(0.08), in: RoundedRectangle(cornerRadius: 16))
                        }
                        .buttonStyle(.plain)
                    }
                }
            }

            if !filteredFeed.isEmpty {
                VStack(alignment: .leading, spacing: 10) {
                    Text("Посты")
                        .font(.headline)
                        .foregroundStyle(.white)

                    ForEach(filteredFeed) { post in
                        PostView(post: post) {
                            Task { await appViewModel.toggleLike(postID: post.id) }
                        } onComment: {
                            selectedCommentsPost = post
                        }
                    }
                }
            }

            if !isSearchingAccounts && accountSearchResults.isEmpty && filteredFeed.isEmpty {
                ContentUnavailableView(
                    "Ничего не найдено",
                    systemImage: "magnifyingglass",
                    description: Text("Попробуй другой запрос")
                )
                .foregroundStyle(.white)
            }
        }
    }

    private func scheduleAccountSearch(_ value: String) {
        accountSearchTask?.cancel()

        let trimmed = value.trimmingCharacters(in: .whitespacesAndNewlines)
        guard trimmed.count >= 2 else {
            accountSearchResults = []
            isSearchingAccounts = false
            return
        }

        isSearchingAccounts = true
        accountSearchTask = Task {
            try? await Task.sleep(nanoseconds: 300_000_000)
            guard !Task.isCancelled else { return }

            do {
                let profiles = try await appViewModel.searchProfiles(query: trimmed)
                guard !Task.isCancelled else { return }
                await MainActor.run {
                    accountSearchResults = profiles
                    isSearchingAccounts = false
                }
            } catch {
                guard !Task.isCancelled else { return }
                await MainActor.run {
                    accountSearchResults = []
                    isSearchingAccounts = false
                }
            }
        }
    }

    private func displayName(for profile: UserProfileModel) -> String {
        let fullName = [profile.name, profile.surname]
            .map { $0.trimmingCharacters(in: .whitespacesAndNewlines) }
            .filter { !$0.isEmpty }
            .joined(separator: " ")
        return fullName.isEmpty ? "@\(profile.username)" : fullName
    }
}
