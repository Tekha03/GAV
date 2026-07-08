import SwiftUI

struct CommentsView: View {
    @EnvironmentObject private var appViewModel: AppViewModel
    @Environment(\.dismiss) private var dismiss

    let post: AppPost
    let onCommentsChanged: ((Int) -> Void)?

    init(post: AppPost, onCommentsChanged: ((Int) -> Void)? = nil) {
        self.post = post
        self.onCommentsChanged = onCommentsChanged
    }

    @State private var comments: [AppComment] = []
    @State private var draft = ""
    @State private var isLoading = false
    @State private var isSending = false
    @State private var errorMessage: String?

    var body: some View {
        NavigationStack {
            VStack(spacing: 0) {
                content
                composer
            }
            .background(Color.black.ignoresSafeArea())
            .navigationTitle("Комментарии")
            .toolbar {
                ToolbarItem(placement: .cancellationAction) {
                    Button("Готово") {
                        dismiss()
                    }
                }
            }
            .task {
                await loadComments()
            }
        }
        .preferredColorScheme(.dark)
    }

    @ViewBuilder
    private var content: some View {
        if isLoading {
            ProgressView()
                .frame(maxWidth: .infinity, maxHeight: .infinity)
        } else if comments.isEmpty {
            VStack(spacing: 12) {
                Image(systemName: "bubble.right")
                    .font(.system(size: 34, weight: .semibold))
                    .foregroundStyle(.orange)
                Text("Комментариев пока нет")
                    .font(.headline)
                Text("Стань первой в обсуждении.")
                    .font(.subheadline)
                    .foregroundStyle(.secondary)
            }
            .frame(maxWidth: .infinity, maxHeight: .infinity)
        } else {
            ScrollView {
                LazyVStack(alignment: .leading, spacing: 12) {
                    ForEach(comments) { comment in
                        commentRow(comment)
                    }
                }
                .padding(16)
            }
        }
    }

    private var composer: some View {
        VStack(spacing: 8) {
            if let errorMessage {
                Text(errorMessage)
                    .font(.footnote)
                    .foregroundStyle(.red)
                    .frame(maxWidth: .infinity, alignment: .leading)
            }

            HStack(alignment: .bottom, spacing: 10) {
                TextField("Комментарий", text: $draft, axis: .vertical)
                    .lineLimit(1...4)
                    .padding(.horizontal, 12)
                    .padding(.vertical, 10)
                    .background(.white.opacity(0.1), in: RoundedRectangle(cornerRadius: 8))

                Button {
                    Task { await sendComment() }
                } label: {
                    if isSending {
                        ProgressView()
                    } else {
                        Image(systemName: "paperplane.fill")
                            .font(.headline)
                    }
                }
                .frame(width: 42, height: 42)
                .background(canSend ? Color.orange : Color.gray.opacity(0.5), in: Circle())
                .foregroundStyle(.black)
                .disabled(!canSend || isSending)
            }
        }
        .padding(12)
        .background(.black)
        .overlay(alignment: .top) {
            Rectangle()
                .fill(.white.opacity(0.1))
                .frame(height: 1)
        }
    }

    private var canSend: Bool {
        !draft.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty
    }

    private func commentRow(_ comment: AppComment) -> some View {
        VStack(alignment: .leading, spacing: 5) {
            HStack {
                Text(comment.authorName)
                    .font(.subheadline.bold())
                Spacer()
                Text(comment.createdAt.formatted(date: .omitted, time: .shortened))
                    .font(.caption2)
                    .foregroundStyle(.secondary)
            }

            Text(comment.content)
                .font(.body)
        }
        .padding(12)
        .background(.white.opacity(0.08), in: RoundedRectangle(cornerRadius: 8))
    }

    private func loadComments() async {
        isLoading = true
        errorMessage = nil
        defer { isLoading = false }

        do {
            comments = try await appViewModel.loadComments(for: post.id)
            appViewModel.setCommentCount(postID: post.id, count: comments.count)
            onCommentsChanged?(comments.count)
        } catch {
            errorMessage = "Не удалось загрузить комментарии"
        }
    }

    private func sendComment() async {
        let content = draft.trimmingCharacters(in: .whitespacesAndNewlines)
        guard !content.isEmpty else { return }

        isSending = true
        errorMessage = nil
        defer { isSending = false }

        do {
            comments = try await appViewModel.addComment(to: post.id, content: content)
            onCommentsChanged?(comments.count)
            draft = ""
        } catch {
            errorMessage = "Не удалось отправить комментарий"
        }
    }
}
