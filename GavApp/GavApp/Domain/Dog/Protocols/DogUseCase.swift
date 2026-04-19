import Domain
import Foundation

public struct DogUseCase {
    private let repository: any DogRepository

    public init(repository: any DogRepository) {
        self.repository = repository
    }

    public func create(ownerID: UUID, input: CreateDogInput) async throws -> Dog {
        return try await repository.create(ownerID: ownerID, input: input)
    }

    public func update(ownerID: UUID, dogID: UUID, input: UpdateDogInput) async throws {
        try await repository.update(ownerID: ownerID, dogID: dogID, input: input)
    }

    public func delete(ownerID: UUID, dogID: UUID) async throws {
        try await repository.delete(ownerID: ownerID, dogID: dogID)
    }

    public func getPublic(dogID: UUID) async throws -> Dog {
        return try await repository.getPublic(dogID: dogID)
    }

    public func getPrivate(ownerID: UUID, dogID: UUID) async throws -> Dog {
        return try await repository.getPrivate(ownerID: ownerID, dogID: dogID)
    }

    public func listByOwnerID(ownerID: UUID) async throws -> [Dog] {
        return try await repository.listByOwnerID(ownerID: ownerID)
    }
}