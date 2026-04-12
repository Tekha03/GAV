import SwiftUI

struct DogDetailView: View {
    let dog: Dog
    @Environment(\.dismiss) private var dismiss

    var body: some View {
        ZStack(alignment: .bottom) { // меняем с .bottomLeading на .bottom
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
            
            VStack(alignment: .center, spacing: 8) { // меняем .leading на .center
                Text(dog.name)
                    .font(.largeTitle.bold())
                    .foregroundStyle(.white)
                    .shadow(color: .black.opacity(0.5), radius: 4)
                    .multilineTextAlignment(.center) // добавляем центрирование текста

                Text("\(DogPresentationMapper.breed(dog.breed)) • \(DogPresentationMapper.gender(dog.gender)) • \(dog.age) года")
                    .font(.subheadline)
                    .foregroundStyle(.white.opacity(0.95))
                    .multilineTextAlignment(.center)

                Text(DogPresentationMapper.character(status: dog.status, gender: dog.gender))
                    .font(.headline)
                    .foregroundStyle(.white.opacity(0.95))
                    .lineLimit(nil)
                    .multilineTextAlignment(.center)
            }
            .padding(24)
            .padding(.bottom, 20)
            .frame(maxWidth: .infinity) // добавляем чтобы VStack занимал всю ширину
        }
    }
}
