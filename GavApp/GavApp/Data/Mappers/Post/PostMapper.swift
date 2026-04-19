// Data/Mappers/Post/PostMapper.swift
import Domain
import SharedModels

public struct PostMapper {
    public static func from(model: PostModel) -> Domain.Post {
        return Domain.Post(
            id: model.id,
            userId: model.userId,
            content: model.content,
            imageUrl: model.imageUrl,
            createdAt: model.createdAt
        )
    }

    public static func to(model: Domain.Post) -> PostModel {
        return PostModel(
            id: model.id,
            userId: model.userId,
            content: model.content,
            imageUrl: model.imageUrl,
            createdAt: model.createdAt
        )
    }
}