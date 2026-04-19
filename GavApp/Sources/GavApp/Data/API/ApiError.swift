import Foundation

enum APIError: Error, LocalizedError, CustomNSError {
    case invalidURL
    case invalidResponse(statusCode: Int)
    case decodingError(Error)
    case networkError(Error)

    var errorDescription: String? {
        switch self {
        case .invalidURL:
            return "Invalid base URL"
        case .invalidResponse(let code):
            return "Invalid server response: \(code)"
        case .decodingError(_):
            return "Failed to decode JSON"
        case .networkError(_):
            return "Network error"
        }
    }
}