import Foundation
import SharedModels

protocol VaccinationServiceAPIProtocol {
    func create(dogID: UUID, input: CreateVaccinationInput) async throws -> VaccinationModel
    func listByDogID(dogID: UUID) async throws -> [VaccinationModel]
    func update(vaccinationID: UUID, input: UpdateVaccinationInput) async throws
    func delete(vaccinationID: UUID) async throws
}

final class VaccinationServiceAPI: VaccinationServiceAPIProtocol {
    private let base: BaseAPI

    init(baseURL: URL, session: URLSession = .shared, authManager: AuthManager) {
        self.base = BaseAPI(baseURL: baseURL, session: session, authManager: authManager)
    }

    func create(dogID: UUID, input: CreateVaccinationInput) async throws -> VaccinationModel {
        let path = "/api/v1/dogs/\(dogID.uuidString)/vaccinations"
        let model = CreateVaccinationModel(
            name: input.name,
            doneAt: input.doneAt,
            nextDueAt: input.nextDueAt,
            notes: input.notes
        )