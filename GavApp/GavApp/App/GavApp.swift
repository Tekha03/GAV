//
//  GavApp.swift
//  
//
//  Created by Виктория Кашуркина on 01.04.2026.
//

import Foundation
import Combine
import SwiftUI

@main
struct GavApp: App {
    var body: some Scene {
        // Shows ProfileView via RootView (ensure ProfileViewModel/ProfileView and repositories exist).
        WindowGroup {
            RootView()
        }
    }
}

private struct RootView: View {
    // Stable userId for the lifetime of RootView
    private let userId = UUID()
    // Keep ViewModel as StateObject so it isn't recreated on re-render
    @StateObject private var vm = ProfileViewModel(
        profileRepo: MockProfileRepository(),
        dogRepo: MockDogRepository(),
        postRepo: MockPostRepository()
    )

    var body: some View {
        ProfileView(vm: vm, userId: userId)
    }
}
