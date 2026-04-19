import Domain
import Foundation

public struct VaccinationUseCase {
    private let repository: any VaccinationRepository

    public init(repository: any VaccinationRepository) {
        self.repository = repository
    }

    public func create(dogID: UUID, input: CreateVaccinationInput) async throws -> Vaccination {
        return try await repository.create(dogID: dogID, input: input)
    }

    public func listByDogID(dogID: UUID) async throws -> [Vaccination] {
        return try await repository.listByDogID(dogID: dogID)
    }

    public func update(vaccinationID: UUID, dogID: UUID, input: UpdateVaccinationInput) async throws {
        try await repository.update(
            vaccinationID: vaccinationID,
            dogID: dogID,
            input: input
        )
    }

    public func delete(vaccinationID: UUID) async throws {
        try await repository.delete(vaccinationID: vaccinationID)
    }
}