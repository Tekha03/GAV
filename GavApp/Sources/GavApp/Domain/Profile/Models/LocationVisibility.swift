public enum LocationVisibility: String, Equatable, Sendable, CaseIterable, Codable {
    case everyone = "VisibilityEveryone"
    case followersOnly = "VisibilityFollowersOnly"
    case noOne = "VisibilityNoOne"
}