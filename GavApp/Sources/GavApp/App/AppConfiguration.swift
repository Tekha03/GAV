import Foundation

enum AppEnvironment: String, Sendable {
    case debug = "Debug"
    case staging = "Staging"
    case production = "Production"
}

struct AppConfiguration : Sendable {
    let environment: AppEnvironment
    let socialBaseUrl: URL
    let messengerBaseUrl: URL

    static let current = AppConfiguration.load()

    private static func load(bundle: Bundle = .main) -> AppConfiguration {
        let environmentName = requiredString("GAVEnvironmentName", bundle: bundle)
        let socialURLString = requiredString("GAVSocialBaseURL", bundle: bundle)
        let messengerURLString = requiredString("GAVMessengerBaseURL", bundle: bundle)

        guard let environment = AppEnvironment(rawValue: environmentName) else {
            preconditionFailure("Invalid GAVEnvironmentName: \(environmentName)")
        }

        guard let socialBaseURL = URL(string: socialURLString), socialBaseURL.scheme != nil else {
            preconditionFailure("Invalid GAVSocialBaseURL: \(socialURLString)")
        }

        guard let messengerBaseURL = URL(string: messengerURLString), messengerBaseURL.scheme != nil else {
            preconditionFailure("Invalid GAVMessengerBaseURL: \(messengerURLString)")
        }

        return AppConfiguration(
            environment: environment,
            socialBaseUrl: socialBaseURL,
            messengerBaseUrl: messengerBaseURL
        )
    }

    private static func requiredString(_ key: String, bundle: Bundle) -> String {
        guard let value = bundle.object(forInfoDictionaryKey: key) as? String else {
            preconditionFailure("Missing Info.plist key: \(key)")
        }

        let trimmed = value.trimmingCharacters(in: .whitespacesAndNewlines)
        guard !trimmed.isEmpty, !trimmed.contains("$(") else {
            preconditionFailure("Invalid Info.plist value for \(key): \(value)")
        }

        return trimmed
    }
}
