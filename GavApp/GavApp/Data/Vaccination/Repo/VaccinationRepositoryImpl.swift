//
//  VaccinationRepositoryImpl.swift
//  GavApp
//
//  Created by Виктория Кашуркина on 10.04.2026.
//
import Foundation

final class VaccinationRepositoryImpl: VaccinationRepository {

    private let api: VaccinationServiceAPI

    init(api: VaccinationServiceAPI) {
        self.api = api
    }

    func create(_ vaccination: Vaccination) async throws {
        try await api.create(vaccination.toDTO())
    }

    func update(_ vaccination: Vaccination) async throws {
        try await api.update(vaccination.toDTO())
    }

    func delete(id: UUID) async throws {
        try await api.delete(id: id)
    }

    func listByDogId(_ dogId: UUID) async throws -> [Vaccination] {
        let dtos = try await api.listByDogId(dogId)
        return dtos.map { $0.toDomain() }
    }

    func getById(_ id: UUID) async throws -> Vaccination {
        let dto = try await api.getById(id)
        return dto.toDomain()
    }
}
