import SwiftUI

struct ChatListView: View {
    @EnvironmentObject private var appViewModel: AppViewModel

    @State private var showAddChat = false
    @State private var state: AppScreenState = .loading(
        message: "Загружаем чаты..."
    )

    var body: some View {
        let vm = appViewModel

        NavigationStack {
            ZStack {
                chatBackground

                switch state {
                case .loading, .error, .offline:
                    AppStatusView(
                        state: state,
                        retryAction: {
                            Task {
                                await loadChats()
                            }
                        }
                    )
                    .foregroundStyle(.white)

                case .content:
                    if appViewModel.chats.isEmpty {
                        emptyState
                    } else {
                        chatList(viewModel: vm)
                    }
                }
            }
            .navigationTitle("Мессенджер")
            .toolbar {
                ToolbarItem(placement: .topBarTrailing) {
                    Button {
                        showAddChat = true
                    } label: {
                        Image(systemName: "plus.message.fill")
                    }
                    .accessibilityLabel("Добавить чат")
                    .tint(.orange)
                }
            }
            .sheet(isPresented: $showAddChat) {
                AddChatView(viewModel: appViewModel)
            }
            .preferredColorScheme(.dark)
            .task {
                await loadChats()
            }
        }
    }

    private var chatBackground: some View {
        LinearGradient(
            colors: [
                Color(red: 0.10, green: 0.08, blue: 0.13),
                .black
            ],
            startPoint: .top,
            endPoint: .bottom
        )
        .ignoresSafeArea()
    }

    private func chatList(viewModel: AppViewModel) -> some View {
        ScrollView {
            LazyVStack(spacing: 12) {
                ForEach(viewModel.chats) { chat in
                    NavigationLink {
                        ChatDetailView(
                            chat: chat.domainChat,
                            currentUserId: viewModel.currentUserId,
                            useCase: viewModel.chatUseCase
                        )
                    } label: {
                        chatRow(chat)
                    }
                    .buttonStyle(.plain)
                }
            }
            .padding(.horizontal, 16)
            .padding(.vertical, 20)
        }
        .refreshable {
            await loadChats(showLoading: false)
        }
    }

    private var emptyState: some View {
        VStack(spacing: 14) {
            Image(systemName: "bubble.left.and.bubble.right")
                .font(.system(size: 42, weight: .semibold))
                .foregroundStyle(.orange)

            Text("Чатов пока нет")
                .font(.headline)
                .foregroundStyle(.white)

            Text("Когда появится диалог, он будет здесь.")
                .font(.subheadline)
                .multilineTextAlignment(.center)
                .foregroundStyle(.white.opacity(0.68))
                .frame(maxWidth: 260)

            Button {
                showAddChat = true
            } label: {
                Label("Добавить чат", systemImage: "plus.message.fill")
                    .font(.headline.weight(.semibold))
                    .padding(.horizontal, 18)
                    .frame(height: 46)
                    .background(
                        Color.orange,
                        in: RoundedRectangle(cornerRadius: 8)
                    )
                    .foregroundStyle(.black)
            }
            .padding(.top, 4)
        }
        .padding(.horizontal, 24)
        .frame(maxWidth: .infinity, maxHeight: .infinity)
    }

    private func chatRow(_ chat: AppChat) -> some View {
        HStack(spacing: 14) {
            Circle()
                .fill(Color.orange.opacity(0.15))
                .frame(width: 46, height: 46)
                .overlay {
                    Image(
                        systemName: "bubble.left.and.bubble.right.fill"
                    )
                    .foregroundStyle(.orange)
                }

            VStack(alignment: .leading, spacing: 4) {
                Text(chat.title)
                    .font(.headline)
                    .foregroundStyle(.white)

                Text(chat.lastMessage)
                    .font(.subheadline)
                    .foregroundStyle(.white.opacity(0.72))
                    .lineLimit(1)
            }

            Spacer()

            if chat.unreadCount > 0 {
                Text("\(chat.unreadCount)")
                    .font(.caption.bold())
                    .foregroundStyle(.white)
                    .padding(8)
                    .background(Color.orange, in: Circle())
            }
        }
        .padding(14)
        .background(
            .white.opacity(0.08),
            in: RoundedRectangle(cornerRadius: 18)
        )
    }

    private func loadChats(showLoading: Bool = true) async {
        if showLoading && appViewModel.chats.isEmpty {
            state = .loading(message: "Загружаем чаты...")
        }

        do {
            try await appViewModel.loadChats()
            state = .content
        } catch {
            if appViewModel.chats.isEmpty {
                state = .from(error)
            } else {
                state = .content
            }
        }
    }
}

private struct AddChatView: View {
    @ObservedObject var viewModel: AppViewModel
    @Environment(\.dismiss) private var dismiss

    @State private var mode: ChatCreationMode = .privateChat
    @State private var title = ""
    @State private var query = ""
    @State private var searchResults: [UserProfileModel] = []
    @State private var selectedProfiles: [UserProfileModel] = []
    @State private var errorMessage: String?
    @State private var isSearching = false
    @State private var isSaving = false
    @State private var searchTask: Task<Void, Never>?

