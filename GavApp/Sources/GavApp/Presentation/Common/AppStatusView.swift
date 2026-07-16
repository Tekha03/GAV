import SwiftUI

struct AppStatusView: View {
    let state: AppScreenState
    var retryAction: (() -> Void)?

    var body: some View {
        switch state {
            case .content:
                EmptyView()

            case .loading(let message):
                VStack(spacing: 14) {
                    ProgressView()
                        .controlSize(.large)

                    Text(message)
                        .font(.subheadline)
                        .foregroundStyle(.secondary)
                }
                .frame(maxWidth: .infinity, maxHeight: .infinity)

            case .error(let message):
                statusContent(
                    icon: "exclamationmark.triangle.fill",
                    title: "Что-то пошло не так",
                    message: message
                )

            case .offline:
                statusContent(
                    icon: "wifi.slash",
                    title: "Нет подключения",
                    message: "Проверьте подключение к интернету и попробуйте снова"
                )
        }
    }

    private func statusContent(
        icon: String,
        title: String,
        message: String
    ) -> some View {
        VStack(spacing: 14) {
            Image(systemName: icon)
                .font(.system(size: 42, weight: .semibold))
                .foregroundStyle(.secondary)

            Text(title)
                .font(.headline)
                .multilineTextAlignment(.center)

            Text(message)
                .font(.subheadline)
                .foregroundStyle(.secondary)
                .multilineTextAlignment(.center)
                .frame(maxWidth: 300)

            if let retryAction {
                Button("Попробовать снова", action: retryAction)
                    .buttonStyle(.borderedProminent)
                    .padding(.top, 4)
            }
        }
        .padding(24)
        .frame(maxWidth: .infinity, maxHeight: .infinity)
    }
}
