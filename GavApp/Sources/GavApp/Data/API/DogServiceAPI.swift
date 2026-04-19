import Foundation

protocol DogServiceAPIProtocol {
    func create(ownerID: UUID, input: CreateDogInput) async throws -> DogModel
    func getPrivate(dogID: UUID) async throws -> DogModel
    func getPublic(dogID: UUID) async throws -> DogModel
    func update(dogID: UUID, input: UpdateDogInput) async throws
    func delete(dogID: UUID) async throws
    func listByOwnerID(ownerID: UUID) async throws -> [DogModel]
}

@available(macOS 12.0, *)
final class DogServiceAPI: DogServiceAPIProtocol {
    private let base: BaseAPI

    init(baseURL: URL, session: URLSession = .shared, authManager: AuthManager) {
        self.base = BaseAPI(baseURL: baseURL, session: session, authManager: authManager)
    }

    func create(ownerID: UUID, input: CreateDogInput) async throws -> DogModel {
        let path = "/api/v1/dogs"
        let model = CreateDogModel(
            ownerId: ownerID,
            name: input.name,
            breed: input.breed,
            status: input.status,
            age: input.age,
            gender: input.gender
        )
        let body = try JSONEncoder().encode(model)
        let data = try await base.request(path, method: "POST", body: body)
        return try JSONDecoder().decode(DogModel.self, from: data)
    }

    func getPrivate(dogID: UUID) async throws -> DogModel {
        let path = "/api/v1/dogs/\(dogID.uuidString)"
        let data = try await base.request(path)
        return try JSONDecoder().decode(DogModel.self, from: data)
    }

    func getPublic(dogID: UUID) async throws -> DogModel {
        let path = "/api/v1/dogs/\(dogID.uuidString)"
        let data = try await base.request(path)
        return try JSONDecoder().decode(DogModel.self, from: data)
    }

    func update(dogID: UUID, input: UpdateDogInput) async throws {
        let path = "/api/v1/dogs/\(dogID.uuidString)"
        let model = UpdateDogModel(
            name: input.name,
            breed: input.breed,
            status: input.status,
            age: input.age,
            gender: input.gender
        )
        let body = try JSONEncoder().encode(model)
        _ = try await base.request(path, method: "PUT", body: body)
    }

    func delete(dogID: UUID) async throws {
        let path = "/api/v1/dogs/\(dogID.uuidString)"
        _ = try await base.request(path, method: "DELETE")
    }

    func listByOwnerID(ownerID: UUID) async throws -> [DogModel] {
        let path = "/api/v1/dogs/owner/\(ownerID.uuidString)"
        let data = try await base.request(path)
        return try JSONDecoder().decode([DogModel].self, from: data)
    }
}
