import Domain

public struct UserUseCase {
    private let repository: any UserRepository

    public init(repository: any UserRepository) {
        self.repository = repository
    }

    public func create(email: String, password: String) async throws -> User {
        return try await repository.create(email: email, password: password)
    }

    public func getByID(id: UUID) async throws -> User {
        return try await repository.getByID(id: id)
    }

    public func getByEmail(email: String) async throws -> User {
        return try await repository.getByEmail(email: email)
    }

    public func update(id: UUID, input: UpdateUserInput) async throws {
        try await repository.update(id: id, input: input)
    }

    public func delete(id: UUID) async throws {
        try await repository.delete(id: id)
    }
}