import Foundation
import SharedModels

protocol UserServiceAPIProtocol {
    func getByID(id: UUID) async throws -> UserModel
    func update(id: UUID, input: UpdateUserInput) async throws
    func delete(id: UUID) async throws
}

final class UserServiceAPI: UserServiceAPIProtocol {
    private let base: BaseAPI

    init(baseURL: URL, session: URLSession = .shared, authManager: AuthManager) {
        self.base = BaseAPI(baseURL: baseURL, session: session, authManager: authManager)
    }

    func getByID(id: UUID) async throws -> UserModel {
        let path = "/api/v1/users/\(id.uuidString)"
        let data = try await base.request(path)
        return try JSONDecoder().decode(UserModel.self, from: data)
    }

    func update(id: UUID, input: UpdateUserInput) async throws {
        let path = "/api/v1/users/\(id.uuidString)"
        let body = try JSONEncoder().encode(input.toModel())
        _ = try await base.request(path, method: "PUT", body: body)
    }

    func delete(id: UUID) async throws {
        let path = "/api/v1/users/\(id.uuidString)"
        _ = try await base.request(path, method: "DELETE")
    }
}
