import Foundation

final class VaccinationRepositoryImpl: VaccinationRepository {
    private let api: any VaccinationServiceAPIProtocol
    private let mapper: VaccinationMapper

    init(api: any VaccinationServiceAPIProtocol, mapper: VaccinationMapper = VaccinationMapper()) {
        self.api = api
        self.mapper = mapper
    }

    func create(dogID: UUID, input: CreateVaccinationInput) async throws -> Vaccination {
        let model = try await api.create(dogID: dogID, input: input)
        return VaccinationMapper.from(model: model)
    }

    func listByDogID(dogID: UUID) async throws -> [Vaccination] {
        let models = try await api.listByDogID(dogID: dogID)
        return models.map { VaccinationMapper.from(model: $0) }
    }

    func update(vaccinationID: UUID, dogID: UUID, input: UpdateVaccinationInput) async throws {
        try await api.update(vaccinationID: vaccinationID, input: input)
    }

    func delete(vaccinationID: UUID) async throws {
        try await api.delete(vaccinationID: vaccinationID)
    }
}