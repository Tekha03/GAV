import SwiftUI
import PhotosUI
#if os(iOS)
import UIKit
#endif

struct ProfileSettingsView: View {
    let appViewModel: AppViewModel
    let uploadService: UploadServiceAPIProtocol
    @Environment(\.dismiss) private var dismiss

    @State private var firstName: String = ""
    @State private var lastName: String = ""
    @State private var username: String = ""
    @State private var bio: String = ""
    @State private var isPrivateProfile: Bool = false
    @State private var showOnMap: Bool = true
    @State private var selectedItem: PhotosPickerItem?
    @State private var selectedAvatarData: Data?
    @State private var avatarPreview: Image?
    @State private var isLoadingAvatar = false
    @State private var isSaving = false
    @State private var errorMessage: String?

    var body: some View {
        NavigationStack {
            Form {
                Section("Профиль") {
                    TextField("Имя", text: $firstName)
                    TextField("Фамилия", text: $lastName)
                    TextField("Никнейм", text: $username)
                        .textInputAutocapitalization(.never)
                        .autocorrectionDisabled()
                    TextField("Био", text: $bio, axis: .vertical)
                        .lineLimit(3...6)
                }

                Section("Фото") {
                    PhotosPicker(selection: $selectedItem, matching: .images) {
                        Label("Выбрать фото", systemImage: "photo")
                    }

                    if let avatarPreview {
                        avatarPreview
                            .resizable()
                            .scaledToFill()
                            .frame(width: 96, height: 96)
                            .clipShape(Circle())
                    }

                    if isLoadingAvatar {
                        ProgressView()
                    }
                }

                if let errorMessage {
                    Section {
                        Text(errorMessage)
                            .foregroundStyle(.red)
                    }
                }

                Section("Приватность") {
                    Toggle("Приватный профиль", isOn: $isPrivateProfile)
                    Toggle("Показывать на карте", isOn: $showOnMap)
                }

                Section {
                    Button {
                        Task { await saveProfile() }
                    } label: {
                        if isSaving {
                            ProgressView()
                        } else {
                            Text("Сохранить")
                        }
                    }
                    .disabled(isSaving)
                }
            }
            .navigationTitle("Настройки")
            .toolbar {
                ToolbarItem(placement: .cancellationAction) {
                    Button("Отмена") {
                        dismiss()
                    }
                }
            }
            .task {
                setupInitialValues()
            }
            .task(id: selectedItem) {
                await loadAvatar()
            }
        }
    }

    private func setupInitialValues() {
        let parts = appViewModel.profile.fullName.split(separator: " ", maxSplits: 1, omittingEmptySubsequences: true)
        firstName = parts.first.map(String.init) ?? ""
        lastName = parts.count > 1 ? String(parts[1]) : ""
        username = appViewModel.profile.handle.trimmingCharacters(in: CharacterSet(charactersIn: "@"))
        bio = appViewModel.profile.bio
        isPrivateProfile = false
        showOnMap = true
    }

    private func loadAvatar() async {
        guard let selectedItem else { return }
        isLoadingAvatar = true
        defer { isLoadingAvatar = false }

        do {
            if let data = try await selectedItem.loadTransferable(type: Data.self) {
                selectedAvatarData = data
                #if os(macOS)
                if let nsImage = NSImage(data: data) {
                    avatarPreview = Image(nsImage: nsImage)
                }
                #elseif os(iOS)
                if let uiImage = UIImage(data: data) {
                    avatarPreview = Image(uiImage: uiImage)
                }
                #endif
            }
        } catch {
            selectedAvatarData = nil
            avatarPreview = nil
        }
    }

    private func saveProfile() async {
        isSaving = true
        errorMessage = nil
        defer { isSaving = false }

        let fullName = [firstName, lastName]
            .filter { !$0.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty }
            .joined(separator: " ")
        let cleanUsername = normalizedUsername(username)

        guard cleanUsername.count >= 3 else {
            errorMessage = "Никнейм должен быть не короче 3 символов"
            return
        }

        let avatarURL = await uploadAvatarIfNeeded()
        if errorMessage != nil { return }

        do {
            try await saveRemoteProfile(username: cleanUsername)
        } catch APIError.invalidResponse(let statusCode) where statusCode == 409 {
            errorMessage = "Этот никнейм уже занят"
            return
        } catch {
            errorMessage = "Не удалось сохранить профиль"
            return
        }

        if !fullName.isEmpty {
            appViewModel.profile.fullName = fullName
        }
        if let avatarURL {
            appViewModel.profile.avatarURL = avatarURL
        }
        appViewModel.profile.handle = "@\(cleanUsername)"
        appViewModel.profile.bio = bio
        dismiss()
    }

    private func uploadAvatarIfNeeded() async -> URL? {
        guard let selectedAvatarData else { return nil }

        do {
            let media = try await uploadService.uploadAvatar(
                selectedAvatarData,
                mimeType: "image/jpeg"
            )
            return MediaURLResolver.resolve(media.url)
        } catch {
            errorMessage = "Не удалось загрузить аватар"
            return nil
        }
    }

    private func saveRemoteProfile(username: String) async throws {
        let cleanFirstName = firstName.trimmingCharacters(in: .whitespacesAndNewlines)
        let cleanLastName = lastName.trimmingCharacters(in: .whitespacesAndNewlines)
        let input = UpdateProfileInput(
            name: cleanFirstName,
            surname: cleanLastName,
            username: username,
            bio: bio
        )

        do {
            try await appViewModel.profileService.update(
                userID: appViewModel.currentUserId,
                input: input
            )
        } catch APIError.invalidResponse(let statusCode) where statusCode == 404 {
            let profile = CreateProfileInput(
                name: cleanFirstName,
                surname: cleanLastName,
                username: username,
                bio: bio
            ).toModel(userID: appViewModel.currentUserId)
            _ = try await appViewModel.profileService.create(
                userID: appViewModel.currentUserId,
                model: profile
            )
        }
    }

    private func normalizedUsername(_ value: String) -> String {
        value
            .trimmingCharacters(in: .whitespacesAndNewlines)
            .trimmingCharacters(in: CharacterSet(charactersIn: "@"))
            .lowercased()
    }
}