    var body: some View {
        NavigationStack {
            Form {
                Section {
                    Picker("Тип", selection: $mode) {
                        ForEach(ChatCreationMode.allCases) { mode in
                            Text(mode.title).tag(mode)
                        }
                    }
                    .pickerStyle(.segmented)
                }

                Section {
                    if mode == .groupChat {
                        TextField("Название группы", text: $title)
                            .textInputAutocapitalization(.sentences)
                    }

                    TextField(mode.searchFieldTitle, text: $query)
                        .textInputAutocapitalization(.never)
                        .autocorrectionDisabled()
                        .onChange(of: query) { _, newValue in
                            scheduleSearch(newValue)
                        }
                } footer: {
                    Text(mode.footer)
                }

                if !selectedProfiles.isEmpty {
                    Section("Выбрано") {
                        ForEach(selectedProfiles, id: \.userId) { profile in
                            profileRow(profile, isSelected: true) {
                                selectedProfiles.removeAll { $0.userId == profile.userId }
                            }
                        }
                    }
                }

                Section("Профили") {
                    if isSearching {
                        HStack {
                            Spacer()
                            ProgressView()
                            Spacer()
                        }
                    } else if searchResults.isEmpty {
                        Text(query.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty ? "Начните вводить имя или ник." : "Ничего не найдено.")
                            .foregroundStyle(.secondary)
                    } else {
                        ForEach(searchResults, id: \.userId) { profile in
                            profileRow(profile, isSelected: isSelected(profile)) {
                                toggleSelection(profile)
                            }
                        }
                    }
                }

                if let errorMessage {
                    Section {
                        Text(errorMessage)
                            .font(.footnote)
                            .foregroundStyle(.red)
                    }
                }
            }
            .navigationTitle("Новый чат")
            .toolbar {
                ToolbarItem(placement: .cancellationAction) {
                    Button("Отмена") {
                        dismiss()
                    }
                }

                ToolbarItem(placement: .confirmationAction) {
                    Button {
                        Task { await createChat() }
                    } label: {
                        if isSaving {
                            ProgressView()
                        } else {
                            Text("Создать")
                        }
                    }
                    .disabled(!isFormValid || isSaving)
                }
            }
            .onDisappear {
                searchTask?.cancel()
            }
        }
    }

    private var isFormValid: Bool {
        switch mode {
        case .privateChat:
            return selectedProfiles.count == 1
        case .groupChat:
            return !title.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty && !selectedProfiles.isEmpty
        }
    }

    private func createChat() async {
        isSaving = true
        errorMessage = nil

        do {
            switch mode {
            case .privateChat:
                guard let profile = selectedProfiles.first else {
                    throw ChatCreationError.noProfileSelected
                }
                try await viewModel.createPrivateChat(with: profile.userId)
            case .groupChat:
                try await viewModel.createGroupChat(
                    title: title.trimmingCharacters(in: .whitespacesAndNewlines),
                    memberIDs: selectedProfiles.map(\.userId)
                )
            }
            dismiss()
        } catch {
            errorMessage = error.localizedDescription
        }

        isSaving = false
    }

    private func scheduleSearch(_ value: String) {
        searchTask?.cancel()
        errorMessage = nil

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
                let profiles = try await viewModel.searchProfiles(query: trimmed)
                guard !Task.isCancelled else { return }
                await MainActor.run {
                    searchResults = profiles
                    isSearching = false
                }
            } catch {
                guard !Task.isCancelled else { return }
                await MainActor.run {
                    searchResults = []
                    errorMessage = error.localizedDescription
                    isSearching = false
                }
            }
        }
    }

    private func isSelected(_ profile: UserProfileModel) -> Bool {
        selectedProfiles.contains { $0.userId == profile.userId }
    }

    private func toggleSelection(_ profile: UserProfileModel) {
        if isSelected(profile) {
            selectedProfiles.removeAll { $0.userId == profile.userId }
        } else if mode == .privateChat {
            selectedProfiles = [profile]
        } else {
            selectedProfiles.append(profile)
        }
    }

    private func profileRow(
        _ profile: UserProfileModel,
        isSelected: Bool,
        action: @escaping () -> Void
    ) -> some View {
        Button(action: action) {
            HStack(spacing: 12) {
                AsyncImage(url: profile.profilePhotoUrl.flatMap(URL.init(string:))) { phase in
                    switch phase {
                    case .success(let image):
                        image.resizable().scaledToFill()
                    default:
                        Circle()
                            .fill(Color.orange.opacity(0.14))
                            .overlay(Image(systemName: "person.fill").foregroundStyle(.orange))
                    }
                }
                .frame(width: 42, height: 42)
                .clipShape(Circle())

                VStack(alignment: .leading, spacing: 2) {
                    Text(displayName(for: profile))
                        .font(.headline)
                        .foregroundStyle(.primary)

                    Text("@\(profile.username)")
                        .font(.subheadline)
                        .foregroundStyle(.secondary)
                }

                Spacer()

                if isSelected {
                    Image(systemName: "checkmark.circle.fill")
                        .foregroundStyle(.orange)
                }
            }
            .contentShape(Rectangle())
        }
        .buttonStyle(.plain)
    }

    private func displayName(for profile: UserProfileModel) -> String {
        let fullName = [profile.name, profile.surname]
            .map { $0.trimmingCharacters(in: .whitespacesAndNewlines) }
            .filter { !$0.isEmpty }
            .joined(separator: " ")

        return fullName.isEmpty ? "@\(profile.username)" : fullName
    }
}

private enum ChatCreationMode: String, CaseIterable, Identifiable {
    case privateChat
    case groupChat

    var id: String { rawValue }

    var title: String {
        switch self {
        case .privateChat:
            return "Личный"
        case .groupChat:
            return "Группа"
        }
    }

    var searchFieldTitle: String {
        switch self {
        case .privateChat:
            return "Имя или ник пользователя"
        case .groupChat:
            return "Имя или ник участника"
        }
    }

    var footer: String {
        switch self {
        case .privateChat:
            return "Найдите профиль и выберите человека, с которым нужно открыть диалог."
        case .groupChat:
            return "Найдите и выберите участников группы."
        }
    }
}

private enum ChatCreationError: LocalizedError {
    case noProfileSelected

    var errorDescription: String? {
        "Выберите профиль из списка."
    }
}
