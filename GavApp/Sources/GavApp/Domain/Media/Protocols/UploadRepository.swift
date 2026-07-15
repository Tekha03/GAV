import Foundation

public protocol UploadRepository {
    func uploadAvatar(
        imageData: Data,
        mimeType: String?
    ) async throws -> MediaInfo

    func uploadPostImage(
        imageData: Data,
        mimeType: String?
    ) async throws -> MediaInfo
}