import SwiftUI

struct FeedView: View {
    @EnvironmentObject private var appViewModel: AppViewModel
    @State private var searchText = ""

    private var filteredFeed: [AppPost] {
        let query = searchText.trimmingCharacters(in: .whitespacesAndNewlines).lowercased()
        guard !query.isEmpty else { return appViewModel.feed }

        return appViewModel.feed.filter {
            $0.content.lowercased().contains(query) ||
            $0.authorName.lowercased().contains(query) ||
            $0.authorHandle.lowercased().contains(query)
        }
    }

    private var filteredAccounts: [AppProfile] {
        let query = searchText.trimmingCharacters(in: .whitespacesAndNewlines).lowercased()
        guard !query.isEmpty else { return [] }

        return appViewModel.accounts.filter {
            $0.fullName.lowercased().contains(query) ||
            $0.handle.lowercased().contains(query) ||
            $0.bio.lowercased().contains(query)
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
                                    appViewModel.toggleLike(postID: post.id)
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
            if !filteredAccounts.isEmpty {
                VStack(alignment: .leading, spacing: 10) {
                    Text("Аккаунты")
                        .font(.headline)
                        .foregroundStyle(.white)

                    ForEach(filteredAccounts, id: \.handle) { account in
                        HStack(spacing: 12) {
                            AsyncImage(url: account.avatarURL) { phase in
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
                                Text(account.fullName)
                                    .font(.subheadline.bold())
                                    .foregroundStyle(.white)
                                Text(account.handle)
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
                }
            }

            if !filteredFeed.isEmpty {
                VStack(alignment: .leading, spacing: 10) {
                    Text("Посты")
                        .font(.headline)
                        .foregroundStyle(.white)

                    ForEach(filteredFeed) { post in
                        PostView(post: post) {
                            appViewModel.toggleLike(postID: post.id)
                        }
                    }
                }
            }

            if filteredAccounts.isEmpty && filteredFeed.isEmpty {
                ContentUnavailableView(
                    "Ничего не найдено",
                    systemImage: "magnifyingglass",
                    description: Text("Попробуй другой запрос")
                )
                .foregroundStyle(.white)
            }
        }
    }
}
