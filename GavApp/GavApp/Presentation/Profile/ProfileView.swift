//
//  ProfileView.swift
//
//
//  Created by Виктория Кашуркина on 01.04.2026.
//

import SwiftUI
import Combine
import Foundation

struct ProfileView: View {

    @ObservedObject var vm: ProfileViewModel
    let userId: Foundation.UUID

    @State private var selectedDog: Dog?

    var body: some View {
        Group {
            if vm.isLoading {
                ProgressView()
            } else if let profile = vm.profile {
                content(profile)
            }
        }
        .task {
            await vm.load(userId: userId)
        }
    }
}

private extension ProfileView {

    func content(_ profile: UserProfile) -> some View {
        NavigationStack {
            ZStack(alignment: .top) {

                Color.black.ignoresSafeArea()
                
                background

                ScrollView {
                    VStack(spacing: 0) {
                        header(profile)
                        dogsRow()
                        posts()
                    }
                }
            }
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .principal) {
                    Text(profile.username)
                        .font(.caption)
                        .foregroundColor(Color(red: 60/255, green: 30/255, blue: 100/255))
                }
            }
            .sheet(item: $selectedDog) { dog in
                DogDetailView(dog: dog)
            }
            .scrollContentBackground(.hidden)
        }
        .environment(\.colorScheme, .dark)
    }

    var background: some View {
        VStack(spacing: 0) {

            Color(red: 223/255, green: 199/255, blue: 242/255)
                .frame(height: 260)

            LinearGradient(
                colors: [
                    Color(red: 223/255, green: 199/255, blue: 242/255),
                    Color.black
                ],
                startPoint: .top,
                endPoint: .bottom
            )
            .frame(height: 180)
        }
        .ignoresSafeArea(edges: .top)
    }

    func header(_ profile: UserProfile) -> some View {
        VStack(spacing: 10) {

            HStack(alignment: .center, spacing: 16) {

                Image(profile.profilePhotoUrl.isEmpty ? "personPlaceholder" : profile.profilePhotoUrl)
                    .resizable()
                    .scaledToFill()
                    .frame(width: 78, height: 78)
                    .clipShape(Circle())
                    .overlay(Circle().stroke(Color(red: 220/255, green: 255/255, blue: 5/255), lineWidth: 2))
                    .shadow(radius: 4)

                // 📦 правая часть
                HStack(alignment: .top, spacing: 16) {

                    VStack(alignment: .leading, spacing: 4) {
                        Text(profile.name)
                            .font(.headline)
                            .foregroundColor(Color(red: 60/255, green: 30/255, blue: 100/255))

                        Text(profile.surname)
                            .font(.subheadline)
                            .foregroundColor(Color(red: 60/255, green: 30/255, blue: 100/255))
                    }

                    Spacer()

                    HStack(spacing: 20) {

                        VStack {
                            Text("\(profile.followersCount)")
                                .font(.headline)
                                .foregroundColor(Color(red: 60/255, green: 30/255, blue: 100/255))
                            Text("подписчики")
                                .font(.caption2)
                                .foregroundColor(Color(red: 60/255, green: 30/255, blue: 100/255))
                        }

                        VStack {
                            Text("\(profile.followingCount)")
                                .font(.headline)
                                .foregroundColor(Color(red: 60/255, green: 30/255, blue: 100/255))
                            Text("подписки")
                                .font(.caption2)
                                .foregroundColor(Color(red: 60/255, green: 30/255, blue: 100/255))
                        }
                    }
                }
            }
            .padding(.horizontal)


            // био
            if !profile.bio.isEmpty {
                Text(profile.bio)
                    .font(.body)
                    .foregroundColor(Color(red: 60/255, green: 30/255, blue: 100/255))
                    .frame(maxWidth: .infinity, alignment: .leading)
                    .padding(.horizontal)
                    .padding(.top, 4)
            }
        }
        .padding(.top, 6)
        .padding(.bottom, 12)
    }

    func dogsRow() -> some View {
        VStack(alignment: .leading, spacing: 8) {
            Text("Мои собаки")
                .font(.headline)
                .foregroundColor(Color(red: 60/255, green: 30/255, blue: 100/255))
                .padding(.horizontal)
            
            ScrollView(.horizontal, showsIndicators: false) {
                HStack(spacing: 16) {
                    ForEach(vm.dogs) { dog in
                        DogStoryView(dog: dog)
                            .onTapGesture {
                                selectedDog = dog
                            }
                    }
                }
                .padding(.horizontal)
            }
        }
        .padding(.vertical, 12)
    }

    func posts() -> some View {
        VStack(alignment: .leading, spacing: 16) {
            Text("Посты")
                .font(.headline)
                .foregroundColor(.white)
                .padding(.horizontal)
                .padding(.top, 16)
            
            ForEach(vm.posts) { post in
                PostView(post: post) {
                    Task {
                        await vm.toggleLike(for: post)
                    }
                }
            }
        }
        .padding(.bottom, 20)
    }
}


