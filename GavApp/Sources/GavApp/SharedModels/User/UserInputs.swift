import Foundation

public struct UpdateUserInput: Encodable {
    public let email: String?
    public let role: String?

    public init(email: String? = nil, role: String? = nil) {
        self.email = email
        self.role = role
    }
}