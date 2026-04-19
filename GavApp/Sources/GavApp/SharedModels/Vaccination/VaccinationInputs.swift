import Foundation

public struct CreateVaccinationInput: Encodable {
    public let name: String
    public let doneAt: Date
    public let nextDueAt: Date?
    public let notes: String?

    public init(
        name: String,
        doneAt: Date,
        nextDueAt: Date? = nil,
        notes: String? = nil
    ) {
        self.name = name
        self.doneAt = doneAt
        self.nextDueAt = nextDueAt
        self.notes = notes
    }
}

public struct UpdateVaccinationInput: Encodable {
    public let name: String?
    public let doneAt: Date?
    public let nextDueAt: Date?
    public let notes: String?

    public init(
        name: String? = nil,
        doneAt: Date? = nil,
        nextDueAt: Date? = nil,
        notes: String? = nil
    ) {
        self.name = name
        self.doneAt = doneAt
        self.nextDueAt = nextDueAt
        self.notes = notes
    }
}