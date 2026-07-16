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
    @State private var screenState: AppScreenState = .loading(message: "Загружаем настройки...")
    @State private var actionErrorMessage: String?

    var body: some View {
        NavigationStack {
            Group {
                switch screenState {
                case .loading, .error, .offline:
                    AppStatusView(
                        state: screenState,
                        retryAction: loadInitialValues
                    )

                case .content:
                    settingsContent
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
                loadInitialValues()
            }
            .task(id: selectedItem) {
                await loadAvatar()
            }
        }
    }

    private var settingsForm: some View {
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
                PhotosPicker(
                    selection: $selectedItem,
                    matching: .images
                ) {
                    Label(
                        "Выбрать фото",
                        systemImage: "photo"
                    )
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

            if let actionErrorMessage {
                Section {
                    HStack(spacing: 10) {
                        Image(
                            systemName: "exclamationmark.circle.fill"
                        )
                        .foregroundStyle(.orange)

                        Text(actionErrorMessage)
                            .foregroundStyle(.red)

                        Spacer()

                        Button {
                            self.actionErrorMessage = nil
                        } label: {
                            Image(systemName: "xmark")
                        }
                        .buttonStyle(.plain)
                    }
                }
            }

            Section("Приватность") {
                Toggle(
                    "Приватный профиль",
                    isOn: $isPrivateProfile
                )

                Toggle(
                    "Показывать на карте",
                    isOn: $showOnMap
                )
            }

            Section {
                Button {
                    Task {
                        await saveProfile()
                    }
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
    }

    private func loadInitialValues() {
        screenState = .loading(
            message: "Загружаем настройки..."
        )

        let parts = appViewModel.profile.fullName.split(
            separator: " ",
            maxSplits: 1,
            omittingEmptySubsequences: true
        )

        firstName = parts.first.map(String.init) ?? ""
        lastName = parts.count > 1
            ? String(parts[1])
            : ""

        username = appViewModel.profile.handle
            .trimmingCharacters(
                in: CharacterSet(charactersIn: "@")
            )

        bio = appViewModel.profile.bio
        isPrivateProfile = false
        showOnMap = true
        screenState = .content
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
            actionErrorMessage = error.localizedDescription
        }
    }

    private func saveProfile() async {
        isSaving = true
        actionErrorMessage = nil
        defer { isSaving = false }

        let fullName = [firstName, lastName]
            .filter { !$0.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty }
            .joined(separator: " ")
        let cleanUsername = normalizedUsername(username)

        guard cleanUsername.count >= 3 else {
            actionErrorMessage = "Никнейм должен быть не короче 3 символов"
            return
        }

        let uploadedAvatar = await uploadAvatarIfNeeded()
        if actionErrorMessage != nil {
            return
        }

        do {
            try await saveRemoteProfile(
                username: cleanUsername,
                profilePhotoUrl: uploadedAvatar?.rawURL
            )
        } catch APIError.invalidResponse(let statusCode) where statusCode == 409 {
            actionErrorMessage = "Этот никнейм уже занят"
            return
        } catch {
            actionErrorMessage = "Не удалось сохранить профиль"
            return
        }

        if !fullName.isEmpty {
            appViewModel.profile.fullName = fullName
        }
        if let uploadedAvatar {
            appViewModel.profile.avatarURL = uploadedAvatar.resolvedURL
        }
        appViewModel.profile.handle = "@\(cleanUsername)"
        appViewModel.profile.bio = bio
        dismiss()
    }

    private func uploadAvatarIfNeeded() async -> UploadedMedia? {
        actionErrorMessage = error.localizedDescription
        guard let selectedAvatarData else { return nil }

        do {
            let media = try await uploadService.uploadAvatar(
                selectedAvatarData,
                mimeType: "image/jpeg"
            )
            return UploadedMedia(
                rawURL: media.url,
                resolvedURL: MediaURLResolver.resolve(media.url)
            )
        } catch {
            actionErrorMessage = "Не удалось загрузить аватар"
            return nil
        }
    }

    private func saveRemoteProfile(username: String, profilePhotoUrl: String?) async throws {
        let cleanFirstName = firstName.trimmingCharacters(in: .whitespacesAndNewlines)
        let cleanLastName = lastName.trimmingCharacters(in: .whitespacesAndNewlines)
        let input = UpdateProfileInput(
            name: cleanFirstName,
            surname: cleanLastName,
            username: username,
            profilePhotoUrl: profilePhotoUrl,
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
                profilePhotoUrl: profilePhotoUrl,
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

private struct UploadedMedia {
    let rawURL: String
    let resolvedURL: URL?
}
