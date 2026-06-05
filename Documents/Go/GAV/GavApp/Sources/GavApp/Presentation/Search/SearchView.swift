import SwiftUI

struct SearchView: View {
    @EnvironmentObject private var appViewModel: AppViewModel
    @State private var searchText = ""

    private var filteredAccounts: [AppProfile] {
        let query = searchText.trimmingCharacters(in: .whitespacesAndNewlines).lowercased()
        guard !query.isEmpty else { return appViewModel.accounts }

        return appViewModel.accounts.filter {
            $0.fullName.lowercased().contains(query) ||
            $0.handle.lowercased().contains(query) ||
            $0.bio.lowercased().contains(query)
        }
    }

    var body: some View {
        NavigationStack {
            List {
                ForEach(filteredAccounts.indices, id: \.self) { index in
                    let account = filteredAccounts[index]

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
                        .frame(width: 48, height: 48)
                        .clipShape(Circle())

                        VStack(alignment: .leading, spacing: 3) {
                            Text(account.fullName)
                                .font(.headline)
                            Text(account.handle)
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
            .listStyle(.plain)
            .navigationTitle("Поиск")
            .searchable(text: $searchText, prompt: "Найти аккаунт")
        }
    }
}