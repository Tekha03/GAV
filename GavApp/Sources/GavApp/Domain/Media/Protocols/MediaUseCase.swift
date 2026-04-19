import Foundation

public struct MediaUseCase {
    private let repository: any MediaRepository

    public init(repository: any MediaRepository) {
        self.repository = repository
    }

    public func uploadImage(
        file: Data,
        mimeType: String? = nil,
        folder: String
    ) async throws -> MediaInfo {
        try await repository.uploadImage(file: file, mimeType: mimeType, folder: folder)
    }

    public func deleteMedia(url: String) async throws {
        try await repository.deleteMedia(url: url)
    }
}