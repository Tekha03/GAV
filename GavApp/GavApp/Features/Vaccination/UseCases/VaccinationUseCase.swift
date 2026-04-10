//
//  VaccinationUseCase.swift
//  GavApp
//
//  Created by Виктория Кашуркина on 10.04.2026.
//


import Foundation

protocol VaccinationUseCase {
    func load(dogId: UUID) async throws -> [Vaccination]
    func create(dogId: UUID, input: CreateVaccinationInput) async throws
    func delete(vaccinationId: UUID) async throws
}