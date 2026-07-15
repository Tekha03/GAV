import Foundation

public struct UserMapper {
    public static func from(model: UserModel) -> User {
        return User(
            id: model.id,
            email: model.email,
            role: role(from: model.role),
            createdAt: model.createdAt,
            updatedAt: model.updatedAt,
            lat: nil,
            lon: nil,
            locationStatus: .forcedOffline,
            locationVisibility: .noOne
        )
    }

    public static func to(model: User) -> UserModel {
        return UserModel(
            id: model.id,
            email: model.email,
            role: model.role.stringValue.lowercased(),
            createdAt: model.createdAt,
            updatedAt: model.updatedAt
        )
    }

    private static func role(from rawValue: String) -> UserRole {
        switch rawValue.lowercased() {
        case "admin":
            return .admin
        default:
            return .user
        }
    }
}
