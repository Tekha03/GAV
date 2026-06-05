import Foundation

public protocol DogRepository {
    func create(ownerID: UUID, input: CreateDogInput) async throws -> Dog
    func update(ownerID: UUID, dogID: UUID, input: UpdateDogInput) async throws
    func delete(ownerID: UUID, dogID: UUID) async throws

    func getPublic(dogID: UUID) async throws -> Dog
    func getPrivate(ownerID: UUID, dogID: UUID) async throws -> Dog
    func listByOwnerID(ownerID: UUID) async throws -> [Dog]
}