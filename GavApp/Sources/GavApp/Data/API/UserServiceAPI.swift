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
        guard let latitude = input.lat, let longitude = input.lon else {
            throw APIError.server(
                statusCode: 422,
                body: APIErrorBody(
                    code: .validation,
                    category: .validation,
                    message: "Location coordinates are required",
                    details: nil
                )
            )
        }

        let updateBody = try JSONEncoder().encode(
            WalkLocationRequest(latitude: latitude, longitude: longitude)
        )

        do {
            _ = try await base.request(
                "/api/v1/walks/current/location",
                method: "PATCH",
                body: updateBody
            )
        } catch let error as APIError where error.category == .notFound {
            let startBody = try JSONEncoder().encode(
                StartWalkRequest(
                    latitude: latitude,
                    longitude: longitude,
                    visibility: input.visibility.apiValue
                )
            )

            do {
                _ = try await base.request(
                    "/api/v1/walks/start",
                    method: "POST",
                    body: startBody
                )
            } catch let startError as APIError where startError.category == .conflict {
                _ = try await base.request(
                    "/api/v1/walks/current/location",
                    method: "PATCH",
                    body: updateBody
                )
            }
        }
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

private struct WalkLocationRequest: Encodable {
    let latitude: Double
    let longitude: Double
}

private struct StartWalkRequest: Encodable {
    let latitude: Double
    let longitude: Double
    let visibility: Int
}
