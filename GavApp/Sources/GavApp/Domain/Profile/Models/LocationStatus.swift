public enum LocationStatus: Int, Equatable, Sendable, CaseIterable, Codable {
    case inactive = 0
    case walking = 1
    case forcedOffline = 2

    public var stringValue: String {
        switch self {
        case .inactive:      return "Inactive"
        case .walking:       return "Walking"
        case .forcedOffline: return "ForcedOffline"
        }
    }

    public static func fromStringValue(_ value: String) -> LocationStatus? {
        switch value {
        case "Inactive":        return .inactive
        case "Walking":         return .walking
        case "ForcedOffline":   return .forcedOffline
        default: return nil
        }
    }
}