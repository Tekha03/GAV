import Foundation

final class UserRepositoryImpl: UserRepository {
    private let api: any UserServiceAPIProtocol
    private let mapper: UserMapper

    init(api: any UserServiceAPIProtocol, mapper: UserMapper = UserMapper()) {
        self.api = api
        self.mapper = mapper
    }

    func create(email: String, password: String) async throws -> User {
        // здесь registration у тебя делает AuthService, а User не возвращается → нужно уточнить API
        fatalError("API не возвращает User напрямую при registration, нужно уточнить")
    }

    func getByID(id: UUID) async throws -> User {
        let model = try await api.getByID(id: id)
        return UserMapper.from(model: model)
    }

    func getByEmail(email: String) async throws -> User {
        let model = try await api.getByEmail(email: email)
        return UserMapper.from(model: model)
    }

    func update(id: UUID, input: UpdateUserInput) async throws {
        try await api.update(id: id, input: input)
    }

    func delete(id: UUID) async throws {
        try await api.delete(id: id)
    }

    func updateUserLocation(
        id: UUID,
        lat: Double,
        lon: Double,
        status: LocationStatus,
        visibility: LocationVisibility
    ) async throws {
        let input = UpdateLocationInput(
            lat: lat,
            lon: lon,
            locationStatus: status,
            visibility: visibility
        )

        try await api.updateLocation(id: id, input: input)
    }

    func findDogsNearby(
        id: UUID,
        centerLat: Double,
        centerLon: Double,
        radiusMeters: Double
    ) async throws -> [Dog] {
        try await api.findDogsNearby(
            id: id,
            centerLat: centerLat,
            centerLon: centerLon,
            radiusMeters: radiusMeters
        ).map(DogMapper.from)
    }
}