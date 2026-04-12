import Domain
import Foundation
import SharedModels

public struct VaccinationMapper {
    public static func from(model: VaccinationModel) -> Domain.Vaccination {
        let status = Domain.VaccinationStatus.resolve(
            doneAt: model.doneAt,
            nextDueAt: model.nextDueAt
        )

        return Domain.Vaccination(
            id: model.id,
            dogId: model.dogId,
            name: model.name,
            doneAt: model.doneAt,
            nextDueAt: model.nextDueAt,
            notes: model.notes ?? "",
            status: status
        )
    }

    public static func to(model: Domain.Vaccination) -> VaccinationModel {
        return VaccinationModel(
            id: model.id,
            dogId: model.dogId,
            name: model.name,
            doneAt: model.doneAt,
            nextDueAt: model.nextDueAt,
            notes: model.notes.isEmpty ? nil : model.notes
        )
    }
}