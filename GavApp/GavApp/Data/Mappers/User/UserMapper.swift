import Domain
import SharedModels

public struct UserMapper {
    public static func from(model: UserModel) -> Domain.User {
        return Domain.User(
            id: model.id,
            email: model.email,
            role: model.role,
            createdAt: model.createdAt,
            updatedAt: model.updatedAt
        )
    }

    public static func to(model: Domain.User) -> UserModel {
        return UserModel(
            id: model.id,
            email: model.email,
            role: model.role,
            createdAt: model.createdAt,
            updatedAt: model.updatedAt
        )
    }
}