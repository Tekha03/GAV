import SwiftUI

enum TabItem: Hashable {
    case profile
    case feed
    case map
    case vaccination
    case chat
    case search
}

struct AppView: View {
    @EnvironmentObject private var appViewModel: AppViewModel
    let session: AppSessionViewModel

    var body: some View {
        TabView {
            ProfileView(
                session: session,
                uploadService: appViewModel.uploadService
            )
                .tabItem {
                    Label("Профиль", systemImage: "person.crop.circle.fill")
                }

            FeedView()
                .tabItem {
                    Label("Лента", systemImage: "house.fill")
                }

            MapView()
                .tabItem {
                    Label("Карта", systemImage: "map.fill")
                }

            VaccinationTabView()
                .tabItem {
                    Label("Прививки", systemImage: "syringe.fill")
                }

            ChatListView()
                .tabItem {
                    Label("Чаты", systemImage: "bubble.left.and.bubble.right.fill")
                }
        }
        .task(id: appViewModel.currentUserId) {
            await appViewModel.loadCurrentProfile()
        }
    }
}
