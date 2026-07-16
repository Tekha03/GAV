import SwiftUI

struct SearchView: View {
    @EnvironmentObject private var appViewModel: AppViewModel

    @State private var searchText = ""
    @State private var searchResults: [UserProfileModel] = []
    @State private var state: AppScreenState = .content
    @State private var searchTask: Task<Void, Never>?

    var body: some View {
        NavigationStack {
            Group {
                switch state {
                case .loading, .error, .offline:
                    AppStatusView(
                        state: state,
                        retryAction: retrySearch
                    )

                case .content:
                    content
                }
            }
            .navigationTitle("Поиск")
            .searchable(text: $searchText, prompt: "Найти аккаунт")
            .onChange(of: searchText) { _, newValue in
                scheduleSearch(newValue)
            }
            .onDisappear {
                searchTask?.cancel()
            }
        }
    }

    @ViewBuilder
    private var content: some View {
        let trimmedSearchText = searchText.trimmingCharacters(
            in: .whitespacesAndNewlines
        )

        if trimmedSearchText.isEmpty {
            ContentUnavailableView(
                "Найди аккаунт",
                systemImage: "magnifyingglass",
                description: Text("Начни вводить имя или никнейм")
            )
        } else if trimmedSearchText.count < 2 {
            ContentUnavailableView(
                "Продолжай ввод",
                systemImage: "text.cursor",
                description: Text("Введите минимум 2 символа")
            )
        } else if searchResults.isEmpty {
            ContentUnavailableView(
                "Ничего не найдено",
                systemImage: "person.crop.circle.badge.questionmark",
                description: Text("Попробуй другой никнейм")
            )
        } else {
            List(searchResults, id: \.userId) { account in
                NavigationLink {
                    UserProfileDetailView(profile: account)
                } label: {
                    accountRow(for: account)
                }
            }
            .listStyle(.plain)
        }
    }

    private func accountRow(_ account: UserProfileModel) -> some View {
        HStack(spacing: 12) {
            AsyncImage(
                url: account.profilePhotoURL.flatMap {
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
                        .fill(.secondary.opacity(0.12))
                        .overlay {
                            Image(systemName: "person.fill")
                                .foregroundStyle(.secondary)
                        }
                }
            }
            .frame(width: 48, height: 48)
            .clipShape(Circle())

            VStack(alignment: .leading, spacing: 3) {
                Text(displayName(for: account))
                    .font(.headline)

                Text("@\(account.username)")
                    .font(.subheadline)
                    .foregroundStyle(.secondary)

                Text(account.bio)
                    .font(.caption)
                    .foregroundStyle(.secondary)
                    .lineLimit(2)
            }
        }
        .padding(.vertical, 4)
    }

    private func scheduleSearch(_ value: String) {
        searchTask?.cancel()

        let trimmed = value.trimmingCharacters(
            in: .whitespacesAndNewlines
        )

        guard trimmed.count >= 2 else {
            searchResults = []
            state = .content
            return
        }

        state  = .loading(message: "Ищем аккаунты...")

        searchTask = Task {
            do {
                try await Task.sleep(nanoseconds: 300_000_000)

                guard !Task.isCancelled else {
                    return
                }

                let profiles = try await appViewModel.searchProfiles(
                    query: trimmed
                )

                guard !Task.isCancelled else {
                    return
                }

                searchResults = profiles
                state = .content
            } catch is CancellationError {
                return
            } catch {
                guard !Task.isCancelled else {
                    return
                }

                searchResults = []
                state = .from(error)
            }
        }
    }

    private func retrySearch() {
        scheduleSearch(searchText)
    }

    private func displayName(for profile: UserProfileModel) -> String {
        let fullName = [profile.name, profile.surname]
            .map {
                $0.trimmingCharacters(in: .whitespacesAndNewlines)
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
