public struct CommentMapper {
    public static func from(model: CommentModel) -> Comment {
        return Comment(
            id: model.id,
            postId: model.postId,
            userId: model.userId,
            content: model.content,
            createdAt: model.createdAt
        )
    }

    public static func to(model: Comment) -> CommentModel {
        return CommentModel(
            id: model.id,
            postId: model.postId,
            userId: model.userId,
            content: model.content,
            createdAt: model.createdAt
        )
    }
}
