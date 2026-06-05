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

        do {
            let (data, response) = try await session.data(for: request)

            guard let httpResponse = response as? HTTPURLResponse else {
                throw APIError.invalidResponse(statusCode: 0)
            }

            guard (200...299).contains(httpResponse.statusCode) else {
                throw decodeError(from: data, statusCode: httpResponse.statusCode)
            }

            return data
        } catch let error as APIError {
            throw error
        } catch {
            throw APIError.networkError(error)
        }
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

        do {
            let (data, response) = try await session.data(for: request)

            guard let httpResponse = response as? HTTPURLResponse else {
                throw APIError.invalidResponse(statusCode: 0)
            }

            guard (200...299).contains(httpResponse.statusCode) else {
                throw decodeError(from: data, statusCode: httpResponse.statusCode)
            }

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

    private func decodeError(from data: Data, statusCode: Int) -> APIError {
        let decoder = JSONDecoder()
        if let payload = try? decoder.decode(APIErrorPayload.self, from: data) {
            return .serverError(
                statusCode: statusCode,
                code: payload.error.code,
                message: payload.error.message
            )
        }

        return .invalidResponse(statusCode: statusCode)
    }
}

private struct APIErrorPayload: Decodable {
    let error: APIErrorBody
}

private struct APIErrorBody: Decodable {
    let code: String
    let message: String
}
