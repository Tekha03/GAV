import SwiftUI
import PhotosUI
import UniformTypeIdentifiers

struct ChatDetailView: View {
    let chat: Chat
    @StateObject private var viewModel: ChatDetailViewModel
    @State private var showingAttachmentOptions = false
    @State private var showingPhotoPicker = false
    @State private var showingFileImporter = false
    @State private var selectedPhoto: PhotosPickerItem?

    init(chat: Chat, currentUserId: UUID, useCase: ChatUseCase) {
        self.chat = chat
        _viewModel = StateObject(
            wrappedValue: ChatDetailViewModel(
                chatID: chat.id,
                currentUserId: currentUserId,
                useCase: useCase
            )
        )
    }

    var body: some View {
        VStack(spacing: 0) {
            header
            messagesView
            if let errorMessage = viewModel.errorMessage {
                Text(errorMessage)
                    .font(.footnote)
                    .foregroundStyle(.red.opacity(0.9))
                    .frame(maxWidth: .infinity, alignment: .leading)
                    .padding(.horizontal, 16)
                    .padding(.vertical, 8)
                    .background(.black.opacity(0.22))
            }
            composer
        }
        .background(
            LinearGradient(
                colors: [
                    Color(red: 0.42, green: 0.22, blue: 0.72),
                    .black
                ],
                startPoint: .top,
                endPoint: .bottom
            )
            .ignoresSafeArea()
        )
        .navigationTitle(chat.title)
        .navigationBarTitleDisplayMode(.inline)
        .preferredColorScheme(.dark)
        .task {
            await viewModel.loadMessages()
        }
        .task {
            await viewModel.runPolling()
        }
        .confirmationDialog("Вложение", isPresented: $showingAttachmentOptions, titleVisibility: .visible) {
            Button("Фото") {
                showingPhotoPicker = true
            }
            Button("Файл") {
                showingFileImporter = true
            }
            Button("Отмена", role: .cancel) {}
        }
        .photosPicker(
            isPresented: $showingPhotoPicker,
            selection: $selectedPhoto,
            matching: .images
        )
        .fileImporter(
            isPresented: $showingFileImporter,
            allowedContentTypes: [.item],
            allowsMultipleSelection: false
        ) { result in
            handleFileImport(result)
        }
        .onChange(of: selectedPhoto) { _, item in
            guard let item else { return }
            Task {
                if let data = try? await item.loadTransferable(type: Data.self) {
                    await viewModel.sendImageData(data)
                }
                selectedPhoto = nil
            }
        }
    }

    private var header: some View {
        HStack(spacing: 12) {
            Circle()
                .fill(Color.orange.opacity(0.18))
                .frame(width: 42, height: 42)
                .overlay(
                    Image(systemName: chat.isGroup ? "person.3.fill" : "bubble.left.and.bubble.right.fill")
                        .foregroundStyle(.orange)
                )

            VStack(alignment: .leading, spacing: 2) {
                Text(chat.title)
                    .font(.headline)
                    .foregroundStyle(.white)

                Text(chat.isGroup ? "Групповой чат" : "Личный чат")
                    .font(.caption)
                    .foregroundStyle(.white.opacity(0.7))
            }

            Spacer()

            Button {
                // открыть info
            } label: {
                Image(systemName: "info.circle")
                    .font(.title3)
                    .foregroundStyle(.white.opacity(0.9))
            }
        }
        .padding(.horizontal, 16)
        .padding(.vertical, 14)
        .background(.white.opacity(0.06))
    }

    private var messagesView: some View {
        ScrollViewReader { proxy in
            ScrollView {
                LazyVStack(spacing: 10) {
                    ForEach(viewModel.messageRows) { row in
                        MessageBubbleView(
                            message: row.message,
                            isMine: row.isMine,
                            isPinned: row.isPinned
                        )
                        .id(row.id)
                    }
                }
                .padding(.horizontal, 16)
                .padding(.vertical, 16)
            }
            .onChange(of: viewModel.messages.count) { _, _ in
                guard let last = viewModel.latestMessage else { return }
                DispatchQueue.main.async {
                    proxy.scrollTo(last.id, anchor: .bottom)
                }
            }
        }
    }

    private var composer: some View {
        HStack(spacing: 10) {
            Button {
                showingAttachmentOptions = true
            } label: {
                Image(systemName: "paperclip")
                    .font(.system(size: 18, weight: .semibold))
                    .foregroundStyle(.white)
                    .frame(width: 36, height: 36)
                    .background(.white.opacity(0.12), in: Circle())
            }

            Button {
                Task { await viewModel.toggleVoiceRecording() }
            } label: {
                Image(systemName: "mic.fill")
                    .font(.system(size: 18, weight: .semibold))
                    .foregroundStyle(.white)
                    .frame(width: 36, height: 36)
                    .background(viewModel.isRecordingVoice ? Color.red : .white.opacity(0.12), in: Circle())
            }

            TextField("Сообщение", text: $viewModel.messageText, axis: .vertical)
                .textFieldStyle(.plain)
                .foregroundStyle(.white)
                .padding(.horizontal, 14)
                .padding(.vertical, 10)
                .background(.white.opacity(0.10), in: RoundedRectangle(cornerRadius: 18))

            Button {
                Task { await viewModel.send() }
            } label: {
                Image(systemName: "arrow.up.circle.fill")
                    .font(.system(size: 28))
                    .foregroundStyle(viewModel.messageText.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty ? .gray : .orange)
            }
            .disabled(viewModel.messageText.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty)
        }
        .padding(.horizontal, 16)
        .padding(.vertical, 12)
        .background(.black.opacity(0.35))
    }

    private func handleFileImport(_ result: Result<[URL], Error>) {
        switch result {
        case .success(let urls):
            guard let url = urls.first else { return }
            let canAccess = url.startAccessingSecurityScopedResource()
            defer {
                if canAccess {
                    url.stopAccessingSecurityScopedResource()
                }
            }

            let fileSize = (try? FileManager.default.attributesOfItem(atPath: url.path)[.size] as? NSNumber)?
                .int64Value ?? 0
            Task {
                await viewModel.sendAttachment(
                    url: url,
                    type: attachmentType(for: url),
                    fileSize: fileSize
                )
            }
        case .failure(let error):
            viewModel.errorMessage = error.localizedDescription
        }
    }

    private func attachmentType(for url: URL) -> AttachmentType {
        guard let type = UTType(filenameExtension: url.pathExtension) else {
            return .document
        }
        if type.conforms(to: .image) {
            return .image
        }
        if type.conforms(to: .movie) || type.conforms(to: .video) {
            return .video
        }
        if type.conforms(to: .audio) {
            return .audio
        }
        return .document
    }
}
