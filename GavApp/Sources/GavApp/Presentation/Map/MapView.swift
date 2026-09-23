import SwiftUI
import MapKit

struct MapView: View {
    @EnvironmentObject private var appViewModel: AppViewModel
    @StateObject private var userLocation = UserLocation()
    @State private var position = MapCameraPosition.region(
        MKCoordinateRegion(
            center: CLLocationCoordinate2D(latitude: 55.751244, longitude: 37.618423),
            span: MKCoordinateSpan(latitudeDelta: 0.03, longitudeDelta: 0.03)
        )
    )
    @State private var screenState: AppScreenState = .content
    @State private var isRefreshing = false
    @State private var actionErrorMessage: String?
    @State private var didCenterOnUser = false

    var body: some View {
        NavigationStack {
            ZStack {
                mapBackground

                switch screenState {
                case .loading, .error, .offline:
                    AppStatusView(
                        state: screenState,
                        retryAction: {
                            Task {
                                await refreshNearby()
                            }
                        }
                    )
                    .foregroundStyle(.white)

                case .content:
                    mapContent
                }
            }
            .navigationTitle("Карта")
            .preferredColorScheme(.dark)
            .onAppear {
                userLocation.startUpdatingLocation()
            }
            .onReceive(userLocation.$location) { location in
                guard
                    !didCenterOnUser,
                    let coordinate = location?.coordinate
                else {
                    return
                }

                didCenterOnUser = true
                centerMap(on: coordinate)
            }
        }
    }

    private var mapBackground: some View {
        LinearGradient(
            colors: [
                Color(red: 0.42, green: 0.22, blue: 0.72),
                .black
            ],
            startPoint: .top,
            endPoint: .bottom
        )
        .ignoresSafeArea()
    }

    private var mapContent: some View {
        ScrollView {
            VStack(spacing: 18) {
                mapCard

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
                                .foregroundStyle(
                                    .white.opacity(0.7)
                                )
                        }
                    }
                    .padding(12)
                    .background(
                        .black.opacity(0.35),
                        in: RoundedRectangle(cornerRadius: 14)
                    )
                }

                infoCard
            }
            .padding(.horizontal, 16)
            .padding(.vertical, 20)
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
                if let currentCoordinate = userLocation.location?.coordinate {
                    Annotation("Вы", coordinate: currentCoordinate) {
                        VStack(spacing: 6) {
                            Image(systemName: "location.circle.fill")
                                .font(.system(size: 34, weight: .semibold))
                                .foregroundStyle(.blue)
                                .background(.white, in: Circle())
                                .shadow(color: .black.opacity(0.25), radius: 4, x: 0, y: 2)

                            Text("Вы")
                                .font(.caption2.bold())
                                .foregroundStyle(.white)
                                .padding(.horizontal, 8)
                                .padding(.vertical, 4)
                                .background(.black.opacity(0.45), in: Capsule())
                        }
                    }
                }

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

            Button {
                Task { await refreshNearby() }
            } label: {
                HStack {
                    if isRefreshing {
                        ProgressView()
                    }
                    Label("Обновить прогулку", systemImage: "location.fill")
                }
                .font(.headline)
                .frame(maxWidth: .infinity)
                .frame(height: 46)
            }
            .buttonStyle(.borderedProminent)
            .tint(.orange)
            .disabled(isRefreshing)

            if let location = userLocation.location {
                HStack(spacing: 8) {
                    Image(systemName: "location.fill")
                    Text(locationText(for: location.coordinate))
                }
                .font(.caption)
                .foregroundStyle(.white.opacity(0.75))
            } else {
                HStack(spacing: 8) {
                    Image(systemName: "location.slash")
                    Text(userLocation.errorMessage ?? "Ждем геолокацию")
                }
                .font(.caption)
                .foregroundStyle(.white.opacity(0.75))
            }
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

    private func refreshNearby() async {
        guard let location = userLocation.location else {
            actionErrorMessage =
                userLocation.errorMessage ??
                "Геолокация ещё не определена"
            return
        }

        if appViewModel.walkers.isEmpty {
            screenState = .loading(
                message: "Ищем прогулки рядом..."
            )
        }

        isRefreshing = true
        actionErrorMessage = nil

        defer {
            isRefreshing = false
        }

        let coordinate = location.coordinate

        do {
            try await appViewModel.shareLocationAndLoadNearby(
                latitude: coordinate.latitude,
                longitude: coordinate.longitude
            )

            centerMap(on: coordinate)
            screenState = .content
        } catch {
            if appViewModel.walkers.isEmpty {
                screenState = .from(error)
            } else {
                screenState = .content
                actionErrorMessage = error.localizedDescription
            }
        }
    }

    private func centerMap(on coordinate: CLLocationCoordinate2D) {
        position = .region(
            MKCoordinateRegion(
                center: coordinate,
                span: MKCoordinateSpan(latitudeDelta: 0.03, longitudeDelta: 0.03)
            )
        )
    }

    private func locationText(for coordinate: CLLocationCoordinate2D) -> String {
        String(format: "Вы: %.5f, %.5f", coordinate.latitude, coordinate.longitude)
    }
}
