import SwiftUI

struct PostView: View {
    @EnvironmentObject private var appViewModel: AppViewModel
    let post: AppPost
    let onLike: () -> Void
    let onComment: () -> Void

    private var isLiked: Bool {
        post.isLiked || appViewModel.likedPostIDs.contains(post.id)
    }

    var body: some View {
        VStack(alignment: .leading, spacing: 12) {
            HStack {
                AsyncImage(url: post.authorPhotoURL) { phase in
                    switch phase {
                    case .success(let image):
                        image.resizable().scaledToFill()
                    default:
                        Circle()
                            .fill(.white.opacity(0.12))
                            .overlay(Image(systemName: "person.fill"))
                    }
                }
                .frame(width: 40, height: 40)
                .clipShape(Circle())

                VStack(alignment: .leading, spacing: 2) {
                    Text(post.authorName)
                        .font(.subheadline.bold())
                    Text(post.authorHandle)
                        .font(.caption)
                        .foregroundStyle(.secondary)
                }

                Spacer()

                Text(post.createdAt.formatted(date: .abbreviated, time: .shortened))
                    .font(.caption2)
                    .foregroundStyle(.secondary)
            }

            Text(post.content)
                .font(.body)

            if let imageURL = post.imageURL {
                AsyncImage(url: imageURL) { phase in
                    switch phase {
                    case .success(let image):
                        image.resizable().scaledToFill()
                    default:
                        RoundedRectangle(cornerRadius: 14)
                            .fill(.white.opacity(0.08))
                            .overlay(Image(systemName: "photo"))
                    }
                }
                .frame(height: 240)
                .clipShape(RoundedRectangle(cornerRadius: 14))
            }

            HStack(spacing: 18) {
                Button {
                    withAnimation(.spring(response: 0.28, dampingFraction: 0.55)) {
                        onLike()
                    }
                } label: {
                    Image(systemName: isLiked ? "heart.fill" : "heart")
                        .foregroundStyle(isLiked ? .red : .white.opacity(0.85))
                        .scaleEffect(isLiked ? 1.18 : 1.0)
                        .contentTransition(.symbolEffect(.replace))
                }
                .buttonStyle(.plain)

                Text("\(post.likes)")
                .foregroundStyle(.white)

                Button {
                    onComment()
                } label: {
                    HStack(spacing: 6) {
                        Image(systemName: "bubble.right")
                        Text("\(post.comments)")
                    }
                    .foregroundStyle(.secondary)
                }
                .buttonStyle(.plain)

                Spacer()
            }
            .font(.subheadline)
        }
        .padding(14)
        .background(.thinMaterial, in: RoundedRectangle(cornerRadius: 18))
    }
}
