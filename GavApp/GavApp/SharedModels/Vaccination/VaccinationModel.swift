import Foundation

public struct VaccinationModel: Codable, Equatable {
    public let id: UUID
    public let dogId: UUID

    public let name: String
    public let doneAt: Date
    public let nextDueAt: Date?
    public let notes: String?
}

extension VaccinationModel {
    private enum CodingKeys: String, CodingKey {
        case id
        case dogId = "dog_id"
        case name
        case doneAt = "done_at"
        case nextDueAt = "next_due_at"
        case notes
    }
}