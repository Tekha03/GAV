import Foundation

public struct Attachment: Identifiable, Codable {
    public let id: UUID
    public let messageID: UUID
    public let url: String
    public let type: AttachmentType
    public let fileName: String
    public let fileSize: Int64
}

public enum AttachmentType: String, Codable {
    case image
    case video
    case audio
    case document
}

public struct AttachmentInput: Codable {
    public let url: String
    public let type: AttachmentType
    public let fileName: String
    public let fileSize: Int64
}