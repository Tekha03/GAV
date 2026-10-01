import Foundation
import UIKit

enum UploadImageProcessor {
    private static let maxDimension: CGFloat = 2_048
    private static let maxFileSize = 4_500_000

    static func jpegData(from sourceData: Data) throws -> Data {
        guard let sourceImage = UIImage(data: sourceData),
              sourceImage.size.width > 0,
              sourceImage.size.height > 0 else {
            throw UploadImageProcessingError.invalidImage
        }

        let longestSide = max(sourceImage.size.width, sourceImage.size.height)
        let scale = min(1, maxDimension / longestSide)
        let targetSize = CGSize(
            width: max(1, (sourceImage.size.width * scale).rounded()),
            height: max(1, (sourceImage.size.height * scale).rounded())
        )

        let format = UIGraphicsImageRendererFormat()
        format.scale = 1
        format.opaque = true

        let image = UIGraphicsImageRenderer(size: targetSize, format: format).image { context in
            UIColor.white.setFill()
            context.fill(CGRect(origin: .zero, size: targetSize))
            sourceImage.draw(in: CGRect(origin: .zero, size: targetSize))
        }

        for quality in [0.85, 0.7, 0.55, 0.4] {
            if let data = image.jpegData(compressionQuality: quality),
               data.count <= maxFileSize {
                return data
            }
        }

        throw UploadImageProcessingError.imageTooLarge
    }
}

private enum UploadImageProcessingError: LocalizedError {
    case invalidImage
    case imageTooLarge

    var errorDescription: String? {
        switch self {
        case .invalidImage:
            return "Не удалось обработать выбранное фото"
        case .imageTooLarge:
            return "Фото слишком большое"
        }
    }
}
