import Domain
import SharedModels

public struct CommentMapper {
    public static func from(model: CommentModel) -> Domain.Comment {
        return Domain.Comment(
            id: model.id,
            postId: model.postId,
            userId: model.userId,
            content: model.content,
            createdAt: model.createdAt
        )
    }

    public static func to(model: Domain.Comment) -> CommentModel {
        return CommentModel(
            id: model.id,
            postId: model.postId,
            userId: model.userId,
            content: model.content,
            createdAt: model.createdAt
        )
    }
}