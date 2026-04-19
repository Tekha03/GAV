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