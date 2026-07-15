import Foundation

public struct UpdateLocationInput: Encodable {
    public let lat: Double?
    public let lon: Double?
    public let locationStatus: LocationStatus
    public let visibility: LocationVisibility

    public init(
        lat: Double? = nil,
        lon: Double? = nil,
        locationStatus: LocationStatus,
        visibility: LocationVisibility
    ) {
        self.lat = lat
        self.lon = lon
        self.locationStatus = locationStatus
        self.visibility = visibility
    }

    public func encode(to encoder: Encoder) throws {
        var container = encoder.container(keyedBy: CodingKeys.self)
        try container.encodeIfPresent(lat, forKey: .lat)
        try container.encodeIfPresent(lon, forKey: .lon)
        try container.encode(locationStatus.rawValue, forKey: .locationStatus)
        try container.encode(visibility.apiValue, forKey: .visibility)
    }

    private enum CodingKeys: String, CodingKey {
        case lat
        case lon
        case locationStatus = "location_status"
        case visibility
    }
}

public struct SetLocationVisibilityInput: Encodable {
    public let visibility: LocationVisibility

    public init(visibility: LocationVisibility) {
        self.visibility = visibility
    }

    public func encode(to encoder: Encoder) throws {
        var container = encoder.container(keyedBy: CodingKeys.self)
        try container.encode(visibility.apiValue, forKey: .visibility)
    }

    private enum CodingKeys: String, CodingKey {
        case visibility
    }
}

extension LocationVisibility {
    var apiValue: Int {
        switch self {
        case .everyone:
            return 0
        case .followersOnly:
            return 1
        case .noOne:
            return 2
        }
    }
}
