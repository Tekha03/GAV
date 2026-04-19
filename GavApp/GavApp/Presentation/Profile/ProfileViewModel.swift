import Foundation
import Combine
import SwiftUI

@MainActor
final class ProfileViewModel: ObservableObject {

    @Published var profile: UserProfile?
    @Published var dogs: [Dog] = []
    @Published var posts: [Post] = []
    @Published var isLoading = false

    private let profileRepo: ProfileRepository
    private let dogRepo: DogRepository
    private let postRepo: PostRepository

    init(profileRepo: ProfileRepository,
         dogRepo: DogRepository,
         postRepo: PostRepository) {

        self.profileRepo = profileRepo
        self.dogRepo = dogRepo
        self.postRepo = postRepo
    }

    func load(userId: UUID) async {
        isLoading = true
        defer { isLoading = false }

        async let p = profileRepo.fetchProfile(userId: userId)
        async let d = dogRepo.fetchDogs(ownerId: userId)
        async let posts = postRepo.fetchPosts(userId: userId)

        let (profile, dogs, postsData) = await (p, d, posts)

        self.profile = profile
        self.dogs = dogs
        self.posts = postsData
    }

    func toggleFollow() async {
        guard let id = profile?.userId else { return }

        let ok = await profileRepo.followUser(userId: id)
        guard ok else { return }

        profile?.isFollowed.toggle()
        profile?.followersCount += profile!.isFollowed ? 1 : -1
    }
    
    func toggleLike(for post: Post) async {
        let ok = await postRepo.likePost(postId: post.id)
        guard ok else { return }

        if let i = posts.firstIndex(where: { $0.id == post.id }) {
            posts[i].isLiked.toggle()
            posts[i].likesCount += posts[i].isLiked ? 1 : -1
        }
    }
}
