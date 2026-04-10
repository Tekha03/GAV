//
//  VaccinationRepository.swift
//  GavApp
//
//  Created by Виктория Кашуркина on 10.04.2026.
//


import Foundation

protocol VaccinationRepository {
    func create(_ vaccination: Vaccination) async throws
    func update(_ vaccination: Vaccination) async throws
    func delete(id: UUID) async throws
    func listByDogId(_ dogId: UUID) async throws -> [Vaccination]
    func getById(_ id: UUID) async throws -> Vaccination
}