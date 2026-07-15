import SwiftUI

struct SearchView: View {
    @EnvironmentObject private var appViewModel: AppViewModel
    @State private var searchText = ""
    @State private var searchResults: [UserProfileModel] = []
    @State private var isSearching = false
    @State private var searchError: String?
    @State private var searchTask: Task<Void, Never>?

    var body: some View {
        NavigationStack {
            List {
                if isSearching {
                    HStack {
                        Spacer()
                        ProgressView()
                        Spacer()
                    }
                } else if let searchError {
                    Text(searchError)
                        .foregroundStyle(.red)
                } else if searchText.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty {
                    ContentUnavailableView(
                        "Найди аккаунт",
                        systemImage: "magnifyingglass",
                        description: Text("Начни вводить имя или никнейм")
                    )
                } else if searchResults.isEmpty {
                    ContentUnavailableView(
                        "Ничего не найдено",
                        systemImage: "person.crop.circle.badge.questionmark",
                        description: Text("Попробуй другой никнейм")
                    )
                } else {
                    ForEach(searchResults, id: \.userId) { account in
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
                    }
                }
            }
            .listStyle(.plain)
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

    private func scheduleSearch(_ value: String) {
        searchTask?.cancel()
        searchError = nil

        let trimmed = value.trimmingCharacters(in: .whitespacesAndNewlines)
        guard trimmed.count >= 2 else {
            searchResults = []
            isSearching = false
            return
        }

        isSearching = true
        searchTask = Task {
            try? await Task.sleep(nanoseconds: 300_000_000)
            guard !Task.isCancelled else { return }

            do {
                let profiles = try await appViewModel.searchProfiles(query: trimmed)
                guard !Task.isCancelled else { return }
                await MainActor.run {
                    searchResults = profiles
                    isSearching = false
                }
            } catch {
                guard !Task.isCancelled else { return }
                await MainActor.run {
                    searchResults = []
                    searchError = error.localizedDescription
                    isSearching = false
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
