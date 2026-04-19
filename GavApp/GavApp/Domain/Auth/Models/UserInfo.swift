import Foundation

public struct UserInfo: Identifiable, Equatable, Sendable {
    public let id: UUID
    public let email: String

    public init(id: UUID, email: String) {
        self.id = id
        self.email = email
    }
}