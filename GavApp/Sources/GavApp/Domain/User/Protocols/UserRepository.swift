import Foundation

public protocol UserRepository {
    func create(email: String, password: String) async throws -> User
    func getByID(id: UUID) async throws -> User
    func getByEmail(email: String) async throws -> User
    func update(id: UUID, input: UpdateUserInput) async throws
    func delete(id: UUID) async throws
    func updateUserLocation(
        id: UUID,
        lat: Double,
        lon: Double,
        status: LocationStatus,
        visibility: LocationVisibility
    ) async throws

    func findDogsNearby(
        id: UUID,
        centerLat: Double,
        centerLon: Double,
        radiusMeters: Double
    ) async throws -> [Dog]
}