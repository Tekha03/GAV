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
        if let apiError = self as ? APIError {
            switch apiError {
                case .networkError(let underlyingError):
                    return underlyingError.isOfflineError

                case .invalidURL, .invalidResponse, .decodingError:
                    return false
            }
        }

        let nsError = self as NSError

        guard nsError.domain == NSURLErrorDomain else {
            return false
        }

        let offlineErrorCodes: Set<Int> = [
            NSURLErrorNotConnectedToInternet,
            NSURLErrorNetworkConnectionLost,
            NSURLErrorCannotFindHost,
            NSURLErrorCannotConnectToHost,
            NSURLErrorDNSLookupFailed,
            NSURLErrorTimedOut
        ]

        return offlineErrorCodes.contains(nsError.code)
    }
}
