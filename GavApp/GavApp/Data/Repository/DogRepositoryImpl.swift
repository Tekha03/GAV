import Foundation
import Domain
import Data
import SharedModels

final class DogRepositoryImpl: DogRepository {
    private let api: any DogServiceAPIProtocol
    private let mapper: DogMapper

    init(api: any DogServiceAPIProtocol, mapper: DogMapper = DogMapper()) {
        self.api = api
        self.mapper = mapper
    }

    func create(ownerID: UUID, input: CreateDogInput) async throws -> Dog {
        let model = try await api.create(ownerID: ownerID, input: input)
        return try DogMapper.from(model: model)
    }

    func update(ownerID: UUID, dogID: UUID, input: UpdateDogInput) async throws {
        try await api.update(dogID: dogID, input: input)
    }

    func delete(ownerID: UUID, dogID: UUID) async throws {
        try await api.delete(dogID: dogID)
    }

    func getPublic(dogID: UUID) async throws -> Dog {
        let model = try await api.getPublic(dogID: dogID)
        return try DogMapper.from(model: model)
    }

    func getPrivate(ownerID: UUID, dogID: UUID) async throws -> Dog {
        let model = try await api.getPrivate(dogID: dogID)
        return try DogMapper.from(model: model)
    }

    func listByOwnerID(ownerID: UUID) async throws -> [Dog] {
        let models = try await api.listByOwnerID(ownerID: ownerID)
        return models.compactMap { try? DogMapper.from(model: $0) }
    }
}