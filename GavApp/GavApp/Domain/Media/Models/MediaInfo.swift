public struct MediaInfo: Equatable, Sendable {
    public let url: String
    public let mimeType: String?

    public init(
        url: String,
        mimeType: String? = nil
    ) {
        self.url = url
        self.mimeType = mimeType
    }
}