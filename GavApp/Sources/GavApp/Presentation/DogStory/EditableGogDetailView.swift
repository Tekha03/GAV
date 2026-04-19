import SwiftUI

@available(macOS 12.0, *)
struct EditableDogStoryView: View {
    let dog: Dog
    let onEditTapped: (Dog) -> Void

    var body: some View {
        VStack(spacing: 6) {
            Image(dog.photoURL)
                .resizable()
                .scaledToFill()
                .frame(width: 60, height: 60)
                .clipShape(Circle())
                .overlay(
                    Circle().stroke(
                        Color(red: 220/255, green: 255/255, blue: 5/255),
                        lineWidth: 2.5
                    )
                )

            Text(dog.name)
                .font(.caption)
                .foregroundColor(.white)

            Image(systemName: "pencil.tip.crop.circle.badge.plus")
                .resizable()
                .scaledToFit()
                .foregroundColor(.accentColor)
                .frame(width: 18, height: 18)
                .offset(x: 0, y: -4)
                .onTapGesture {
                    onEditTapped(dog)
                }
        }
        .frame(width: 70)
    }
}