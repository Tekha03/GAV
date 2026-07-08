import Foundation

struct CreateDogModel: Encodable {
    let ownerId: UUID
    let name: String
    let breed: String
    let status: String
    let age: String
    let gender: String
    let photoUrl: String
    let notes: String

    private enum CodingKeys: String, CodingKey {
        case ownerId = "owner_id"
        case name
        case breed
        case status
        case age
        case gender
        case photoUrl = "photo_url"
        case notes
    }
}

struct UpdateDogModel: Encodable {
    let name: String?
    let breed: String?
    let status: String?
    let age: String?
    let gender: String?
    let notes: String?
}
