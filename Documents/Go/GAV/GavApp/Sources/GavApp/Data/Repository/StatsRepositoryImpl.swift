import Foundation

final class StatsRepositoryImpl: StatsRepository {
    private let api: any StatsServiceAPIProtocol
    private let mapper: UserStatsMapper
    private let profileMapper: ProfileStatsMapper
    private let postMapper: PostStatsMapper

    init(
        api: any StatsServiceAPIProtocol,
        mapper: UserStatsMapper = UserStatsMapper(),
        profileMapper: ProfileStatsMapper = ProfileStatsMapper(),
        postMapper: PostStatsMapper = PostStatsMapper()
    ) {
        self.api = api
        self.mapper = mapper
        self.profileMapper = profileMapper
        self.postMapper = postMapper
    }

    func userStats(userID: UUID) async throws -> UserStats {
        let model = try await api.userStats(userID: userID)
        return mapper.from(model: model)
    }

    func profileStats(userID: UUID) async throws -> ProfileStats {
        let model = try await api.profileStats(userID: userID)
        return ProfileStatsMapper.from(model: model)
    }

    func postStats(postID: UUID) async throws -> PostStats {
        let model = try await api.postStats(postID: postID)
        return PostStatsMapper.from(model: model)
    }

    func incrementPosts(userID: UUID) async throws { fatalError("API endpoint missing") }
    func incrementFollowers(userID: UUID) async throws { fatalError("API endpoint missing") }
    func incrementFollowings(userID: UUID) async throws { fatalError("API endpoint missing") }
    func incrementDogs(userID: UUID) async throws { fatalError("API endpoint missing") }
    func decrementPosts(userID: UUID) async throws { fatalError("API endpoint missing") }
    func decrementFollowers(userID: UUID) async throws { fatalError("API endpoint missing") }
    func decrementFollowings(userID: UUID) async throws { fatalError("API endpoint missing") }
    func decrementDogs(userID: UUID) async throws { fatalError("API endpoint missing") }
    func incrementPostLikes(postID: UUID) async throws { fatalError("API endpoint missing") }
    func incrementPostComments(postID: UUID) async throws { fatalError("API endpoint missing") }
    func decrementPostLikes(postID: UUID) async throws { fatalError("API endpoint missing") }
    func decrementPostComments(postID: UUID) async throws { fatalError("API endpoint missing") }
}