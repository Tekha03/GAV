import Foundation

protocol UserServiceAPIProtocol: Sendable {
    func getByID(id: UUID) async throws -> UserModel
    func update(id: UUID, input: UpdateUserInput) async throws
    func delete(id: UUID) async throws
    func getByEmail(email: String) async throws -> UserModel
    func updateLocation(id: UUID, input: UpdateLocationInput) async throws
    func findDogsNearby(
        id: UUID,
        centerLat: Double,
        centerLon: Double,
        radiusMeters: Double
    ) async throws -> [DogModel]
}

@available(macOS 12.0, *)
final class UserServiceAPI: UserServiceAPIProtocol, @unchecked Sendable {
    private let base: BaseAPI

    init(baseURL: URL, session: URLSession = .shared, authManager: AuthManager) {
        self.base = BaseAPI(baseURL: baseURL, session: session, authManager: authManager)
    }

    func getByID(id: UUID) async throws -> UserModel {
        let path = "/api/v1/users/\(id.uuidString)"
        let data = try await base.request(path)
        return try JSONDecoder().decode(UserModel.self, from: data)
    }

    func getByEmail(email: String) async throws -> UserModel {
        let path = "/api/v1/users/email/\(email)"
        let data = try await base.request(path)
        return try JSONDecoder().decode(UserModel.self, from: data)
    }

    func update(id: UUID, input: UpdateUserInput) async throws {
        let path = "/api/v1/users/\(id.uuidString)"
        let body = try JSONEncoder().encode(input)
        _ = try await base.request(path, method: "PUT", body: body)
    }

    func delete(id: UUID) async throws {
        let path = "/api/v1/users/\(id.uuidString)"
        _ = try await base.request(path, method: "DELETE")
    }

    func updateLocation(id: UUID, input: UpdateLocationInput) async throws {
        let path = "/api/v1/users/\(id.uuidString)/location"
        let body = try JSONEncoder().encode(input)
        _ = try await base.request(path, method: "PUT", body: body)
    }

    func findDogsNearby(
        id: UUID,
        centerLat: Double,
        centerLon: Double,
        radiusMeters: Double
    ) async throws -> [DogModel] {

        let path = "/api/v1/dogs/nearby?lat=\(centerLat)&lon=\(centerLon)&radius=\(radiusMeters)"

        let data = try await base.request(path)

        return try JSONDecoder().decode([DogModel].self, from: data)
    }
}
