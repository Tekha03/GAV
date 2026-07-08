import Foundation

public struct UserProfileModel: Codable, Equatable {
    public let userId: UUID
    public let name: String
    public let surname: String
    public let username: String
    public let profilePhotoUrl: String?
    public let bio: String
    public let address: String?
    public let birthDate: String?

    public let lat: Double?
    public let lon: Double?
    public let locationStatus: Double
    public let locationVisibility: Double

    public let showLocation: Bool
    public let isProfilePublic: Bool

    public init(
        userId: UUID,
        name: String,
        surname: String,
        username: String,
        profilePhotoUrl: String?,
        bio: String,
        address: String?,
        birthDate: String?,
        lat: Double?,
        lon: Double?,
        locationStatus: Double,
        locationVisibility: Double,
        showLocation: Bool,
        isProfilePublic: Bool
    ) {
        self.userId = userId
        self.name = name
        self.surname = surname
        self.username = username
        self.profilePhotoUrl = profilePhotoUrl
        self.bio = bio
        self.address = address
        self.birthDate = birthDate
        self.lat = lat
        self.lon = lon
        self.locationStatus = locationStatus
        self.locationVisibility = locationVisibility
        self.showLocation = showLocation
        self.isProfilePublic = isProfilePublic
    }

    public init(from decoder: Decoder) throws {
        let container = try decoder.container(keyedBy: CodingKeys.self)
        userId = try container.decode(UUID.self, forKey: .userId)
        name = try container.decodeIfPresent(String.self, forKey: .name) ?? ""
        surname = try container.decodeIfPresent(String.self, forKey: .surname) ?? ""
        username = try container.decodeIfPresent(String.self, forKey: .username) ?? ""
        profilePhotoUrl = try container.decodeIfPresent(String.self, forKey: .profilePhotoUrl)
        bio = try container.decodeIfPresent(String.self, forKey: .bio) ?? ""
        address = try container.decodeIfPresent(String.self, forKey: .address)
        birthDate = try container.decodeIfPresent(String.self, forKey: .birthDate)
        lat = try container.decodeIfPresent(Double.self, forKey: .lat)
        lon = try container.decodeIfPresent(Double.self, forKey: .lon)
        locationStatus = try container.decodeIfPresent(Double.self, forKey: .locationStatus) ?? 0
        locationVisibility = try container.decodeIfPresent(Double.self, forKey: .locationVisibility) ?? 0
        showLocation = try container.decodeIfPresent(Bool.self, forKey: .showLocation) ?? false
        isProfilePublic = try container.decodeIfPresent(Bool.self, forKey: .isProfilePublic) ?? false
    }
}

extension UserProfileModel {
    private enum CodingKeys: String, CodingKey {
        case userId = "user_id"
        case name
        case surname
        case username
        case profilePhotoUrl = "profile_photo_url"
        case bio
        case address
        case birthDate = "birth_date"
        case lat
        case lon
        case locationStatus = "location_status"
        case locationVisibility = "location_visibility"
        case showLocation = "show_location"
        case isProfilePublic = "is_profile_public"
    }
}
