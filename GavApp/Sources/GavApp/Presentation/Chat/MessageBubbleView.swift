import SwiftUI
import AVFoundation

struct MessageBubbleView: View {
    let message: Message
    let isMine: Bool
    let isPinned: Bool
    let isRead: Bool
    let onDelete: (() -> Void)?

    @State private var previewImage: ChatImagePreview?

    var body: some View {
        HStack {
            if isMine { Spacer(minLength: 40) }

            VStack(alignment: .leading, spacing: 8) {
                if isPinned {
                    HStack(spacing: 6) {
                        Image(systemName: "pin.fill")
                        Text("Закреплено")
                    }
                    .font(.caption2)
                    .foregroundStyle(.orange)
                }

                if let text = message.text, !text.isEmpty {
                    Text(text)
                        .font(.body)
                        .foregroundStyle(.white)
                }

                if !message.attachments.isEmpty {
                    VStack(alignment: .leading, spacing: 6) {
                        ForEach(message.attachments) { attachment in
                            attachmentRow(attachment)
                        }
                    }
                }

                HStack {
                    Text(message.createdAt.formatted(date: .omitted, time: .shortened))
                        .font(.caption2)
                        .foregroundStyle(.white.opacity(0.6))

                    if message.editedAt != nil {
                        Text("изменено")
                            .font(.caption2)
                            .foregroundStyle(.white.opacity(0.5))
                    }

                    if isMine {
                        if isRead {
                            Text("Прочитано")
                                .font(.caption2)
                                .foregroundStyle(Color.blue)
                        } else {
                            Image(systemName: "checkmark")
                                .font(.caption2)
                                .foregroundStyle(.white.opacity(0.55))
                                .accessibilityLabel("Отправлено")
                        }
                    }

                    Spacer()

                    if !message.reactions.isEmpty {
                        reactionRow(message.reactions)
                    }
                }
            }
            .padding(12)
            .background(
                isMine ? Color.orange.opacity(0.22) : Color.white.opacity(0.10),
                in: RoundedRectangle(cornerRadius: 18, style: .continuous)
            )
            .overlay(
                RoundedRectangle(cornerRadius: 18, style: .continuous)
                    .strokeBorder(.white.opacity(0.08), lineWidth: 1)
            )

            if !isMine { Spacer(minLength: 40) }
        }
        .contextMenu {
            if let onDelete {
                Button(role: .destructive, action: onDelete) {
                    Label("Удалить", systemImage: "trash")
                }
            }
        }
        .fullScreenCover(item: $previewImage) { preview in
            ChatImageViewer(url: preview.url)
        }
    }

    private func attachmentRow(_ attachment: Attachment) -> some View {
        Group {
            if attachment.type == .image,
               let url = MediaURLResolver.resolve(attachment.url) {
                Button {
                    previewImage = ChatImagePreview(url: url)
                } label: {
                    VStack(alignment: .leading, spacing: 6) {
                        AsyncImage(url: url) { phase in
                            switch phase {
                            case .success(let image):
                                image.resizable().scaledToFill()
                            default:
                                RoundedRectangle(cornerRadius: 10)
                                    .fill(.white.opacity(0.08))
                                    .overlay { ProgressView().tint(.white) }
                            }
                        }
                        .frame(width: 220, height: 160)
                        .clipShape(RoundedRectangle(cornerRadius: 10))

                    }
                }
            } else if attachment.type == .audio,
                      let url = MediaURLResolver.resolve(attachment.url) {
                ChatAudioAttachmentView(
                    url: url,
                    fileName: attachment.fileName,
                    fileSize: attachment.fileSize
                )
            } else {
                attachmentContent(attachment)
            }
        }
        .buttonStyle(.plain)
    }

    private func attachmentContent(_ attachment: Attachment) -> some View {
        HStack(spacing: 8) {
            Image(systemName: iconName(for: attachment.type))
                .foregroundStyle(.orange)

            VStack(alignment: .leading, spacing: 2) {
                Text(attachment.fileName)
                    .font(.subheadline)
                    .foregroundStyle(.white)

                Text(formatSize(attachment.fileSize))
                    .font(.caption2)
                    .foregroundStyle(.white.opacity(0.6))
            }

            Spacer()

        }
        .padding(10)
        .background(.white.opacity(0.08), in: RoundedRectangle(cornerRadius: 12))
    }

