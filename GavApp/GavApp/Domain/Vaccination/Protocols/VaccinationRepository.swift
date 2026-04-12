import Foundation

public protocol VaccinationRepository {
    func create(dogID: UUID, input: CreateVaccinationInput) async throws -> Vaccination
    func listByDogID(dogID: UUID) async throws -> [Vaccination]
    func update(vaccinationID: UUID, dogID: UUID, input: UpdateVaccinationInput) async throws
    func delete(vaccinationID: UUID) async throws
}