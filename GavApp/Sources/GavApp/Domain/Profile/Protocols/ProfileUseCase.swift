import Foundation

public struct ProfileUseCase {
    private let repository: any UserProfileRepository

    public init(repository: any UserProfileRepository) {
        self.repository = repository
    }

    public func create(userID: UUID, input: CreateProfileInput) async throws -> UserProfile {
        return try await repository.create(userID: userID, input: input)
    }

    public func getByUserID(userID: UUID) async throws -> UserProfile {
        return try await repository.getByUserID(userID: userID)
    }

    public func getStats(userID: UUID) async throws -> ProfileStats {
        return try await repository.getStats(userID: userID)
    }

    public func update(userID: UUID, input: UpdateProfileInput) async throws {
        try await repository.update(userID: userID, input: input)
    }

    public func delete(userID: UUID) async throws {
        try await repository.delete(userID: userID)
    }
}