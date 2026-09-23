import Foundation

enum AppScreenState: Equatable {
    case content
    case loading(message: String)
    case error(message: String)
    case offline

    static var loading: AppScreenState {
        return .loading(message: "Загрузка...")
    }
}

extension AppScreenState {
    static func from(_ error: Error) -> AppScreenState {
        if error.isOfflineError {
            return .offline
        }

        return .error(message: error.localizedDescription)
    }
}

private extension Error {
    var isOfflineError: Bool {
        if let apiError = self as? APIError {
            switch apiError {
                case .networkError(let underlyingError):
                    return underlyingError.isOfflineError

                case .invalidURL, .server, .invalidResponse, .decodingError:
                    return false
            }
        }

        let nsError = self as NSError

        guard nsError.domain == NSURLErrorDomain else {
            return false
        }

        // A reachable network and a reachable backend are different states.
        // Host lookup, connection and timeout failures usually mean that the
        // configured API is unavailable, not that the iPhone is offline.
        return nsError.code == NSURLErrorNotConnectedToInternet
    }
}
