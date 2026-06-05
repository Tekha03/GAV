import SwiftUI

struct MessageBubbleView: View {
    let message: Message
    let isMine: Bool
    let isPinned: Bool

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
    }

    private func attachmentRow(_ attachment: Attachment) -> some View {
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