    private func reactionRow(_ reactions: [Reaction]) -> some View {
        let grouped = Dictionary(grouping: reactions, by: { $0.emoji })
        return HStack(spacing: 4) {
            ForEach(grouped.keys.sorted(), id: \.self) { emoji in
                Text("\(emoji) \(grouped[emoji]?.count ?? 0)")
                    .font(.caption2)
                    .padding(.horizontal, 6)
                    .padding(.vertical, 3)
                    .background(.white.opacity(0.10), in: Capsule())
            }
        }
    }

    private func iconName(for type: AttachmentType) -> String {
        switch type {
        case .image: return "photo"
        case .video: return "video"
        case .audio: return "waveform"
        case .document: return "doc"
        }
    }

    private func formatSize(_ size: Int64) -> String {
        if size > 1_000_000 {
            return String(format: "%.1f MB", Double(size) / 1_000_000)
        } else if size > 1_000 {
            return String(format: "%.1f KB", Double(size) / 1_000)
        } else {
            return "\(size) B"
        }
    }
}

private struct ChatImagePreview: Identifiable {
    let id = UUID()
    let url: URL
}

private struct ChatImageViewer: View {
    let url: URL
    @Environment(\.dismiss) private var dismiss

    var body: some View {
        ZStack {
            Color.black.ignoresSafeArea()

            AsyncImage(url: url) { phase in
                switch phase {
                case .success(let image):
                    image
                        .resizable()
                        .scaledToFit()
                        .frame(maxWidth: .infinity, maxHeight: .infinity)
                case .failure:
                    ContentUnavailableView(
                        "Не удалось открыть фото",
                        systemImage: "photo.badge.exclamationmark"
                    )
                    .foregroundStyle(.white)
                default:
                    ProgressView().tint(.white)
                }
            }

            VStack {
                HStack {
                    Spacer()
                    Button {
                        dismiss()
                    } label: {
                        Image(systemName: "xmark.circle.fill")
                            .font(.system(size: 30))
                            .symbolRenderingMode(.palette)
                            .foregroundStyle(.white, .black.opacity(0.55))
                    }
                    .padding(16)
                }
                Spacer()
            }
        }
        .preferredColorScheme(.dark)
    }
}

private struct ChatAudioAttachmentView: View {
    let url: URL
    let fileName: String
    let fileSize: Int64

    @State private var player: AVPlayer?
    @State private var isPlaying = false

    var body: some View {
        HStack(spacing: 10) {
            Button {
                togglePlayback()
            } label: {
                Image(systemName: isPlaying ? "pause.fill" : "play.fill")
                    .frame(width: 34, height: 34)
                    .background(Color.orange, in: Circle())
                    .foregroundStyle(.black)
            }
            .buttonStyle(.plain)

            VStack(alignment: .leading, spacing: 3) {
                Text("Голосовое сообщение")
                    .font(.subheadline)
                    .foregroundStyle(.white)
                Text(fileSizeText)
                    .font(.caption2)
                    .foregroundStyle(.white.opacity(0.6))
            }

            Spacer()
        }
        .padding(10)
        .background(.white.opacity(0.08), in: RoundedRectangle(cornerRadius: 12))
        .onDisappear {
            player?.pause()
            isPlaying = false
        }
        .onReceive(NotificationCenter.default.publisher(for: .AVPlayerItemDidPlayToEndTime)) { notification in
            guard let item = notification.object as? AVPlayerItem,
                  item === player?.currentItem else { return }
            player?.seek(to: .zero)
            isPlaying = false
        }
    }

    private var fileSizeText: String {
        if fileSize > 1_000_000 {
            return String(format: "%.1f MB", Double(fileSize) / 1_000_000)
        }
        return String(format: "%.1f KB", Double(fileSize) / 1_000)
    }

    private func togglePlayback() {
        if isPlaying {
            player?.pause()
            isPlaying = false
            return
        }

        if player == nil {
            player = AVPlayer(url: url)
        }
        player?.play()
        isPlaying = true
    }
}
