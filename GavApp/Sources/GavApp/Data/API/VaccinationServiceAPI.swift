import Foundation

protocol VaccinationServiceAPIProtocol {
    func create(dogID: UUID, input: CreateVaccinationInput) async throws -> VaccinationModel
    func listByDogID(dogID: UUID) async throws -> [VaccinationModel]
    func update(vaccinationID: UUID, input: UpdateVaccinationInput) async throws
    func delete(vaccinationID: UUID) async throws
}

@available(macOS 12.0, *)
final class VaccinationServiceAPI: VaccinationServiceAPIProtocol {
    private let base: BaseAPI

    init(baseURL: URL, session: URLSession = .shared, authManager: AuthManager) {
        self.base = BaseAPI(baseURL: baseURL, session: session, authManager: authManager)
    }

    // POST /api/v1/dogs/{dogId}/vaccinations
    func create(dogID: UUID, input: CreateVaccinationInput) async throws -> VaccinationModel {
        let path = "/api/v1/dogs/\(dogID.uuidString)/vaccinations"
        let body = try JSONEncoder().encode(input)
        let data = try await base.request(path, method: "POST", body: body)
        return try JSONDecoder().decode(VaccinationModel.self, from: data)
    }

    // GET /api/v1/dogs/{dogId}/vaccinations
    func listByDogID(dogID: UUID) async throws -> [VaccinationModel] {
        let path = "/api/v1/dogs/\(dogID.uuidString)/vaccinations"
        let data = try await base.request(path)
        return try JSONDecoder().decode([VaccinationModel].self, from: data)
    }

    // PUT /api/v1/dogs/{dogId}/vaccinations/{vaccinationID}
    func update(vaccinationID: UUID, input: UpdateVaccinationInput) async throws {
        // TODO: передать в Продукте `dogID` — или сделать другой путь
        // Пока включи в `input` или передай отдельно.
        let path = "/api/v1/dogs/...TODO.../vaccinations/\(vaccinationID.uuidString)"
        let body = try JSONEncoder().encode(input)
        _ = try await base.request(path, method: "PUT", body: body)
    }

    // DELETE /api/v1/vaccinations/{id}
    func delete(vaccinationID: UUID) async throws {
        let path = "/api/v1/vaccinations/\(vaccinationID.uuidString)"
        _ = try await base.request(path, method: "DELETE")
    }
}