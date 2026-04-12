import Foundation
import SharedModels

struct BaseAPI {
    let baseURL: URL
    let session: URLSession
    let authManager: any AuthManager

    init(
        baseURL: URL,
        session: URLSession = .shared,
        authManager: any AuthManager
    ) {
        self.baseURL = baseURL
        self.session = session
        self.authManager = authManager
    }

    func request(
        _ path: String,
        method: String = "GET",
        body: Data? = nil,
        requiresAuth: Bool = true
    ) async throws -> Data {
        var components = URLComponents(url: baseURL, resolvingAgainstBaseURL: true)
        guard components != nil else {
            throw APIError.invalidURL
        }
        components?.path = path

        guard let url = components?.url else {
            throw APIError.invalidURL
        }

        var request = URLRequest(url: url)
        request.httpMethod = method
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")

        if requiresAuth, let token = authManager.currentToken {
            request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        }

        if let body = body {
            request.httpBody = body
        }

        do {
            let (data, response) = try await session.data(for: request)

            guard let httpResponse = response as? HTTPURLResponse,
                  (200...299).contains(httpResponse.statusCode)
            else {
                throw APIError.invalidResponse(statusCode: (response as? HTTPURLResponse)?.statusCode ?? 0)
            }

            return data
        } catch {
            throw APIError.networkError(error)
        }
    }
}