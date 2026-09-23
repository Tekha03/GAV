import Foundation

@available(macOS 12.0, *)
struct BaseAPI: Sendable {
    let baseURL: URL
    let session: URLSession
    let authManager: AuthManager

    init(
        baseURL: URL,
        session: URLSession = .shared,
        authManager: AuthManager
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
        guard let url = makeURL(path) else {
            throw APIError.invalidURL
        }
        var request = URLRequest(url: url)
        request.httpMethod = method
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")

        if requiresAuth, let token = authManager.currentToken() {
            request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        }

        if let body = body {
            request.httpBody = body
        }

        return try await perform(request, requiresAuth: requiresAuth)
    }

    func upload(
        _ path: String,
        fileData: Data,
        mimeType: String?,
        fieldName: String = "file",
        fileName: String = "image.jpg",
        requiresAuth: Bool = true
    ) async throws -> Data {

        guard let url = makeURL(path) else {
            throw APIError.invalidURL
        }
        var request = URLRequest(url: url)
        request.httpMethod = "POST"

        let boundary = UUID().uuidString

        request.setValue(
            "multipart/form-data; boundary=\(boundary)",
            forHTTPHeaderField: "Content-Type"
        )

        if requiresAuth, let token = authManager.currentToken() {
            request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        }

        var body = Data()

        let type = mimeType ?? "image/jpeg"

        body.append("--\(boundary)\r\n".data(using: .utf8)!)
        body.append(
            "Content-Disposition: form-data; name=\"\(fieldName)\"; filename=\"\(fileName)\"\r\n"
                .data(using: .utf8)!
        )
        body.append("Content-Type: \(type)\r\n\r\n".data(using: .utf8)!)
        body.append(fileData)
        body.append("\r\n".data(using: .utf8)!)
        body.append("--\(boundary)--\r\n".data(using: .utf8)!)

        request.httpBody = body

        return try await perform(request, requiresAuth: requiresAuth)
    }

    private func perform(_ request: URLRequest, requiresAuth: Bool) async throws -> Data {
        do {
            let (data, response) = try await session.data(for: request)
            guard let response = response as? HTTPURLResponse else {
                throw APIError.invalidResponse(statusCode: 0)
            }

            if response.statusCode == 401 && requiresAuth {
                let failedToken = request.value(forHTTPHeaderField: "Authorization")
                    .flatMap { $0.hasPrefix("Bearer ") ? String($0.dropFirst(7)) : nil }
                let token = try await authManager.refreshAccessToken(
                    after: failedToken,
                    using: session
                )
                var retry = request
                retry.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
                let (retryData, retryResponse) = try await session.data(for: retry)
                guard let retryResponse = retryResponse as? HTTPURLResponse else {
                    throw APIError.invalidResponse(statusCode: 0)
                }
                try validate(retryResponse, data: retryData)
                return retryData
            }

            try validate(response, data: data)
            return data
        } catch let error as APIError {
            throw error
        } catch {
            throw APIError.networkError(error)
        }
    }

    private func makeURL(_ path: String) -> URL? {
        URL(string: path, relativeTo: baseURL)?.absoluteURL
    }

    private func validate(_ response: HTTPURLResponse, data: Data) throws {
        guard !(200...299).contains(response.statusCode) else { return }

        if let payload = try? JSONDecoder().decode(APIErrorResponse.self, from: data) {
            throw APIError.server(statusCode: response.statusCode, body: payload.error)
        }

        throw APIError.invalidResponse(statusCode: response.statusCode)
    }
}
