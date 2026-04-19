import Foundation

public struct UserUseCase {
    private let repository: UserRepository

    public func updateUserLocation(
        id: UUID,
        lat: Double,
        lon: Double,
        status: LocationStatus,
        visibility: LocationVisibility
    ) async throws {
        return try await repository.updateUserLocation(id: id, lat: lat, lon: lon, status: status, visibility: visibility)
    }

    public func findDogsNearby(
        userId: UUID,
        centerLat: Double,
        centerLon: Double,
        radiusMeters: Double
    ) async throws -> [Dog] {
        return try await repository.findDogsNearby(id: userId, centerLat: centerLat, centerLon: centerLon, radiusMeters: radiusMeters)
    }
}