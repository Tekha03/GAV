import Domain
import SharedModels

public struct LikeMapper {
    public static func from(model: LikeModel) -> Domain.Like {
        return Domain.Like(userId: model.userId, postId: model.postId)
    }

    public static func to(model: Domain.Like) -> LikeModel {
        return LikeModel(userId: model.userId, postId: model.postId)
    }
}