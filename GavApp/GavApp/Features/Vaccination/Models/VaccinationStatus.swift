//
//  VaccinationStatus.swift
//  GavApp
//
//  Created by Виктория Кашуркина on 10.04.2026.
//


import Foundation

enum VaccinationStatus {
    case actual
    case dueSoon
    case overdue

    static func resolve(doneAt: Date, nextDueAt: Date?) -> VaccinationStatus {
        guard let next = nextDueAt else {
            return .actual
        }

        let now = Date()

        if next < now {
            return .overdue
        }

        if let soon = Calendar.current.date(byAdding: .day, value: 30, to: now),
           next < soon {
            return .dueSoon
        }

        return .actual
    }
}