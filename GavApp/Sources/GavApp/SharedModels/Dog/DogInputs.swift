import Foundation

public struct CreateDogInput: Encodable, Sendable {
    public let name: String
    public let breed: String
    public let age: String
    public let status: String
    public let gender: String
    public let photoUrl: String
    public let notes: String

    public init(name: String, breed: String, age: String, status: String, gender: String, photoUrl: String, notes: String = "") {
        self.name = name
        self.breed = breed
        self.age = age
        self.status = status
        self.gender = gender
        self.photoUrl = photoUrl
        self.notes = notes
    }
}

public struct UpdateDogInput: Encodable {
    public let name: String?
    public let breed: String?
    public let age: String?
    public let status: String?
    public let gender: String?
    public let notes: String?

    public init(
        name: String? = nil,
        breed: String? = nil,
        age: String? = nil,
        status: String? = nil,
        gender: String? = nil,
        notes: String? = nil
    ) {
        self.name = name
        self.breed = breed
        self.age = age
        self.status = status
        self.gender = gender
        self.notes = notes
    }
}
