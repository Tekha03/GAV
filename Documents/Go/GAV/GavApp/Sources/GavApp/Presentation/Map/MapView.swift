import SwiftUI
import MapKit

struct MapView: View {
    @EnvironmentObject private var appViewModel: AppViewModel
    @State private var position = MapCameraPosition.region(
        MKCoordinateRegion(
            center: CLLocationCoordinate2D(latitude: 55.751244, longitude: 37.618423),
            span: MKCoordinateSpan(latitudeDelta: 0.03, longitudeDelta: 0.03)
        )
    )

    var body: some View {
        NavigationStack {
            ScrollView {
                VStack(spacing: 18) {
                    mapCard

                    infoCard
                }
                .padding(.horizontal, 16)
                .padding(.vertical, 20)
            }
            .background(
                LinearGradient(
                    colors: [
                        Color(red: 0.42, green: 0.22, blue: 0.72),
                        .black
                    ],
                    startPoint: .top,
                    endPoint: .bottom
                )
                .ignoresSafeArea()
            )
            .navigationTitle("Карта")
            .preferredColorScheme(.dark)
        }
    }

    private var mapCard: some View {
        VStack(alignment: .leading, spacing: 12) {
            HStack {
                Text("Прогулки рядом")
                    .font(.headline)
                    .foregroundStyle(.white)

                Spacer()

                Text("\(appViewModel.walkers.count) рядом")
                    .font(.caption.weight(.semibold))
                    .foregroundStyle(.white.opacity(0.8))
                    .padding(.horizontal, 10)
                    .padding(.vertical, 6)
                    .background(.white.opacity(0.10), in: Capsule())
            }

            Map(position: $position) {
                ForEach(appViewModel.walkers) { walker in
                    Annotation(
                        walker.dogName,
                        coordinate: CLLocationCoordinate2D(latitude: walker.latitude, longitude: walker.longitude)
                    ) {
                        VStack(spacing: 6) {
                            Image(systemName: "flag.fill")
                                .font(.title3)
                                .foregroundStyle(walker.mood.color)
                                .shadow(color: .black.opacity(0.25), radius: 4, x: 0, y: 2)

                            Text(walker.dogName)
                                .font(.caption2.bold())
                                .foregroundStyle(.white)
                                .padding(.horizontal, 8)
                                .padding(.vertical, 4)
                                .background(.black.opacity(0.45), in: Capsule())
                        }
                    }
                }
            }
            .frame(height: 540)
            .clipShape(RoundedRectangle(cornerRadius: 28, style: .continuous))
            .overlay(
                RoundedRectangle(cornerRadius: 28, style: .continuous)
                    .stroke(.white.opacity(0.08), lineWidth: 1)
            )
            .shadow(color: .black.opacity(0.28), radius: 18, x: 0, y: 10)
        }
        .padding(18)
        .background(
            RoundedRectangle(cornerRadius: 28, style: .continuous)
                .fill(.white.opacity(0.08))
        )
        .overlay(
            RoundedRectangle(cornerRadius: 28, style: .continuous)
                .stroke(.white.opacity(0.08), lineWidth: 1)
        )
    }

    private var infoCard: some View {
        VStack(alignment: .leading, spacing: 8) {
            Text("Значения флажков")
                .font(.headline)
                .foregroundStyle(.white)

            HStack(spacing: 8) {
                Circle().fill(.red).frame(width: 10, height: 10)
                Text("Красный флаг — агрессивный.")
                    .foregroundStyle(.white.opacity(0.75))
            }

            HStack(spacing: 8) {
                Circle().fill(.yellow).frame(width: 10, height: 10)
                Text("Жёлтый — настороженный.")
                    .foregroundStyle(.white.opacity(0.75))
            }

            HStack(spacing: 8) {
                Circle().fill(.green).frame(width: 10, height: 10)
                Text("Зелёный — дружелюбный.")
                    .foregroundStyle(.white.opacity(0.75))
            }
        }
        .padding(18)
        .background(.white.opacity(0.08), in: RoundedRectangle(cornerRadius: 24, style: .continuous))
    }
}