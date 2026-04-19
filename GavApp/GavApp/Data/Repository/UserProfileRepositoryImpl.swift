import Foundation
import Domain
import Data
import SharedModels

final class UserProfileRepositoryImpl: UserProfileRepository {
    private let api: any UserProfileServiceAPIProtocol
    private let mapper: UserProfileMapper

    init(api: any UserProfileServiceAPIProtocol, mapper: UserProfileMapper = UserProfileMapper()) {
        self.api = api
        self.mapper = mapper
    }

    func create(userID: UUID, input: CreateProfileInput) async throws -> UserProfile {
        let model = input.toModel()
        _ = try await api.create(userID: userID, model: model)
        fatalError("API не возвращает модель после create, нужно уточнить эндпоинт и модель")
    }

    func getByUserID(userID: UUID) async throws -> UserProfile {
        let model = try await api.getByUserID(userID: userID)
        return try UserProfileMapper.from(model: model)
    }

    func update(userID: UUID, input: UpdateProfileInput) async throws {
        try await api.update(userID: userID, input: input)
    }

    func delete(userID: UUID) async throws {
        try await api.delete(userID: userID)
    }
}