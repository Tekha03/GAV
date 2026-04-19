import Foundation

public struct MediaInfoModel: Codable, Equatable {
    public let url: String
    public let mimeType: String?
}

extension MediaInfoModel {
    private enum CodingKeys: String, CodingKey {
        case url
        case mimeType = "mime_type"
    }
}