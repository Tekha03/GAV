import Foundation

public struct UserModel: Codable, Equatable, Sendable {
    public let id: UUID
    public let email: String
    public let role: String
    public let createdAt: Date
    public let updatedAt: Date
}

extension UserModel {
    private enum CodingKeys: String, CodingKey {
        case id
        case email
        case role
        case createdAt = "created_at"
        case updatedAt = "updated_at"
    }
}