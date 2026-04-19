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
}

public struct SetLocationVisibilityInput: Encodable {
    public let visibility: LocationVisibility

    public init(visibility: LocationVisibility) {
        self.visibility = visibility
    }
}