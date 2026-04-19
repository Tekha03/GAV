import Foundation

final class PostRepositoryImpl: PostRepository {
    private let api: any PostServiceAPIProtocol
    private let mapper: PostMapper

    init(api: any PostServiceAPIProtocol, mapper: PostMapper = PostMapper()) {
        self.api = api
        self.mapper = mapper
    }

    func create(userID: UUID, content: String, imageUrl: String?) async throws -> Post {
        let model = try await api.create(userID: userID, content: content, imageUrl: imageUrl)
        return mapper.from(model: model)
    }

    func get(id: UUID) async throws -> Post {
        let model = try await api.getByID(id: id)
        return mapper.from(model: model)
    }

    func listByUser(userID: UUID) async throws -> [Post] {
        let models = try await api.listByUser(userID: userID)
        return models.map { mapper.from(model: $0) }
    }

    func delete(userID: UUID, id: UUID) async throws {
        try await api.delete(id: id)
    }
}