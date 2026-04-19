import SwiftUI

@main
@available(macOS 12.0, *)
struct GavApp: App {
    private let container = DependencyContainer()

    var body: some Scene {
        WindowGroup {
            AppView(
                profileViewModel: container.profileViewModel,
                feedViewModel: container.feedViewModel,
                mapViewModel: container.mapViewModel,
                vaccinationViewModel: container.vaccinationViewModel,
                chatListViewModel: container.chatListViewModel,
                appViewModel: container.appViewModel
            )
            .environment(\.colorScheme, .dark)
        }
    }
}