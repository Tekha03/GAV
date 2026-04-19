import Foundation
import SharedModels

protocol UserProfileServiceAPIProtocol {
    func getByUserID(userID: UUID) async throws -> UserProfileModel
    func update(userID: UUID, input: UpdateProfileInput) async throws
    func delete(userID: UUID) async throws
}

final class UserProfileServiceAPI: UserProfileServiceAPIProtocol {
    private let base: BaseAPI

    init(baseURL: URL, session: URLSession = .shared, authManager: AuthManager) {
        self.base = BaseAPI(baseURL: baseURL, session: session, authManager: authManager)
    }

    func getByUserID(userID: UUID) async throws -> UserProfileModel {
        let path = "/api/v1/users/\(userID.uuidString)/profile"
        let data = try await base.request(path)
        return try JSONDecoder().decode(UserProfileModel.self, from: data)
    }

    func update(userID: UUID, input: UpdateProfileInput) async throws {
        let path = "/api/v1/profiles/\(userID.uuidString)"
        let body = try JSONEncoder().encode(input.toModel())
        _ = try await base.request(path, method: "PUT", body: body)
    }

    func delete(userID: UUID) async throws {
        let path = "/api/v1/profiles/\(userID.uuidString)"
        _ = try await base.request(path, method: "DELETE")
    }
}