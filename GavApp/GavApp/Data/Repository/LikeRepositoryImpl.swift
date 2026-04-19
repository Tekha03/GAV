import Foundation
import Domain
import Data
import SharedModels

final class LikeRepositoryImpl: LikeRepository {
    private let api: any PostServiceAPIProtocol
    private let mapper: LikeMapper

    init(api: any PostServiceAPIProtocol, mapper: LikeMapper = LikeMapper()) {
        self.api = api
        self.mapper = mapper
    }

    func add(like: Like) async throws {
        try await api.addLike(postID: like.postId)
    }

    func remove(like: Like) async throws {
        try await api.removeLike(postID: like.postId)
    }

    func exists(like: Like) async throws -> Bool {
        fatalError("API не проверяет существование like напрямую")
    }
}