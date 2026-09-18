import Foundation

protocol VaccinationServiceAPIProtocol: Sendable {
    func create(dogID: UUID, input: CreateVaccinationInput) async throws -> VaccinationModel
    func listByDogID(dogID: UUID) async throws -> [VaccinationModel]
    func update(vaccinationID: UUID, dogID: UUID, input: UpdateVaccinationInput) async throws
    func delete(vaccinationID: UUID) async throws
}

@available(macOS 12.0, *)
final class VaccinationServiceAPI: VaccinationServiceAPIProtocol, @unchecked Sendable {
    private let base: BaseAPI
    private let encoder: JSONEncoder
    private let decoder: JSONDecoder

    init(baseURL: URL, session: URLSession = .shared, authManager: AuthManager) {
        self.base = BaseAPI(baseURL: baseURL, session: session, authManager: authManager)
        self.encoder = JSONEncoder()
        self.encoder.dateEncodingStrategy = .iso8601
        self.decoder = JSONDecoder()
        self.decoder.dateDecodingStrategy = .custom { decoder in
            let value = try decoder.singleValueContainer().decode(String.self)
            let fractional = ISO8601DateFormatter()
            fractional.formatOptions = [.withInternetDateTime, .withFractionalSeconds]
            if let date = fractional.date(from: value) { return date }
            if let date = ISO8601DateFormatter().date(from: value) { return date }
            throw DecodingError.dataCorruptedError(
                in: try decoder.singleValueContainer(),
                debugDescription: "Invalid date: \(value)"
            )
        }
    }

    func create(dogID: UUID, input: CreateVaccinationInput) async throws -> VaccinationModel {
        let path = "/api/v1/dogs/\(dogID.uuidString)/vaccinations"
        let body = try encoder.encode(input)
        let data = try await base.request(path, method: "POST", body: body)
        return try decoder.decode(VaccinationModel.self, from: data)
    }

    func listByDogID(dogID: UUID) async throws -> [VaccinationModel] {
        let path = "/api/v1/dogs/\(dogID.uuidString)/vaccinations"
        let data = try await base.request(path)
        return try decoder.decode([VaccinationModel].self, from: data)
    }

    func update(vaccinationID: UUID, dogID: UUID, input: UpdateVaccinationInput) async throws {
        let path = "/api/v1/dogs/\(dogID.uuidString)/vaccinations/\(vaccinationID.uuidString)"
        let body = try encoder.encode(input)
        _ = try await base.request(path, method: "PUT", body: body)
    }

    func delete(vaccinationID: UUID) async throws {
        let path = "/api/v1/vaccinations/\(vaccinationID.uuidString)"
        _ = try await base.request(path, method: "DELETE")
    }
}
