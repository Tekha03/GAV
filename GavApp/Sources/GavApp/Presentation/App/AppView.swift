import SwiftUI
import Combine

enum TabItem: Hashable {
    case profile
    case feed
    case map
    case vaccination 
    case chat
}

@available(macOS 12.0, *)
struct AppView: View {
    @State private var selection: TabItem = .profile

    var body: some View {
        TabView(selection: $selection) {
            ProfileView(
                viewModel: profileViewModel,
                userId: appViewModel.currentUser?.id ?? UUID()
            )
            .tabItem {
                Label("Профиль", systemImage: "person.crop.circle.fill")
            }
            .tag(TabItem.profile)

            FeedView(
                viewModel: feedViewModel
            )
            .tabItem {
                Label("Лента", systemImage: "house.fill")
            }
            .tag(TabItem.feed)

            MapView(
                viewModel: mapViewModel
            )
            .tabItem {
                Label("Карта", systemImage: "map.fill")
            }
            .tag(TabItem.map)

            VaccinationTabView(
                viewModel: vaccinationViewModel
            )
            .tabItem {
                Label("Прививки", systemImage: "syringe.fill")
            }
            .tag(TabItem.vaccination)

            ChatListView(
                viewModel: chatListViewModel
            )
            .tabItem {
                Label("Чаты", systemImage: "bubble.left.and.bubble.right.fill")
            }
            .tag(TabItem.chat)
        }
    }

    @StateObject private var profileViewModel: ProfileViewModel
    @StateObject private var feedViewModel: FeedViewModel
    @StateObject private var mapViewModel: MapViewModel
    @StateObject private var vaccinationViewModel: VaccinationTabViewModel
    @StateObject private var chatListViewModel: ChatListViewModel

    @ObservedObject private var appViewModel: AppViewModel

    init(
        profileViewModel: ProfileViewModel,
        feedViewModel: FeedViewModel,
        mapViewModel: MapViewModel,
        vaccinationViewModel: VaccinationTabViewModel,
        chatListViewModel: ChatListViewModel,
        appViewModel: AppViewModel
    ) {
        self._profileViewModel = StateObject(
            wrappedValue: profileViewModel
        )
        self._feedViewModel = StateObject(
            wrappedValue: feedViewModel
        )
        self._mapViewModel = StateObject(
            wrappedValue: mapViewModel
        )
        self._vaccinationViewModel = StateObject(
            wrappedValue: vaccinationViewModel
        )
        self._chatListViewModel = StateObject(
            wrappedValue: chatListViewModel
        )
        self._appViewModel = ObservedObject(
            wrappedValue: appViewModel
        )
    }
}