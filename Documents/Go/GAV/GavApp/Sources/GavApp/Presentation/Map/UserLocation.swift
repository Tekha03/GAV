import CoreLocation
import Foundation
import Combine

@MainActor
final class UserLocation: NSObject, ObservableObject, CLLocationManagerDelegate {

    @Published var location: CLLocation?
    @Published var locationStatus: LocationStatus = .inactive
    @Published var locationVisibility: LocationVisibility = .everyone
    @Published var isLocationEnabled = false

    private let locationManager = CLLocationManager()

    override init() {
        super.init()
        locationManager.delegate = self
        locationManager.requestAlwaysAuthorization()
    }

    func startUpdatingLocation() {
        isLocationEnabled = true
        locationManager.startUpdatingLocation()
    }

    func stopUpdatingLocation() {
        isLocationEnabled = false
        locationManager.stopUpdatingLocation()
    }

    nonisolated func locationManager(
        _ manager: CLLocationManager,
        didUpdateLocations locations: [CLLocation]
    ) {
        guard let location = locations.last else { return }

        Task { @MainActor in
            self.location = location
        }
    }
}
