public enum DogStatusMapper {
    public static func from(string: String) throws -> DogStatus {
        switch string {
        case "friendly":   return .friendly
        case "cautious":   return .cautious
        case "aggressive": return .aggressive
        default:
            throw DogMapperError.invalidStatus
        }
    }

    public static func toString(_ status: DogStatus) -> String {
        switch status {
        case .friendly:   return "friendly"
        case .cautious:   return "cautious"
        case .aggressive: return "aggressive"
        }
    }
}