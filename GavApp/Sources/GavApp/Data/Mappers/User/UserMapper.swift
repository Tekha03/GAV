public struct UserMapper {
    public static func from(model: UserModel) -> User {
        return User(
            id: model.id,
            email: model.email,
            role: model.role,
            createdAt: model.createdAt,
            updatedAt: model.updatedAt
        )
    }

    public static func to(model: User) -> UserModel {
        return UserModel(
            id: model.id,
            email: model.email,
            role: model.role,
            createdAt: model.createdAt,
            updatedAt: model.updatedAt
        )
    }
}