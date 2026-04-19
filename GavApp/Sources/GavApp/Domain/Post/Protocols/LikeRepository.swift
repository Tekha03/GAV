import Foundation

public protocol LikeRepository {
    func add(like: Like) async throws
    func remove(like: Like) async throws
    func exists(like: Like) async throws -> Bool
}