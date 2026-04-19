import Foundation

public struct User: Equatable, Sendable, Codable {
    public let id: UUID
    public let email: String
    public let role: UserRole
    public let createdAt: Date
    public let updatedAt: Date

    public let lat: Double?
    public let lon: Double?
    public let locationStatus: LocationStatus
    public let locationVisibility: LocationVisibility
}

public enum UserRole: Int, Equatable, Sendable, CaseIterable, Codable {
    case user = 0
    case admin = 1

    public var stringValue: String {
        switch self {
        case .user:      return "User"
        case .admin:       return "Admin"
        }
    }
}