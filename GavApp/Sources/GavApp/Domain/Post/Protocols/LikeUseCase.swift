import Foundation

public struct LikeUseCase {
    private let repository: any LikeRepository

    public init(repository: any LikeRepository) {
        self.repository = repository
    }

    public func add(like: Like) async throws {
        try await repository.add(like: like)
    }

    public func remove(like: Like) async throws {
        try await repository.remove(like: like)
    }

    public func exists(like: Like) async throws -> Bool {
        return try await repository.exists(like: like)
    }
}