import SwiftUI

struct CommentsView: View {
    @EnvironmentObject private var appViewModel: AppViewModel
    @Environment(\.dismiss) private var dismiss

    let post: AppPost
    let onCommentsChanged: ((Int) -> Void)?

    @State private var comments: [AppComment] = []
    @State private var draft = ""
    @State private var state: AppScreenState = .loading(
        message: "Загружаем комментарии..."
    )
    @State private var isSending = false
    @State private var actionErrorMessage: String?

    init(
        post: AppPost,
        onCommentsChanged: ((Int) -> Void)? = nil
    ) {
        self.post = post
        self.onCommentsChanged = onCommentsChanged
    }

    var body: some View {
        NavigationStack {
            VStack(spacing: 0) {
                switch state {
                case .loading, .error, .offline:
                    AppStatusView(
                        state: state,
                        retryAction: {
                            Task {
                                await loadComments()
                            }
                        }
                    )
                    .foregroundStyle(.white)

                case .content:
                    commentsContent
                    composer
                }
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
    private var commentsContent: some View {
        if comments.isEmpty {
            emptyState
        } else {
            commentsList
        }
    }

    private var emptyState: some View {
        VStack(spacing: 12) {
            Image(systemName: "bubble.right")
                .font(.system(size: 34, weight: .semibold))
                .foregroundStyle(.orange)

            Text("Комментариев пока нет")
                .font(.headline)
                .foregroundStyle(.white)

            Text("Стань первой в обсуждении.")
                .font(.subheadline)
                .foregroundStyle(.secondary)
        }
        .frame(maxWidth: .infinity, maxHeight: .infinity)
    }

    private var commentsList: some View {
        ScrollView {
            LazyVStack(
                alignment: .leading,
                spacing: 12
            ) {
                ForEach(comments) { comment in
                    commentRow(comment)
                }
            }
            .padding(16)
        }
        .refreshable {
            await loadComments(
                showLoading: false
            )
        }
    }

    private var composer: some View {
        VStack(spacing: 8) {
            if let actionErrorMessage {
                HStack(spacing: 10) {
                    Image(
                        systemName: "exclamationmark.circle.fill"
                    )
                    .foregroundStyle(.orange)

                    Text(actionErrorMessage)
                        .font(.footnote)
                        .foregroundStyle(.white)
                        .frame(
                            maxWidth: .infinity,
                            alignment: .leading
                        )

                    Button {
                        self.actionErrorMessage = nil
                    } label: {
                        Image(systemName: "xmark")
                            .font(.caption.bold())
                            .foregroundStyle(.white.opacity(0.7))
                    }
                }
            }

            HStack(alignment: .bottom, spacing: 10) {
                TextField(
                    "Комментарий",
                    text: $draft,
                    axis: .vertical
                )
                .lineLimit(1...4)
                .padding(.horizontal, 12)
                .padding(.vertical, 10)
                .background(
                    .white.opacity(0.1),
                    in: RoundedRectangle(cornerRadius: 8)
                )

                Button {
                    Task {
                        await sendComment()
                    }
                } label: {
                    if isSending {
                        ProgressView()
                            .tint(.black)
                    } else {
                        Image(systemName: "paperplane.fill")
                            .font(.headline)
                    }
                }
                .frame(width: 42, height: 42)
                .background(
                    canSend
                        ? Color.orange
                        : Color.gray.opacity(0.5),
                    in: Circle()
                )
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
        !draft.trimmingCharacters(
            in: .whitespacesAndNewlines
        )
        .isEmpty
    }

    private func commentRow(
        _ comment: AppComment
    ) -> some View {
        VStack(alignment: .leading, spacing: 5) {
            HStack {
                Text(comment.authorName)
                    .font(.subheadline.bold())

                Spacer()

                Text(
                    comment.createdAt.formatted(
                        date: .omitted,
                        time: .shortened
                    )
                )
                .font(.caption2)
                .foregroundStyle(.secondary)
            }

            Text(comment.content)
                .font(.body)
        }
        .padding(12)
        .background(
            .white.opacity(0.08),
            in: RoundedRectangle(cornerRadius: 8)
        )
    }

    private func loadComments(
        showLoading: Bool = true
    ) async {
        if showLoading && comments.isEmpty {
            state = .loading(
                message: "Загружаем комментарии..."
            )
        }

        actionErrorMessage = nil

        do {
            let loadedComments = try await appViewModel.loadComments(
                for: post.id
            )

            comments = loadedComments

            appViewModel.setCommentCount(
                postID: post.id,
                count: loadedComments.count
            )

            onCommentsChanged?(
                loadedComments.count
            )

            state = .content
        } catch {
            if comments.isEmpty {
                state = .from(error)
            } else {
                state = .content
                actionErrorMessage = error.localizedDescription
            }
        }
    }

    private func sendComment() async {
        let content = draft.trimmingCharacters(
            in: .whitespacesAndNewlines
        )

        guard !content.isEmpty else {
            return
        }

        isSending = true
        actionErrorMessage = nil

        defer {
            isSending = false
        }

        do {
            comments = try await appViewModel.addComment(
                to: post.id,
                content: content
            )

            onCommentsChanged?(
                comments.count
            )

            draft = ""
            state = .content
        } catch {
            actionErrorMessage = error.localizedDescription
        }
    }
}
