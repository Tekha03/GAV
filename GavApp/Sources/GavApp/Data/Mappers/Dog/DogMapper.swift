public enum DogMapperError: Error {
    case invalidStatus
    case invalidAge
    case invalidGender
}

public struct DogMapper {
    public static func from(model: DogModel) throws -> Dog {
        let status = try DogStatusMapper.from(string: model.status)
        let age    = try DogAgeMapper.from(string: model.age)
        let gender = try DogGenderMapper.from(string: model.gender)

        return Dog(
            id: model.id,
            ownerId: model.ownerId,
            name: model.name,
            breed: model.breed,
            photoURL: model.photoUrl,
            status: status,
            age: age,
            gender: gender
        )
    }

    public static func to(model: Dog) -> DogModel {
        return DogModel(
            id: model.id,
            ownerId: model.ownerId,
            name: model.name,
            breed: model.breed,
            photoUrl: model.photoURL,
            status: DogStatusMapper.toString(model.status),
            age: DogAgeMapper.toString(model.age),
            gender: DogGenderMapper.toString(model.gender)
        )
    }
}
