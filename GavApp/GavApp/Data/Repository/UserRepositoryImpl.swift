// Data/Repositories/UserRepositoryImpl.swift
import Foundation
import Domain
import Data
import SharedModels

final class UserRepositoryImpl: UserRepository {
    private let api: any UserServiceAPIProtocol
    private let mapper: UserMapper

    init(api: any UserServiceAPIProtocol, mapper: UserMapper = UserMapper()) {
        self.api = api
        self.mapper = mapper
    }

    func create(email: String, password: String) async throws -> User {
        // Если на бэкенде нет прямого /api/v1/auth/register → User, можно в AuthRepository делать register → User
        fatalError("API не возвращает User напрямую при registration, нужно уточнить")
    }

    func getByID(id: UUID) async throws -> User {
        let model = try await api.getByID(id: id)
        return UserMapper.from(model: model)
    }

    func update(id: UUID, input: UpdateUserInput) async throws {
        try await api.update(id: id, input: input)
    }

    func delete(id: UUID) async throws {
        try await api.delete(id: id)
    }
}