import Foundation

public struct UserProfile: Identifiable, Equatable, Sendable {
    public let userId: UUID
    public let name: String
    public let surname: String
    public let username: String
    public let profilePhotoURL: String?
    public let bio: String
    public let address: String?
    public let birthDate: Date?

    public let lat: Double?
    public let lon: Double?
    public let locationStatus: LocationStatus
    public let locationVisibility: LocationVisibility

    public let showLocation: Bool
    public let isProfilePublic: Bool

    public var id: UUID { userId }

    public init(
        userId: UUID,
        name: String,
        surname: String,
        username: String,
        profilePhotoURL: String? = nil,
        bio: String,
        address: String? = nil,
        birthDate: Date? = nil,
        lat: Double? = nil,
        lon: Double? = nil,
        locationStatus: LocationStatus = .forcedOffline,
        locationVisibility: LocationVisibility = .noOne,
        showLocation: Bool = false,
        isProfilePublic: Bool = true
    ) {
        self.userId = userId
        self.name = name
        self.surname = surname
        self.username = username
        self.profilePhotoURL = profilePhotoURL
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
}