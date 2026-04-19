import Foundation

public struct UpdateProfileInput: Encodable {
    public let name: String?
    public let surname: String?
    public let bio: String?

    public init(name: String? = nil, surname: String? = nil, bio: String? = nil) {
        self.name = name
        self.surname = surname
        self.bio = bio
    }
}

public struct CreateProfileInput: Encodable {
    public let name: String
    public let surname: String
    public let bio: String

    public init(name: String, surname: String, bio: String) {
        self.name = name
        self.surname = surname
        self.bio = bio
    }
}

extension CreateProfileInput {
    func toModel(userID: UUID) -> UserProfileModel {
        UserProfileModel(
            userId: userID,
            name: name,
            surname: surname,
            username: "",
            profilePhotoUrl: nil,
            bio: bio,
            address: nil,
            birthDate: nil,
            lat: nil,
            lon: nil,
            locationStatus: 0,
            locationVisibility: 0,
            showLocation: false,
            isProfilePublic: false
        )
    }
}