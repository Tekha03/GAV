//
//  VaccinationDTO.swift
//  GavApp
//
//  Created by Виктория Кашуркина on 10.04.2026.
//
import Foundation

struct VaccinationDTO: Codable {
    let id: UUID
    let dogId: UUID
    let name: String
    let doneAt: Date
    let nextDueAt: Date?
    let notes: String?

    enum CodingKeys: String, CodingKey {
        case id
        case dogId = "dog_id"
        case name
        case doneAt = "done_at"
        case nextDueAt = "next_due_at"
        case notes
    }
}

extension VaccinationDTO {
    func toDomain() -> Vaccination {
        .init(
            id: id,
            dogId: dogId,
            name: name,
            doneAt: doneAt,
            nextDueAt: nextDueAt,
            notes: notes ?? ""
        )
    }
}
