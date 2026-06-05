import Foundation

public struct VaccinationMapper {
    public static func from(model: VaccinationModel) -> Vaccination {
        let status = VaccinationStatus.resolve(
            doneAt: model.doneAt,
            nextDueAt: model.nextDueAt
        )

        return Vaccination(
            id: model.id,
            dogId: model.dogId,
            name: model.name,
            doneAt: model.doneAt,
            nextDueAt: model.nextDueAt,
            notes: model.notes ?? "",
            status: status
        )
    }

    public static func to(model: Vaccination) -> VaccinationModel {
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