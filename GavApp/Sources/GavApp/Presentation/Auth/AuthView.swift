import SwiftUI

struct AuthView: View {
    @ObservedObject var session: AppSessionViewModel

    @State private var mode: AuthMode = .login
    @State private var email = ""
    @State private var firstName = ""
    @State private var lastName = ""
    @State private var username = ""
    @State private var password = ""
    @State private var passwordConfirmation = ""

    var body: some View {
        VStack(spacing: 0) {
            Spacer(minLength: 28)

            VStack(alignment: .leading, spacing: 28) {
                VStack(alignment: .leading, spacing: 8) {
                    Text("Gav")
                        .font(.system(size: 42, weight: .bold))
                        .foregroundStyle(.white)

                    Text(mode.title)
                        .font(.title3.weight(.semibold))
                        .foregroundStyle(.white.opacity(0.78))
                }

                Picker("Режим", selection: $mode) {
                    ForEach(AuthMode.allCases) { item in
                        Text(item.segmentTitle).tag(item)
                    }
                }
                .pickerStyle(.segmented)

                VStack(spacing: 12) {
                    authField(
                        title: "Email",
                        text: $email,
                        keyboard: .emailAddress,
                        contentType: .emailAddress
                    )

                    secureField(title: "Пароль", text: $password)

                    if mode == .register {
                        authField(
                            title: "Имя",
                            text: $firstName,
                            keyboard: .default,
                            contentType: .givenName,
                            autocapitalization: .words
                        )

                        authField(
                            title: "Фамилия",
                            text: $lastName,
                            keyboard: .default,
                            contentType: .familyName,
                            autocapitalization: .words
                        )

                        authField(
                            title: "Никнейм",
                            text: $username,
                            keyboard: .default,
                            contentType: nil
                        )

                        secureField(title: "Повтор пароля", text: $passwordConfirmation)
                    }
                }

                if let errorMessage = session.errorMessage {
                    Text(errorMessage)
                        .font(.footnote)
                        .foregroundStyle(.red.opacity(0.9))
                        .fixedSize(horizontal: false, vertical: true)
                }

                Button {
                    Task { await submit() }
                } label: {
                    HStack(spacing: 10) {
                        if session.isLoading {
                            ProgressView()
                                .tint(.black)
                        }

                        Text(mode.buttonTitle)
                            .font(.headline)
                    }
                    .frame(maxWidth: .infinity)
                    .frame(height: 50)
                    .foregroundStyle(.black)
                    .background(isSubmitEnabled ? Color.orange : Color.gray, in: RoundedRectangle(cornerRadius: 8))
                }
                .disabled(!isSubmitEnabled || session.isLoading)
            }
            .padding(.horizontal, 22)
            .padding(.vertical, 28)

            Spacer(minLength: 28)
        }
        .background(
            LinearGradient(
                colors: [
                    Color(red: 0.42, green: 0.22, blue: 0.72),
                    .black
                ],
                startPoint: .top,
                endPoint: .bottom
            )
            .ignoresSafeArea()
        )
        .preferredColorScheme(.dark)
        .onChange(of: mode) { _, _ in
            session.errorMessage = nil
        }
    }

    private var isSubmitEnabled: Bool {
        let trimmedEmail = email.trimmingCharacters(in: .whitespacesAndNewlines)
        guard trimmedEmail.contains("@"), password.count >= 6 else { return false }
        guard mode == .register else { return true }

        return password == passwordConfirmation &&
            !firstName.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty &&
            isUsernameValid
    }

    private var isUsernameValid: Bool {
        let clean = normalizedUsername(username)
        return clean.count >= 3 && clean.allSatisfy { $0.isLetter || $0.isNumber || $0 == "_" || $0 == "." }
    }

    private func submit() async {
        let trimmedEmail = email.trimmingCharacters(in: .whitespacesAndNewlines)
        switch mode {
        case .login:
            await session.login(email: trimmedEmail, password: password)
        case .register:
            await session.register(
                email: trimmedEmail,
                password: password,
                firstName: firstName,
                lastName: lastName,
                username: normalizedUsername(username)
            )
        }
    }

    private func normalizedUsername(_ value: String) -> String {
        value
            .trimmingCharacters(in: .whitespacesAndNewlines)
            .trimmingCharacters(in: CharacterSet(charactersIn: "@"))
            .lowercased()
    }

    private func authField(
        title: String,
        text: Binding<String>,
        keyboard: UIKeyboardType,
        contentType: UITextContentType?,
        autocapitalization: TextInputAutocapitalization = .never
    ) -> some View {
        TextField(title, text: text)
            .keyboardType(keyboard)
            .textContentType(contentType)
            .textInputAutocapitalization(autocapitalization)
            .autocorrectionDisabled()
            .padding(.horizontal, 14)
            .frame(height: 48)
            .foregroundStyle(.white)
            .background(.white.opacity(0.11), in: RoundedRectangle(cornerRadius: 8))
    }

    private func secureField(title: String, text: Binding<String>) -> some View {
        SecureField(title, text: text)
            .textContentType(.password)
            .padding(.horizontal, 14)
            .frame(height: 48)
            .foregroundStyle(.white)
            .background(.white.opacity(0.11), in: RoundedRectangle(cornerRadius: 8))
    }
}

private enum AuthMode: String, CaseIterable, Identifiable {
    case login
    case register

    var id: String { rawValue }

    var title: String {
        switch self {
        case .login:
            return "Вход"
        case .register:
            return "Регистрация"
        }
    }

    var segmentTitle: String {
        switch self {
        case .login:
            return "Вход"
        case .register:
            return "Регистрация"
        }
    }

    var buttonTitle: String {
        switch self {
        case .login:
            return "Войти"
        case .register:
            return "Создать аккаунт"
        }
    }
}
