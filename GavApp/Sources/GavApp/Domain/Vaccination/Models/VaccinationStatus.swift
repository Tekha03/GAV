import Foundation

public enum VaccinationStatus: Equatable, Sendable {
    case upToDate
    case outdated
    case upcoming
    case expired

    public static func resolve(
        doneAt: Date,
        nextDueAt: Date?
    ) -> VaccinationStatus {
        let now = Date()
        guard let next = nextDueAt else {
            return .expired
        }
        if now < doneAt { return .upcoming }
        if next > now { return .upToDate }
        return .expired
    }
}
