import Foundation
import SwiftUI
import Combine
import AVFoundation

@MainActor
final class ChatDetailViewModel: ObservableObject {
    @Published var messages: [Message] = []
    @Published var messageText = ""
    @Published var pinnedMessages: [PinnedMessage] = []
    @Published var isRecordingVoice = false
    @Published var screenState: AppScreenState = .loading(
        message: "Загружаем сообщения..."
    )
    @Published var actionErrorMessage: String?

    private let chatID: UUID
    private let currentUserId: UUID
    private let useCase: ChatUseCase

    private var recorder: AVAudioRecorder?
    private var recordingURL: URL?

    var messageRows: [ChatMessageRowModel] {
        sortedMessages.map { message in
            ChatMessageRowModel(
                id: message.id,
                message: message,
                isMine: message.senderId == currentUserId,
                isPinned: pinnedMessages.contains {
                    $0.messageID == message.id
                }
            )
        }
    }

    private var sortedMessages: [Message] {
        messages.sorted {
            if $0.createdAt == $1.createdAt {
                return $0.id.uuidString < $1.id.uuidString
            }

            return $0.createdAt < $1.createdAt
        }
    }

    var latestMessage: Message? {
        sortedMessages.last
    }

    init(
        chatID: UUID,
        currentUserId: UUID,
        useCase: ChatUseCase
    ) {
        self.chatID = chatID
        self.currentUserId = currentUserId
        self.useCase = useCase
    }

    func loadMessages(showLoading: Bool = true) async {
        if showLoading && messages.isEmpty {
            screenState = .loading(
                message: "Загружаем сообщения..."
            )
        }

        do {
            messages = try await useCase.getMessages(
                chatID: chatID,
                limit: 50,
                before: nil
            )
            .sorted {
                $0.createdAt < $1.createdAt
            }

            screenState = .content
        } catch {
            if messages.isEmpty {
                screenState = .from(error)
            } else {
                screenState = .content
                actionErrorMessage = error.localizedDescription
            }
        }
    }

    func runPolling() async {
        while !Task.isCancelled {
            try? await Task.sleep(
                nanoseconds: 2_000_000_000
            )

            guard !Task.isCancelled else {
                return
            }

            await refreshMessagesSilently()
        }
    }

    func send() async {
        let text = messageText.trimmingCharacters(
            in: .whitespacesAndNewlines
        )

        guard !text.isEmpty else {
            return
        }

        actionErrorMessage = nil

        do {
            let message = try await useCase.sendMessage(
                chatID: chatID,
                text: text,
                attachments: nil,
                replyToId: nil
            )

            messageText = ""
            upsertMessage(message)
        } catch {
            actionErrorMessage = error.localizedDescription
        }
    }

    func sendAttachment(
        url: URL,
        type: AttachmentType,
        fileSize: Int64
    ) async {
        actionErrorMessage = nil

        do {
            let message = try await useCase.sendMessage(
                chatID: chatID,
                text: nil,
                attachments: [
                    AttachmentInput(
                        url: url.absoluteString,
                        type: type,
                        fileName: url.lastPathComponent,
                        fileSize: fileSize
                    )
                ],
                replyToId: nil
            )

            upsertMessage(message)
        } catch {
            actionErrorMessage = error.localizedDescription
        }
    }

    func sendImageData(_ data: Data) async {
        do {
            let url = try writeTemporaryFile(
                data: data,
                fileExtension: "jpg"
            )

            await sendAttachment(
                url: url,
                type: .image,
                fileSize: Int64(data.count)
            )
        } catch {
            actionErrorMessage = error.localizedDescription
        }
    }

    func toggleVoiceRecording() async {
        if isRecordingVoice {
            await stopVoiceRecording()
        } else {
            await startVoiceRecording()
        }
    }

    func setActionError(_ error: Error) {
        actionErrorMessage = error.localizedDescription
    }

    func clearActionError() {
        actionErrorMessage = nil
    }

    private func startVoiceRecording() async {
        do {
            let session = AVAudioSession.sharedInstance()

            try session.setCategory(
                .playAndRecord,
                mode: .default
            )
            try session.setActive(true)

            let allowed = await requestRecordPermission(
                session: session
            )

            guard allowed else {
                actionErrorMessage = "Нет доступа к микрофону"
                return
            }

            let url = try makeRecordingURL()

            let settings: [String: Any] = [
                AVFormatIDKey: Int(kAudioFormatMPEG4AAC),
                AVSampleRateKey: 44_100,
                AVNumberOfChannelsKey: 1,
                AVEncoderAudioQualityKey: AVAudioQuality.medium.rawValue
            ]

            recorder = try AVAudioRecorder(
                url: url,
                settings: settings
            )
            recorder?.record()

            recordingURL = url
            isRecordingVoice = true
        } catch {
            actionErrorMessage = error.localizedDescription
        }
    }

    private func stopVoiceRecording() async {
        recorder?.stop()
        recorder = nil
        isRecordingVoice = false

        guard let url = recordingURL else {
            return
        }

        recordingURL = nil

        let fileSize = (
            try? FileManager.default
                .attributesOfItem(atPath: url.path)[.size] as? NSNumber
        )?.int64Value ?? 0

        await sendAttachment(
            url: url,
            type: .audio,
            fileSize: fileSize
        )
    }

    private func requestRecordPermission(
        session: AVAudioSession
    ) async -> Bool {
        await withCheckedContinuation { continuation in
            if #available(iOS 17.0, *) {
                AVAudioApplication.requestRecordPermission { allowed in
                    continuation.resume(returning: allowed)
                }
            } else {
                session.requestRecordPermission { allowed in
                    continuation.resume(returning: allowed)
                }
            }
        }
    }

    private func writeTemporaryFile(
        data: Data,
        fileExtension: String
    ) throws -> URL {
        let url = FileManager.default.temporaryDirectory
            .appendingPathComponent(UUID().uuidString)
            .appendingPathExtension(fileExtension)

        try data.write(to: url, options: .atomic)

        return url
    }

    private func makeRecordingURL() throws -> URL {
        FileManager.default.temporaryDirectory
            .appendingPathComponent(UUID().uuidString)
            .appendingPathExtension("m4a")
    }

    private func upsertMessage(_ message: Message) {
        actionErrorMessage = nil

        if let index = messages.firstIndex(
            where: {
                $0.id == message.id
            }
        ) {
            messages[index] = message
        } else {
            messages.append(message)
        }

        messages.sort {
            $0.createdAt < $1.createdAt
        }

        screenState = .content
    }

    private func refreshMessagesSilently() async {
        do {
            let loaded = try await useCase.getMessages(
                chatID: chatID,
                limit: 50,
                before: nil
            )

            messages = loaded.sorted {
                $0.createdAt < $1.createdAt
            }

            if screenState != .content {
                screenState = .content
            }
        } catch {
            // Ошибка polling не скрывает уже загруженные сообщения
        }
    }
}
