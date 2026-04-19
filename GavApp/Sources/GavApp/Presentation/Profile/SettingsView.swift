import SwiftUI

@available(macOS 12.0, *)
struct ProfileSettingsView: View {
    let viewModel: ProfileViewModel
    @Environment(\\.presentationMode) private var presentationMode

    @State private var name: String = ""
    @State private var surname: String = ""
    @State private var bio: String = ""

    @State private var profilePrivacy: Bool = true
    @State private var showLocation: Bool = false
    @State private var allowMessages: Bool = true

    @State private var isProfileEditing = true

    var body: some View {
        NavigationStack {
            Form {
                Section("Профиль") {
                    TextField("Имя", text: $name)
                    TextField("Фамилия", text: $surname)
                    TextField("Био", text: $bio)

                    Toggle("Приватный профиль", isOn: $profilePrivacy)
                    Toggle("Показывать местоположение", isOn: $showLocation)
                    Toggle("Разрешить сообщения", isOn: $allowMessages)
                }

                Section("Собаки") {
                    Button("Добавить собаку") {
                        isProfileEditing = false
                    }
                }

                Button("Сохранить") {
                    Task { await saveProfile() }
                }
            }
            .navigationTitle("Настройки профиля")
            .onAppear {
                setupInitialValues()
            }
            .sheet(isPresented: $isProfileEditing) {
                AddDogView(viewModel: viewModel)
            }
        }
    }

    private func setupInitialValues() {
        guard let profile = viewModel.profile else { return }

        name = profile.name
        surname = profile.surname
        bio = profile.bio

        profilePrivacy = viewModel.currentUser?.settings.profilePrivacy ?? true
        showLocation = viewModel.currentUser?.settings.showLocation ?? false
        allowMessages = viewModel.currentUser?.settings.allowMessages ?? true
    }

    private func saveProfile() async {
        let profileInput = UpdateProfileInput(
            name: name,
            surname: surname,
            bio: bio
        )
        let userSettings = UserSettings(
            userId: viewModel.currentUserId,
            profilePrivacy: profilePrivacy,
            showLocation: showLocation,
            allowMessages: allowMessages
        )

        await viewModel.updateProfile(input: profileInput)
        await viewModel.updateUserSettings(settings: userSettings)
        presentationMode.wrappedValue.dismiss()
    }
}
