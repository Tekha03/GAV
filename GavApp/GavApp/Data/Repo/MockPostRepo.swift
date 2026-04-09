import Foundation

final class MockPostRepository: PostRepository {

    private var posts: [Post] = []

    init() {
        let userId = UUID()

        posts = [
            Post(id: UUID(), userId: userId, date: Date(),
                 content: "Сегодня гуляли в парке! 🐕",
                 imageUrl: "dog_walk",
                 likesCount: 15,
                 commentsCount: 3,
                 isLiked: false),

            Post(id: UUID(), userId: userId, date: Date().addingTimeInterval(-86400),
                 content: "Новая игрушка 🧸",
                 imageUrl: "toy",
                 likesCount: 8,
                 commentsCount: 1,
                 isLiked: true)
        ]
    }

    func fetchPosts(userId: UUID) async -> [Post] {
        posts
    }

    func likePost(postId: UUID) async -> Bool {
        guard let i = posts.firstIndex(where: { $0.id == postId }) else { return false }

        posts[i].isLiked.toggle()
        posts[i].likesCount += posts[i].isLiked ? 1 : -1
        return true
    }
}
