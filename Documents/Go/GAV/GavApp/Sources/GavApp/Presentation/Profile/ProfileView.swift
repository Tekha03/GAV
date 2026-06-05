import SwiftUI

struct ProfileView: View {
    @EnvironmentObject private var appViewModel: AppViewModel
    @ObservedObject var session: AppSessionViewModel
    let uploadService: UploadServiceAPIProtocol

    @State private var selectedDog: AppDog?
    @State private var showSettings = false
    @State private var showAddPostSheet = false
    @State private var showAddDogSheet = false
    @State private var showLogoutConfirmation = false
    @State private var editingDog: AppDog?

    var body: some View {
        NavigationStack {
            ScrollView {
                VStack(spacing: 20) {
                    header
                    dogsRow
                    postsSection
                }
                .padding(.horizontal, 16)
                .padding(.bottom, 24)
            }
            .scrollIndicators(.hidden)
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
            .navigationTitle("")
            .toolbarTitleDisplayMode(.inline)
            .toolbar {
                if appViewModel.canEditProfile {
                    ToolbarItem(placement: .topBarLeading) {
                        Button {
                            showLogoutConfirmation = true
                        } label: {
                            Image(systemName: "rectangle.portrait.and.arrow.right")
                        }
                        .accessibilityLabel("Выйти")
                        .tint(.white)
                    }

                    ToolbarItem(placement: .automatic) {
                        Button {
                            showAddPostSheet = true
                        } label: {
                            Label("Пост", systemImage: "plus.circle.fill")
                        }
                        .tint(.white)
                    }

                    ToolbarItem(placement: .automatic) {
                        Button {
                            showSettings = true
                        } label: {
                            Image(systemName: "gearshape.fill")
                        }
                        .tint(.white)
                    }
                }
            }
            .sheet(isPresented: $showSettings) {
                ProfileSettingsView(
                    appViewModel: appViewModel,
                    settingsService: appViewModel.settingsService,
                    uploadService: uploadService
                )
            }
            .sheet(isPresented: $showAddPostSheet) {
                AddPostView(
                    viewModel: appViewModel,
                    uploadService: uploadService
                )
            }
            .sheet(isPresented: $showAddDogSheet) {
                AddDogView(
                    viewModel: appViewModel,
                    uploadService: uploadService
                )
            }
            .sheet(item: $selectedDog) { dog in
                dogDetail(dog)
            }
            .sheet(item: $editingDog) { dog in
                AddDogView(
                    viewModel: appViewModel,
                    editingDog: dog,
                    uploadService: uploadService
                )
            }
            .confirmationDialog(
                "Выйти из профиля?",
                isPresented: $showLogoutConfirmation,
                titleVisibility: .visible
            ) {
                Button("Выйти", role: .destructive) {
                    session.logout()
                }

                Button("Отмена", role: .cancel) {}
            }
        }
        .preferredColorScheme(.dark)
    }

    private var header: some View {
        VStack(alignment: .leading, spacing: 14) {
            Text(appViewModel.profile.handle)
                .font(.headline)
                .foregroundStyle(.white.opacity(0.95))
                .padding(.top, 12)

            HStack(alignment: .top, spacing: 16) {
                avatar
                    .frame(width: 92, height: 92)

                VStack(alignment: .leading, spacing: 12) {
                    HStack(spacing: 12) {
                        statCard(title: "Подписчики", value: appViewModel.profile.followers)
                        statCard(title: "Подписки", value: appViewModel.profile.following)
                        statCard(title: "Собаки", value: appViewModel.dogs.count)
                    }

                    Text(appViewModel.profile.bio)
                        .font(.footnote)
                        .foregroundStyle(.white.opacity(0.9))
                        .lineLimit(nil)
                        .fixedSize(horizontal: false, vertical: true)
                        .frame(maxWidth: .infinity, alignment: .leading)
                }
            }
            .padding(.top, 4)
        }
        .padding(18)
        .background(
            RoundedRectangle(cornerRadius: 28, style: .continuous)
                .fill(
                    LinearGradient(
                        colors: [
                            Color(red: 0.62, green: 0.45, blue: 0.93).opacity(0.92),
                            Color.black.opacity(0.75)
                        ],
                        startPoint: .topLeading,
                        endPoint: .bottomTrailing
                    )
                )
        )
        .overlay(
            RoundedRectangle(cornerRadius: 28, style: .continuous)
                .stroke(.white.opacity(0.08), lineWidth: 1)
        )
        .shadow(color: .black.opacity(0.25), radius: 18, x: 0, y: 10)
    }

