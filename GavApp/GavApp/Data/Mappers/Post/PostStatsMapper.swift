import Domain
import SharedModels

public struct PostStatsMapper {
    public static func from(model: PostStatsModel) -> Domain.PostStats {
        return Domain.PostStats(
            id: model.id,
            postId: model.postId,
            likesCount: model.likesCount,
            commentsCount: model.commentsCount,
            createdAt: model.createdAt,
            updatedAt: model.updatedAt
        )
    }

    public static func to(model: Domain.PostStats) -> PostStatsModel {
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