import Foundation

public struct Vaccination: Identifiable, Equatable, Sendable {
    public let id: UUID
    public let dogId: UUID
    public let name: String
    public let doneAt: Date
    public let nextDueAt: Date?
    public let notes: String
    public let status: VaccinationStatus

    public init(
        id: UUID,
        dogId: UUID,
        name: String,
        doneAt: Date,
        nextDueAt: Date?,
        notes: String,
        status: VaccinationStatus
    ) {
        self.id = id
        self.dogId = dogId
        self.name = name
        self.doneAt = doneAt
        self.nextDueAt = nextDueAt
        self.notes = notes
        self.status = status
    }
}