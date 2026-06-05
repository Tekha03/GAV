import Foundation

public enum MemberRole: String, Codable {
    case owner
    case admin
    case participant
}

public struct ChatMember: Codable {
    public let userId: UUID
    public let chatId: UUID
    public let role: MemberRole
    public let muted: Bool
    public let lastReadMessageId: UUID?
}