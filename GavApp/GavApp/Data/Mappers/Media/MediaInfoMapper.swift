import Domain
import SharedModels

public struct MediaInfoMapper {
    public static func from(model: MediaInfoModel) -> Domain.MediaInfo {
        return Domain.MediaInfo(
            url: model.url,
            mimeType: model.mimeType
        )
    }

    public static func to(model: Domain.MediaInfo) -> MediaInfoModel {
        return MediaInfoModel(
            url: model.url,
            mimeType: model.mimeType
        )
    }
}