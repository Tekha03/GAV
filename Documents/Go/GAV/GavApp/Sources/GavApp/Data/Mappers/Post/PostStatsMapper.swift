public struct PostStatsMapper {
    public static func from(model: PostStatsModel) -> PostStats {
        return PostStats(
            id: model.id,
            postId: model.postId,
            likesCount: model.likesCount,
            commentsCount: model.commentsCount,
            createdAt: model.createdAt,
            updatedAt: model.updatedAt
        )
    }

    public static func to(model: PostStats) -> PostStatsModel {
        return PostStatsModel(
            id: model.id,
            postId: model.postId,
            likesCount: model.likesCount,
            commentsCount: model.commentsCount,
            createdAt: model.createdAt,
            updatedAt: model.updatedAt
        )
    }
}