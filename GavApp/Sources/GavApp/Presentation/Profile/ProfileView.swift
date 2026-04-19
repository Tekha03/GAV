import SwiftUI

@available(macOS 12.0, *)
struct ProfileView: View {
    @StateObject private var viewModel: ProfileViewModel
    let userId: UUID

    @State private var selectedDog: Dog?

    var body: some View {
        Group {
            if viewModel.isLoading {
                ProgressView()
            } else if let profile = viewModel.profile {
                content(profile)
            }
        }
        .task {
            await viewModel.load(userId: userId)
        }
    }

    @ViewBuilder
    private func content(_ profile: UserProfile) -> some View {
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
                if viewModel.isOwner {
                    EditableDogDetailView(
                        dog: dog,
                        onSave: { updatedDog in
                            Task { await viewModel.updateDog(dogId: dog.id, input: updatedDog) }
                        }
                    )
                } else {
                    DogDetailView(dog: dog)
                }
            }
            .scrollContentBackground(.hidden)
        }
        .environment(\.colorScheme, .dark)
    }

    private var background: some View {
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

    private func header(_ profile: UserProfile) -> some View {
        VStack(spacing: 10) {
            HStack(alignment: .center, spacing: 16) {
                Image(profile.profilePhotoUrl.isEmpty ? "personPlaceholder" : profile.profilePhotoUrl)
                    .resizable()
                    .scaledToFill()
                    .frame(width: 78, height: 78)
                    .clipShape(Circle())
                    .overlay(
                        Circle().stroke(Color(red: 220/255, green: 255/255, blue: 5/255), lineWidth: 2)
                    )
                    .shadow(radius: 4)

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
            .padding(.horizontal)

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

    private func dogsRow() -> some View {
        VStack(alignment: .leading, spacing: 8) {
            Text("Мои собаки")
                .font(.headline)
                .foregroundColor(Color(red: 60/255, green: 30/255, blue: 100/255))
                .padding(.horizontal)

            ScrollView(.horizontal, showsIndicators: false) {
                HStack(spacing: 16) {
                    ForEach(viewModel.dogs) { dog in
                        if viewModel.isOwner {
                            EditableDogStoryView(
                                dog: dog,
                                onEditTapped: { editedDog in
                                    selectedDog = editedDog
                                }
                            )
                        } else {
                            DogStoryView(
                                dog: dog,
                                onTapped: { storyDog in
                                    selectedDog = storyDog
                                }
                            )
                        }
                    }

                    // Добавляем плюс‑стори для владельца
                    if viewModel.isOwner {
                        AddDogStoryView {
                            // TODO: Открыть AddDogView или модальное окно
                        }
                    }
                }
                .padding(.horizontal)
            }
        }
        .padding(.vertical, 12)
    }

    private func posts() -> some View {
        VStack(alignment: .leading, spacing: 16) {
            if viewModel.isOwner {
                Button("Добавить пост") {
                    Task {
                        await viewModel.createPost(
                            content: "Тестовый пост",
                            imageUrl: nil
                        )
                    }
                }
                .padding(.horizontal)
            }

            Text("Посты")
                .font(.headline)
                .foregroundColor(.white)
                .padding(.horizontal)
                .padding(.top, 8)

            ForEach(viewModel.posts) { post in
                PostView(
                    post: post,
                    toggleLike: {
                        Task { await viewModel.toggleLike(for: post) }
                    }
                )
            }
        }
        .padding(.bottom, 20)
    }
}