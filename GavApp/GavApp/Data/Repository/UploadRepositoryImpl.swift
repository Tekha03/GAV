import Foundation
import Domain
import Data
import SharedModels

final class UploadRepositoryImpl: UploadRepository {
    private let api: any UploadServiceAPIProtocol
    private let mapper: MediaInfoMapper

    init(api: any UploadServiceAPIProtocol, mapper: MediaInfoMapper = MediaInfoMapper()) {
        self.api = api
        self.mapper = mapper
    }

    func uploadAvatar(imageData: Data, mimeType: String? = nil) async throws -> MediaInfo {
        let model = try await api.uploadAvatar(imageData, mimeType: mimeType)
        return MediaInfoMapper.from(model: model)
    }

    func uploadPostImage(imageData: Data, mimeType: String? = nil) async throws -> MediaInfo {
        let model = try await api.uploadPostImage(imageData, mimeType: mimeType)
        return MediaInfoMapper.from(model: model)
    }
}