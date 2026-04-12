import Foundation

public struct Dog: Identifiable, Equatable, Sendable {
    public let id: UUID
    public let ownerId: UUID

    public let name: String
    public let breed: String
    public let photoURL: String

    public let status: DogStatus
    public let age: DogAge
    public let gender: DogGender

    public init(
        id: UUID,
        ownerId: UUID,
        name: String,
        breed: String,
        photoURL: String,
        status: DogStatus,
        age: DogAge,
        gender: DogGender,
    ) {
        self.id = id
        self.ownerId = ownerId
        self.name = name
        self.breed = breed
        self.photoURL = photoURL
        self.status = status
        self.age = age
        self.gender = DogGender(rawValue: gender.rawValue) ?? .unknown
    }
}