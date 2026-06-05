public struct MediaInfoMapper {
    public func from(model: MediaInfoModel) -> MediaInfo {
        return MediaInfo(
            url: model.url,
            mimeType: model.mimeType
        )
    }

    public func to(model: MediaInfo) -> MediaInfoModel {
        return MediaInfoModel(
            url: model.url,
            mimeType: model.mimeType
        )
    }
}
