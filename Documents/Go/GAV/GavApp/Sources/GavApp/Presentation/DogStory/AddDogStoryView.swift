import SwiftUI

@available(macOS 12.0, *)
struct AddDogStoryView: View {
    let onAddTapped: () -> Void

    var body: some View {
        VStack(spacing: 6) {
            Image(systemName: "plus.circle.fill")
                .resizable()
                .scaledToFit()
                .frame(width: 60, height: 60)
                .foregroundColor(.accentColor)

            Text("Добавить собаку")
                .font(.caption)
                .foregroundColor(.white)
        }
        .frame(width: 70)
        .onTapGesture {
            onAddTapped()
        }
    }
}