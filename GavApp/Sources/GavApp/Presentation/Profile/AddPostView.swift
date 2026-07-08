import SwiftUI
import PhotosUI
#if os(iOS)
import UIKit
#endif

@available(macOS 12.0, *)
struct AddPostView: View {
    let viewModel: AppViewModel
    let uploadService: UploadServiceAPIProtocol
    @Environment(\.dismiss) private var dismiss

    @State private var content: String = ""
    @State private var selectedItem: PhotosPickerItem?
    @State private var selectedImageData: Data?
    @State private var imagePreview: Image?
    @State private var isLoadingImage = false
    @State private var isPublishing = false
    @State private var errorMessage: String?

    var body: some View {
        NavigationStack {
            Form {
                Section("Новый пост") {
                    TextEditor(text: $content)
                        .frame(minHeight: 140)

                    PhotosPicker(selection: $selectedItem, matching: .images) {
                        Label("Выбрать фото", systemImage: "photo")
                    }

                    if let imagePreview {
                        imagePreview
                            .resizable()
                            .scaledToFill()
                            .frame(height: 240)
                            .clipShape(RoundedRectangle(cornerRadius: 16))
                    }

                    if isLoadingImage {
                        ProgressView()
                    }
                }

                if let errorMessage {
                    Section {
                        Text(errorMessage)
                            .foregroundStyle(.red)
                    }
                }

                Section {
                    Button {
                        Task { await createPost() }
                    } label: {
                        if isPublishing {
                            ProgressView()
                        } else {
                            Text("Опубликовать")
                        }
                    }
                    .disabled(
                        isPublishing ||
                        content.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty
                    )
                }
            }
            .navigationTitle("Новый пост")
            .toolbar {
                ToolbarItem(placement: .automatic) {
                    Button("Готово") {
                        dismiss()
                    }
                }
            }
            .task(id: selectedItem) {
                await loadImage()
            }
        }
    }

    private func loadImage() async {
        guard let selectedItem else { return }
        isLoadingImage = true
        defer { isLoadingImage = false }

        do {
            if let data = try await selectedItem.loadTransferable(type: Data.self) {
                selectedImageData = data
                #if os(macOS)
                if let nsImage = NSImage(data: data) {
                    imagePreview = Image(nsImage: nsImage)
                } else {
                    imagePreview = nil
                }
                #else
                if let uiImage = UIImage(data: data) {
                    imagePreview = Image(uiImage: uiImage)
                } else {
                    imagePreview = nil
                }
                #endif
            }
        } catch {
            selectedImageData = nil
            imagePreview = nil
        }
    }

    private func createPost() async {
        isPublishing = true
        errorMessage = nil
        defer { isPublishing = false }

        var rawImageURL: String?
        if let selectedImageData {
            do {
                let media = try await uploadService.uploadPostImage(
                    selectedImageData,
                    mimeType: "image/jpeg"
                )
                rawImageURL = media.url
            } catch {
                errorMessage = "Не удалось загрузить фото поста"
                return
            }
        }

        do {
            try await viewModel.createPost(
                content: content.trimmingCharacters(in: .whitespacesAndNewlines),
                imageUrl: rawImageURL
            )
            dismiss()
        } catch {
            errorMessage = "Не удалось опубликовать пост"
        }
    }
}
