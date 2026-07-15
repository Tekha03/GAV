import Foundation

enum DogGenderMapper {

    static func from(string: String) throws -> DogGender {
        guard let value = DogGender(rawValue: string.lowercased()) else {
            throw NSError(
                domain: "DogGenderMapper",
                code: 1,
                userInfo: [NSLocalizedDescriptionKey: "Invalid DogGender: \(string)"]
            )
        }

        return value
    }

    static func toString(_ value: DogGender) -> String {
        value.rawValue
    }
}