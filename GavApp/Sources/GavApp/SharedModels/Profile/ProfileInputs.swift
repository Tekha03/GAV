import Foundation

public struct UpdateProfileInput: Encodable, Sendable {
    public let name: String?
    public let surname: String?
    public let username: String?
    public let profilePhotoUrl: String?
    public let bio: String?

    public init(
        name: String? = nil,
        surname: String? = nil,
        username: String? = nil,
        profilePhotoUrl: String? = nil,
        bio: String? = nil
    ) {
        self.name = name
        self.surname = surname
        self.username = username
        self.profilePhotoUrl = profilePhotoUrl
        self.bio = bio
    }

    private enum CodingKeys: String, CodingKey {
        case name
        case surname
        case username
        case profilePhotoUrl = "profile_photo_url"
        case bio
    }
}

public struct CreateProfileInput: Encodable {
    public let name: String
    public let surname: String
    public let username: String
    public let profilePhotoUrl: String?
    public let bio: String

    public init(
        name: String,
        surname: String,
        username: String,
        profilePhotoUrl: String? = nil,
        bio: String
    ) {
        self.name = name
        self.surname = surname
        self.username = username
        self.profilePhotoUrl = profilePhotoUrl
        self.bio = bio
    }

    private enum CodingKeys: String, CodingKey {
        case name
        case surname
        case username
        case profilePhotoUrl = "profile_photo_url"
        case bio
    }
}

extension CreateProfileInput {
    func toModel(userID: UUID) -> UserProfileModel {
        UserProfileModel(
            userId: userID,
            name: name,
            surname: surname,
            username: username,
            profilePhotoUrl: profilePhotoUrl,
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
