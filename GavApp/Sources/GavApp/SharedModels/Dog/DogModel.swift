import Foundation

public struct DogModel: Codable, Equatable {
    public let id: UUID
    public let ownerId: UUID

    public let name: String
    public let breed: String
    public let photoUrl: String
    public let notes: String

    public let status: String
    public let age: String
    public let gender: String
    public let lat: Double?
    public let lon: Double?

    public init(
        id: UUID,
        ownerId: UUID,
        name: String,
        breed: String,
        photoUrl: String,
        notes: String = "",
        status: String,
        age: String,
        gender: String,
        lat: Double? = nil,
        lon: Double? = nil
    ) {
        self.id = id
        self.ownerId = ownerId
        self.name = name
        self.breed = breed
        self.photoUrl = photoUrl
        self.notes = notes
        self.status = status
        self.age = age
        self.gender = gender
        self.lat = lat
        self.lon = lon
    }
}

extension DogModel {
    private enum CodingKeys: String, CodingKey {
        case id
        case ownerId = "owner_id"
        case name
        case breed
        case photoUrl = "photo_url"
        case notes
        case status
        case age
        case gender
        case lat
        case lon
    }
}
