import Foundation

enum MediaURLResolver {
    private static let socialBaseURL = URL(string: "http://localhost:8080")!

    static func resolve(_ rawURL: String?) -> URL? {
        guard let rawURL, !rawURL.isEmpty else { return nil }

        if let url = URL(string: rawURL), url.scheme != nil {
            return url
        }

        if rawURL.hasPrefix("/") {
            return URL(string: rawURL, relativeTo: socialBaseURL)?.absoluteURL
        }

        return URL(string: "/" + rawURL, relativeTo: socialBaseURL)?.absoluteURL
    }
}
