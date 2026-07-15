import Foundation

protocol UploadServiceAPIProtocol: Sendable {
    func uploadAvatar(
        _ imageData: Data,
        mimeType: String?
    ) async throws -> MediaInfoModel

    func uploadPostImage(
        _ imageData: Data,
        mimeType: String?
    ) async throws -> MediaInfoModel

    func uploadDogImage(_ imageData: Data, mimeType: String?) async throws -> MediaInfoModel
}

@available(macOS 12.0, *)
final class UploadServiceAPI: UploadServiceAPIProtocol, @unchecked Sendable {
    private let base: BaseAPI

    init(
        baseURL: URL,
        session: URLSession = .shared,
        authManager: AuthManager
    ) {
        self.base = BaseAPI(
            baseURL: baseURL,
            session: session,
            authManager: authManager
        )
    }

    func uploadAvatar(
        _ imageData: Data,
        mimeType: String?
    ) async throws -> MediaInfoModel {

        let data = try await base.upload(
            "/api/v1/upload/avatar",
            fileData: imageData,
            mimeType: mimeType,
            fieldName: "avatar",
            fileName: "avatar.jpg"
        )

        return try JSONDecoder().decode(MediaInfoModel.self, from: data)
    }

    func uploadPostImage(
        _ imageData: Data,
        mimeType: String?
    ) async throws -> MediaInfoModel {

        let data = try await base.upload(
            "/api/v1/upload/post-image",
            fileData: imageData,
            mimeType: mimeType,
            fieldName: "image",
            fileName: "post.jpg"
        )

        return try JSONDecoder().decode(MediaInfoModel.self, from: data)
    }

    func uploadDogImage(_ imageData: Data, mimeType: String?) async throws -> MediaInfoModel {
        let data = try await base.upload(
            "/api/v1/upload/dog-image",
            fileData: imageData,
            mimeType: mimeType,
            fieldName: "image",
            fileName: "dog.jpg"
        )
        return try JSONDecoder().decode(MediaInfoModel.self, from: data)
    }
}
