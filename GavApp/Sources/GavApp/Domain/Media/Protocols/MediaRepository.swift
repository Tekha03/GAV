import Foundation

public protocol MediaRepository {
    func uploadImage(file: Data, mimeType: String?, folder: String) async throws -> MediaInfo
    func deleteMedia(url: String) async throws
}