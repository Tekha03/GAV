import Foundation

public struct DogModel: Codable, Equatable {
    public let id: UUID
    public let ownerId: UUID

    public let name: String
    public let breed: String
    public let photoUrl: String

    public let status: String
    public let age: String
    public let gender: String
}

extension DogModel {
    private enum CodingKeys: String, CodingKey {
        case id
        case ownerId = "owner_id"
        case name
        case breed
        case photoUrl = "photo_url"
        case status
        case age
        case gender
    }
}