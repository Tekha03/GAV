import SwiftUI

@available(macOS 12.0, *)
struct DogDetailView: View {
    let dog: Dog
    let onSave: (Dog) -> Void

    @Environment(\.dismiss) private var dismiss
    @State private var showEdit = false

    var body: some View {
        ZStack(alignment: .bottom) {
            Image(dog.photoURL)
                .resizable()
                .scaledToFill()
                .frame(maxWidth: .infinity, maxHeight: .infinity)
                .clipped()
                .ignoresSafeArea()

            LinearGradient(
                colors: [.clear, .black.opacity(0.3), .black.opacity(0.85)],
                startPoint: .center,
                endPoint: .bottom
            )
            .ignoresSafeArea()

            VStack(alignment: .center, spacing: 8) {
                Text(dog.name)
                    .font(.largeTitle.bold())
                    .foregroundStyle(.white)
                    .shadow(color: .black.opacity(0.5), radius: 4)
                    .multilineTextAlignment(.center)

                Text("\(DogPresentationMapper.breed(dog.breed)) • \(DogPresentationMapper.gender(dog.gender.rawValue)) • \(dog.age) года")
                    .font(.subheadline)
                    .foregroundStyle(.white.opacity(0.95))
                    .multilineTextAlignment(.center)

                Text(DogPresentationMapper.character(status: dog.status.rawValue, gender: dog.gender.rawValue))
                    .font(.headline)
                    .foregroundStyle(.white.opacity(0.95))
                    .multilineTextAlignment(.center)

                Button {
                    showEdit = true
                } label: {
                    Label("Редактировать", systemImage: "pencil")
                        .font(.headline.weight(.semibold))
                        .padding(.horizontal, 16)
                        .padding(.vertical, 10)
                        .background(.white.opacity(0.18), in: Capsule())
                        .foregroundStyle(.white)
                }
                .buttonStyle(.plain)
                .padding(.top, 12)
            }
            .padding(24)
            .padding(.bottom, 20)
            .frame(maxWidth: .infinity)
        }
        .sheet(isPresented: $showEdit) {
            EditDogView(dog: dog, onSave: onSave)
        }
    }
}