import SwiftUI

@main
struct GavApp: App {
    @StateObject private var appViewModel: AppViewModel
    @StateObject private var sessionViewModel: AppSessionViewModel

    init() {
        let container = DependencyContainer()
        _appViewModel = StateObject(wrappedValue: container.appViewModel)
        _sessionViewModel = StateObject(
            wrappedValue: AppSessionViewModel(
                authService: container.authService,
                authManager: container.authManager,
                appViewModel: container.appViewModel
            )
        )
    }

    var body: some Scene {
        WindowGroup {
            Group {
                if sessionViewModel.isLoading {
                    ProgressView()
                } else if sessionViewModel.isAuthenticated {
                    AppView(session: sessionViewModel)
                } else {
                    AuthView(session: sessionViewModel)
                }
            }
            .environmentObject(appViewModel)
            .task {
                await sessionViewModel.restoreSavedSessionIfNeeded()
            }
        }
    }
}
