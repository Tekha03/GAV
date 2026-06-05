import Foundation

struct CreateDogModel: Encodable {
    let ownerId: UUID
    let name: String
    let breed: String
    let status: String
    let age: Int
    let gender: String
}

struct UpdateDogModel: Encodable {
    let name: String?
    let breed: String?
    let status: String?
    let age: Int?
    let gender: String?
}