//
//  Vaccination.swift
//  GavApp
//
//  Created by Виктория Кашуркина on 09.04.2026.
//
import Foundation

struct Vaccination: Identifiable {
    let id: UUID
    let dogId: UUID

    let name: String
    let doneAt: Date
    let nextDueAt: Date?
    let notes: String

    var status: VaccinationStatus {
        VaccinationStatus.resolve(doneAt: doneAt, nextDueAt: nextDueAt)
    }
}

extension Vaccination {
    func toDTO() -> VaccinationDTO {
        .init(
            id: id,
            dogId: dogId,
            name: name,
            doneAt: doneAt,
            nextDueAt: nextDueAt,
            notes: notes.isEmpty ? nil : notes
        )
    }
}
