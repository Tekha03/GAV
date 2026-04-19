import SwiftUI
import MapKit

@available(macOS 12.0, *)
struct MapView: View {
    @StateObject private var viewModel: MapViewModel

    var body: some View {
        NavigationStack {
            Map(
                coordinateRegion: $userLocation.map { .(
                    MapCoordinateRegion(
                        center: CLLocationCoordinate2D(
                            latitude: $0.coordinate.latitude,
                            longitude: $0.coordinate.longitude
                        ),
                        span: .(
                            latitudeDelta: 0.01,
                            longitudeDelta: 0.01
                        )
                    )
                ) }
            )
            .annotation(
                .dot,
                coordinate: .(
                    CLLocationCoordinate2D(
                        latitude: viewModel.userLocation?.location?.coordinate.latitude ?? 0,
                        longitude: viewModel.userLocation?.location?.coordinate.longitude ?? 0
                    )
                )
            )

            if viewModel.isLoading {
                ProgressView()
            }
        }
    }
}