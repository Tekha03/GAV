//
//  VaccinationServiceAPI.swift
//  GavApp
//
//  Created by Виктория Кашуркина on 10.04.2026.
//
import Foundation

protocol VaccinationServiceAPI {
    func create(_ dto: VaccinationDTO) async throws
    func update(_ dto: VaccinationDTO) async throws
    func delete(id: UUID) async throws
    func listByDogId(_ dogId: UUID) async throws -> [VaccinationDTO]
    func getById(_ id: UUID) async throws -> VaccinationDTO
}
