public struct LikeMapper {
    public static func from(model: LikeModel) -> Like {
        return Like(userId: model.userId, postId: model.postId)
    }

    public static func to(model: Like) -> LikeModel {
        return LikeModel(userId: model.userId, postId: model.postId)
    }
}