import SwiftUI
import UIKit

struct ImageCropEditor: View {
    let imageData: Data
    let aspectRatio: CGFloat
    let onCancel: () -> Void
    let onSave: (Data) -> Void

    @State private var zoom: CGFloat = 1
    @State private var lastZoom: CGFloat = 1
    @State private var offset: CGSize = .zero
    @State private var lastOffset: CGSize = .zero
    @State private var viewportSize: CGSize = .zero

    private var image: UIImage? { UIImage(data: imageData) }

    var body: some View {
        NavigationStack {
            GeometryReader { geometry in
                let width = geometry.size.width
                let size = CGSize(width: width, height: width / aspectRatio)

                ZStack {
                    Color.black

                    if let image {
                        Image(uiImage: image)
                            .resizable()
                            .scaledToFill()
                            .frame(width: size.width, height: size.height)
                            .scaleEffect(zoom)
                            .offset(offset)
                            .gesture(cropGesture(image: image, viewport: size))
                    }
                }
                .frame(width: size.width, height: size.height)
                .clipped()
                .overlay(Rectangle().stroke(.white.opacity(0.8), lineWidth: 1))
                .position(x: geometry.size.width / 2, y: geometry.size.height / 2)
                .onAppear { viewportSize = size }
                .onChange(of: size) { _, newSize in viewportSize = newSize }
            }
            .background(Color.black.ignoresSafeArea())
            .navigationTitle("Кадрирование")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .cancellationAction) {
                    Button("Отмена", action: onCancel)
                }
                ToolbarItem(placement: .confirmationAction) {
                    Button("Готово") { saveCrop() }
                }
            }
        }
        .preferredColorScheme(.dark)
    }

    private func cropGesture(image: UIImage, viewport: CGSize) -> some Gesture {
        SimultaneousGesture(
            MagnificationGesture()
                .onChanged { value in
                    zoom = min(max(lastZoom * value, 1), 5)
                    offset = clamped(offset, image: image, viewport: viewport, zoom: zoom)
                }
                .onEnded { _ in
                    lastZoom = zoom
                    lastOffset = offset
                },
            DragGesture()
                .onChanged { value in
                    let proposed = CGSize(
                        width: lastOffset.width + value.translation.width,
                        height: lastOffset.height + value.translation.height
                    )
                    offset = clamped(proposed, image: image, viewport: viewport, zoom: zoom)
                }
                .onEnded { _ in lastOffset = offset }
        )
    }

    private func clamped(
        _ proposed: CGSize,
        image: UIImage,
        viewport: CGSize,
        zoom: CGFloat
    ) -> CGSize {
        let baseScale = max(viewport.width / image.size.width, viewport.height / image.size.height)
        let displayedWidth = image.size.width * baseScale * zoom
        let displayedHeight = image.size.height * baseScale * zoom
        let maxX = max(0, (displayedWidth - viewport.width) / 2)
        let maxY = max(0, (displayedHeight - viewport.height) / 2)
        return CGSize(
            width: min(max(proposed.width, -maxX), maxX),
            height: min(max(proposed.height, -maxY), maxY)
        )
    }

    private func saveCrop() {
        guard let image, viewportSize.width > 0 else { return }

        let outputWidth: CGFloat = aspectRatio >= 1 ? 2_048 : 2_048 * aspectRatio
        let outputSize = CGSize(width: outputWidth, height: outputWidth / aspectRatio)
        let baseScale = max(
            viewportSize.width / image.size.width,
            viewportSize.height / image.size.height
        )
        let displayScale = baseScale * zoom
        let displayedSize = CGSize(
            width: image.size.width * displayScale,
            height: image.size.height * displayScale
        )
        let origin = CGPoint(
            x: (viewportSize.width - displayedSize.width) / 2 + offset.width,
            y: (viewportSize.height - displayedSize.height) / 2 + offset.height
        )
        let outputScale = outputSize.width / viewportSize.width

        let format = UIGraphicsImageRendererFormat()
        format.scale = 1
        format.opaque = true
        let cropped = UIGraphicsImageRenderer(size: outputSize, format: format).image { context in
            UIColor.black.setFill()
            context.fill(CGRect(origin: .zero, size: outputSize))
            image.draw(in: CGRect(
                x: origin.x * outputScale,
                y: origin.y * outputScale,
                width: displayedSize.width * outputScale,
                height: displayedSize.height * outputScale
            ))
        }

        if let data = cropped.jpegData(compressionQuality: 0.88) {
            onSave(data)
        }
    }
}
