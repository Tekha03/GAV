import Foundation

final class UserProfileRepositoryImpl: UserProfileRepository {
    private let api: any UserProfileServiceAPIProtocol
    private let mapper: UserProfileMapper

    init(
        api: any UserProfileServiceAPIProtocol,
        mapper: UserProfileMapper = UserProfileMapper()
    ) {
        self.api = api
        self.mapper = mapper
    }

    func create(userID: UUID, input: CreateProfileInput) async throws -> UserProfile {
        let requestModel = input.toModel(userID: userID)
        let responseModel = try await api.create(userID: userID, model: requestModel)
        return try mapper.from(model: responseModel)
    }

    func getByUserID(userID: UUID) async throws -> UserProfile {
        let model = try await api.getByUserID(userID: userID)
        return try mapper.from(model: model)
    }

    func getStats(userID: UUID) async throws -> ProfileStats {
        fatalError("Inject StatsRepository сюда отдельно")
    }

    func update(userID: UUID, input: UpdateProfileInput) async throws {
        try await api.update(userID: userID, input: input)
    }

    func delete(userID: UUID) async throws {
        try await api.delete(userID: userID)
    }
}