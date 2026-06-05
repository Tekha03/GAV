import Foundation

public struct CreateDogInput: Encodable, Sendable {
    public let name: String
    public let breed: String
    public let age: Int
    public let status: String
    public let gender: String
    

    public init(name: String, breed: String, age: Int, status: String, gender: String) {
        self.name = name
        self.breed = breed
        self.age = age
        self.status = status
        self.gender = gender
    }
}

public struct UpdateDogInput: Encodable {
    public let name: String?
    public let breed: String?
    public let age: Int?
    public let status: String?
    public let gender: String?

    public init(
        name: String? = nil,
        breed: String? = nil,
        age: Int? = nil,
        status: String? = nil,
        gender: String? = nil
    ) {
        self.name = name
        self.breed = breed
        self.age = age
        self.status = status
        self.gender = gender
    }
}