    private var avatar: some View {
        AsyncImage(url: appViewModel.profile.avatarURL) { phase in
            switch phase {
            case .success(let image):
                image.resizable().scaledToFill()
            default:
                ZStack {
                    Circle().fill(.white.opacity(0.12))
                    Image(systemName: "person.fill")
                        .font(.title)
                        .foregroundStyle(.white.opacity(0.9))
                }
            }
        }
        .clipShape(Circle())
        .overlay(Circle().stroke(.white.opacity(0.25), lineWidth: 2))
    }

    private func statCard(title: String, value: Int) -> some View {
        VStack(alignment: .leading, spacing: 3) {
            Text("\(value)")
                .font(.headline.bold())
                .foregroundStyle(.white)

            Text(title)
                .font(.caption2)
                .foregroundStyle(.white.opacity(0.75))
        }
        .frame(maxWidth: .infinity, alignment: .leading)
    }

    private var dogsRow: some View {
        VStack(alignment: .leading, spacing: 12) {
            HStack {
                Text("Собаки")
                    .font(.headline)
                    .foregroundStyle(.white)

                Spacer()

                if appViewModel.canEditProfile {
                    Button {
                        showAddDogSheet = true
                    } label: {
                        Label("Добавить", systemImage: "plus")
                            .font(.subheadline.weight(.semibold))
                    }
                    .tint(.white)
                }
            }

            ScrollView(.horizontal, showsIndicators: false) {
                HStack(spacing: 14) {
                    ForEach(appViewModel.dogs) { dog in
                        Button {
                            selectedDog = dog
                        } label: {
                            dogCard(dog)
                        }
                        .buttonStyle(.plain)
                    }
                }
                .padding(.vertical, 2)
            }
        }
    }

    private var postsSection: some View {
        VStack(alignment: .leading, spacing: 12) {
            Text("Посты")
                .font(.headline)
                .foregroundStyle(.white)

            ForEach(appViewModel.posts) { post in
                PostView(post: post) {
                    appViewModel.toggleLike(postID: post.id)
                }
            }
        }
    }

    @ViewBuilder
    private func dogDetail(_ dog: AppDog) -> some View {
        VStack(alignment: .leading, spacing: 16) {
            AsyncImage(url: dog.photoURL) { phase in
                switch phase {
                case .success(let image):
                    image.resizable().scaledToFill()
                default:
                    RoundedRectangle(cornerRadius: 24)
                        .fill(Color.white.opacity(0.12))
                        .overlay(Image(systemName: "dog.fill").font(.largeTitle))
                }
            }
            .frame(height: 280)
            .clipShape(RoundedRectangle(cornerRadius: 24))

            Text(dog.name)
                .font(.largeTitle.bold())

            Text("\(dog.breed) · \(dog.ageText)")
                .foregroundStyle(.secondary)

            Text(dog.notes)

            Text("Характер: \(dog.mood.title)")
                .foregroundStyle(dog.mood.color)

            if appViewModel.canEditProfile {
                Button {
                    editingDog = dog
                    selectedDog = nil
                } label: {
                    Label("Редактировать", systemImage: "pencil")
                        .font(.headline.weight(.semibold))
                        .frame(maxWidth: .infinity)
                        .padding(.vertical, 12)
                        .background(.white.opacity(0.12), in: RoundedRectangle(cornerRadius: 16))
                }
                .buttonStyle(.plain)
            }

            Spacer()
        }
        .padding(20)
        .background(.black)
        .preferredColorScheme(.dark)
    }

    private func dogCard(_ dog: AppDog) -> some View {
        VStack(alignment: .leading, spacing: 8) {
            AsyncImage(url: dog.photoURL) { phase in
                switch phase {
                case .success(let image):
                    image
                        .resizable()
                        .scaledToFill()
                default:
                    ZStack {
                        LinearGradient(
                            colors: [
                                Color.white.opacity(0.18),
                                Color.white.opacity(0.06)
                            ],
                            startPoint: .topLeading,
                            endPoint: .bottomTrailing
                        )
                        Image(systemName: "dog.fill")
                            .font(.title2)
                            .foregroundStyle(.white.opacity(0.85))
                    }
                }
            }
            .frame(width: 132, height: 180)
            .clipped()
            .overlay(
                LinearGradient(
                    colors: [.clear, .black.opacity(0.82)],
                    startPoint: .center,
                    endPoint: .bottom
                )
            )
            .overlay(alignment: .bottomLeading) {
                VStack(alignment: .leading, spacing: 4) {
                    Text(dog.name)
                        .font(.headline.weight(.semibold))
                        .foregroundStyle(.white)
                    Text(dog.breed)
                        .font(.caption)
                        .foregroundStyle(.white.opacity(0.8))
                }
                .padding(10)
            }
            .clipShape(RoundedRectangle(cornerRadius: 22, style: .continuous))
            .shadow(color: .black.opacity(0.22), radius: 10, x: 0, y: 6)
        }
    }
}
