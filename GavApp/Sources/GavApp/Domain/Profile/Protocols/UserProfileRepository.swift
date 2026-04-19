import Foundation

public protocol UserProfileRepository {
    func create(userID: UUID, input: CreateProfileInput) async throws -> UserProfile
    func getByUserID(userID: UUID) async throws -> UserProfile
    func getStats(userID: UUID) async throws -> ProfileStats
    func update(userID: UUID, input: UpdateProfileInput) async throws
    func delete(userID: UUID) async throws
}