public struct PostMapper {
    public func from(model: PostModel) -> Post {
        return Post(
            id: model.id,
            userId: model.userId,
            content: model.content,
            imageUrl: model.imageUrl,
            createdAt: model.createdAt
        )
    }

    public func to(model: Post) -> PostModel {
        return PostModel(
            id: model.id,
            userId: model.userId,
            content: model.content,
            imageUrl: model.imageUrl,
            createdAt: model.createdAt
        )
    }
}