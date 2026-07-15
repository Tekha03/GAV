import Foundation

enum DogAgeMapper {

    static func from(string: String) throws -> DogAge {
        guard let value = DogAge(rawValue: string.lowercased()) else {
            throw NSError(
                domain: "DogAgeMapper",
                code: 1,
                userInfo: [NSLocalizedDescriptionKey: "Invalid DogAge: \(string)"]
            )
        }

        return value
    }

    static func toString(_ value: DogAge) -> String {
        value.rawValue
    }
}