import Foundation
import SwiftUI
import Combine
import AVFoundation

@MainActor
final class ChatDetailViewModel: ObservableObject {
    @Published var messages: [Message] = []
    @Published var messageText = ""
    @Published var pinnedMessages: [PinnedMessage] = []
    @Published private(set) var members: [ChatMember] = []
    @Published var isRecordingVoice = false
    @Published var screenState: AppScreenState = .loading(
        message: "Загружаем сообщения..."
    )
    @Published var actionErrorMessage: String?

    private let chatID: UUID
    private let currentUserId: UUID
    private let useCase: ChatUseCase
    private let uploadService: UploadServiceAPIProtocol

    private var recorder: AVAudioRecorder?
    private var recordingURL: URL?

    var messageRows: [ChatMessageRowModel] {
        let orderedMessages = sortedMessages
        let messageIndices = Dictionary(
            uniqueKeysWithValues: orderedMessages.enumerated().map { ($1.id, $0) }
        )
        let otherMembers = members.filter { $0.userId != currentUserId }

        return orderedMessages.map { message in
            let messageIndex = messageIndices[message.id]
            let isRead = !otherMembers.isEmpty && otherMembers.allSatisfy { member in
                guard
                    let lastReadID = member.lastReadMessageId,
                    let lastReadIndex = messageIndices[lastReadID],
                    let messageIndex
                else {
                    return false
                }
                return lastReadIndex >= messageIndex
            }

            return ChatMessageRowModel(
                id: message.id,
                message: message,
                isMine: message.senderId == currentUserId,
                isPinned: pinnedMessages.contains {
                    $0.messageID == message.id
                },
                isRead: isRead
            )
        }
    }

    var messageSections: [ChatDaySection] {
        let calendar = Calendar.autoupdatingCurrent
        let grouped = Dictionary(grouping: messageRows) { row in
            calendar.startOfDay(for: row.message.createdAt)
        }

        return grouped.keys.sorted().map { date in
            ChatDaySection(date: date, rows: grouped[date] ?? [])
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
        useCase: ChatUseCase,
        uploadService: UploadServiceAPIProtocol
    ) {
        self.chatID = chatID
        self.currentUserId = currentUserId
        self.useCase = useCase
        self.uploadService = uploadService
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

            await markCurrentMessagesAsRead()
            await refreshMembersSilently()

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
            await refreshMembersSilently()
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
        data: Data,
        fileName: String,
        type: AttachmentType,
        mimeType: String?
    ) async {
        actionErrorMessage = nil

        do {
            let media = try await uploadService.uploadChatAttachment(
                data,
                fileName: fileName,
                mimeType: mimeType
            )
            let message = try await useCase.sendMessage(
                chatID: chatID,
                text: nil,
                attachments: [
                    AttachmentInput(
                        url: media.url,
                        type: type,
                        fileName: fileName,
                        fileSize: Int64(data.count)
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
            let jpegData = try UploadImageProcessor.jpegData(from: data)

            await sendAttachment(
                data: jpegData,
                fileName: "photo-\(UUID().uuidString).jpg",
                type: .image,
                mimeType: "image/jpeg"
            )
        } catch {
            actionErrorMessage = error.localizedDescription
        }
    }

    func deleteMessage(_ messageID: UUID) async {
        actionErrorMessage = nil
        do {
            try await useCase.deleteMessage(messageID: messageID)
            messages.removeAll { $0.id == messageID }
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

        do {
            let data = try Data(contentsOf: url)
            await sendAttachment(
                data: data,
                fileName: url.lastPathComponent,
                type: .audio,
                mimeType: "audio/mp4"
            )
        } catch {
            actionErrorMessage = error.localizedDescription
        }
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

            await markCurrentMessagesAsRead()

            if screenState != .content {
                screenState = .content
            }
        } catch {
            // Ошибка polling не скрывает уже загруженные сообщения
        }
    }

    private func markCurrentMessagesAsRead() async {
        guard messages.contains(where: { $0.senderId != currentUserId }) else {
            return
        }
        try? await useCase.markAsRead(chatID: chatID, userID: currentUserId)
    }

    private func refreshMembersSilently() async {
        if let loaded = try? await useCase.getChatMembers(chatID: chatID) {
            members = loaded
        }
    }
}
