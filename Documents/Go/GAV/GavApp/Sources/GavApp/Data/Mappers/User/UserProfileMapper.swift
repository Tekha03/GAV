import Foundation

public enum UserProfileMapperError: Error {
    case invalidLocationStatus
    case invalidLocationVisibility
}

public struct UserProfileMapper {
    public func from(model: UserProfileModel) throws -> UserProfile {
        let locationStatus: LocationStatus
        switch Int8(model.locationStatus) {
        case 0: locationStatus = .inactive
        case 1: locationStatus = .walking
        case 2: locationStatus = .forcedOffline
        default: throw UserProfileMapperError.invalidLocationStatus
        }

        let locationVisibility: LocationVisibility
        switch Int8(model.locationVisibility) {
        case 0: locationVisibility = .everyone
        case 1: locationVisibility = .followersOnly
        case 2: locationVisibility = .noOne
        default: throw UserProfileMapperError.invalidLocationVisibility
        }

        let birthDate: Date? = nil

        return UserProfile(
            userId: model.userId,
            name: model.name,
            surname: model.surname,
            username: model.username,
            profilePhotoURL: model.profilePhotoUrl,
            bio: model.bio,
            address: model.address,
            birthDate: birthDate,
            lat: model.lat,
            lon: model.lon,
            locationStatus: locationStatus,
            locationVisibility: locationVisibility,
            showLocation: model.showLocation,
            isProfilePublic: model.isProfilePublic
        )
    }

    public func to(model: UserProfile) -> UserProfileModel {
        let locationStatus: Double
        switch model.locationStatus {
        case .inactive:      locationStatus = 0
        case .walking:       locationStatus = 1
        case .forcedOffline: locationStatus = 2
        }

        let locationVisibility: Double
        switch model.locationVisibility {
        case .everyone:      locationVisibility = 0
        case .followersOnly: locationVisibility = 1
        case .noOne:         locationVisibility = 2
        }

        let birthDate: String? = nil

        return UserProfileModel(
            userId: model.userId,
            name: model.name,
            surname: model.surname,
            username: model.username,
            profilePhotoUrl: model.profilePhotoURL,
            bio: model.bio,
            address: model.address,
            birthDate: birthDate,
            lat: model.lat,
            lon: model.lon,
            locationStatus: locationStatus,
            locationVisibility: locationVisibility,
            showLocation: model.showLocation,
            isProfilePublic: model.isProfilePublic
        )
    }
}