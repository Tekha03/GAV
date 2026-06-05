import SwiftUI
import PhotosUI

@available(macOS 12.0, *)
struct AddDogView: View {
    @EnvironmentObject private var viewModel: AppViewModel
    @Environment(\.dismiss) private var dismiss

    let editingDog: AppDog?
    let uploadService: UploadServiceAPIProtocol

    @State private var name: String
    @State private var breed: String
    @State private var ageText: String
    @State private var mood: DogMood
    @State private var notes: String

    @State private var selectedPhoto: PhotosPickerItem?
    @State private var selectedImageData: Data?
    @State private var isUploading = false
    @State private var uploadError: String?

    init(viewModel: AppViewModel,
         editingDog: AppDog? = nil,
         uploadService: UploadServiceAPIProtocol) {
        self.editingDog = editingDog
        self.uploadService = uploadService

        _name = State(initialValue: editingDog?.name ?? "")
        _breed = State(initialValue: editingDog?.breed ?? "")
        _ageText = State(initialValue: editingDog?.ageText ?? "")
        _mood = State(initialValue: editingDog?.mood ?? .friendly)
        _notes = State(initialValue: editingDog?.notes ?? "")
    }

    var body: some View {
        NavigationStack {
            Form {
                Section("Фото") {
                    PhotosPicker(selection: $selectedPhoto, matching: .images) {
                        HStack {
                            Image(systemName: "photo")
                            Text(selectedImageData == nil ? "Выбрать фото" : "Фото выбрано")
                        }
                    }

                    if let selectedImageData,
                       let uiImage = UIImage(data: selectedImageData) {
                        Image(uiImage: uiImage)
                            .resizable()
                            .scaledToFill()
                            .frame(height: 180)
                            .clipShape(RoundedRectangle(cornerRadius: 16))
                    } else if let editingDog,
                              let url = editingDog.photoURL {
                        AsyncImage(url: url) { phase in
                            switch phase {
                            case .success(let image):
                                image
                                    .resizable()
                                    .scaledToFill()
                            default:
                                ZStack {
                                    RoundedRectangle(cornerRadius: 16)
                                        .fill(Color.gray.opacity(0.2))
                                    Image(systemName: "dog.fill")
                                        .font(.largeTitle)
                                }
                            }
                        }
                        .frame(height: 180)
                        .clipShape(RoundedRectangle(cornerRadius: 16))
                    }
                }

                Section("Собака") {
                    TextField("Имя", text: $name)
                    TextField("Порода", text: $breed)
                    TextField("Возраст", text: $ageText)
                    Picker("Настроение", selection: $mood) {
                        ForEach(DogMood.allCases) { mood in
                            Text(mood.title).tag(mood)
                        }
                    }
                    TextField("Заметки", text: $notes, axis: .vertical)
                }

                if let uploadError {
                    Section {
                        Text(uploadError)
                            .foregroundStyle(.red)
                    }
                }

                Section {
                    Button {
                        Task { await save() }
                    } label: {
                        HStack {
                            Spacer()
                            if isUploading {
                                ProgressView()
                            } else {
                                Text(editingDog == nil ? "Добавить" : "Сохранить")
                            }
                            Spacer()
                        }
                    }
                    .disabled(isUploading)
                }
            }
            .navigationTitle(editingDog == nil ? "Новая собака" : "Редактировать")
            .toolbar {
                ToolbarItem(placement: .cancellationAction) {
                    Button("Отмена") { dismiss() }
                }
            }
            .onChange(of: selectedPhoto) { _, newValue in
                Task {
                    if let data = try? await newValue?.loadTransferable(type: Data.self) {
                        selectedImageData = data
                    }
                }
            }
        }
    }

    private func save() async {
        isUploading = true
        uploadError = nil
        defer { isUploading = false }

        var photoURL: URL? = editingDog?.photoURL

        if let selectedImageData {
            do {
                let media = try await uploadService.uploadDogImage(
                    selectedImageData,
                    mimeType: "image/jpeg"
                )

                guard let newURL = MediaURLResolver.resolve(media.url) else {
                    uploadError = "Некорректный URL фото"
                    return
                }

                photoURL = newURL
            } catch {
                uploadError = "Не удалось загрузить фото"
                return
            }
        }

        let dog = AppDog(
            id: editingDog?.id ?? UUID(),
            name: name,
            breed: breed,
            ageText: ageText,
            mood: mood,
            photoURL: photoURL,
            notes: notes
        )

        if let editingDog,
           let index = viewModel.dogs.firstIndex(where: { $0.id == editingDog.id }) {
            viewModel.dogs[index] = dog
        } else {
            viewModel.dogs.insert(dog, at: 0)
        }

        dismiss()
    }
}
