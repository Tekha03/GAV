import CoreLocation
import Foundation
import Combine

@MainActor
final class UserLocation: NSObject, ObservableObject, CLLocationManagerDelegate {

    @Published var location: CLLocation?
    @Published var locationStatus: LocationStatus = .inactive
    @Published var locationVisibility: LocationVisibility = .everyone
    @Published var isLocationEnabled = false
    @Published var authorizationStatus: CLAuthorizationStatus = .notDetermined
    @Published var errorMessage: String?

    private let locationManager = CLLocationManager()

    override init() {
        super.init()
        locationManager.delegate = self
        locationManager.desiredAccuracy = kCLLocationAccuracyHundredMeters
        authorizationStatus = locationManager.authorizationStatus
    }

    func startUpdatingLocation() {
        isLocationEnabled = true
        errorMessage = nil

        if locationManager.authorizationStatus == .notDetermined {
            locationManager.requestWhenInUseAuthorization()
            return
        }

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
            self.errorMessage = nil
        }
    }

    nonisolated func locationManagerDidChangeAuthorization(_ manager: CLLocationManager) {
        let status = manager.authorizationStatus

        Task { @MainActor in
            self.authorizationStatus = status

            switch status {
            case .authorizedAlways, .authorizedWhenInUse:
                self.errorMessage = nil
                self.locationManager.startUpdatingLocation()
            case .denied, .restricted:
                self.errorMessage = "Геолокация выключена для приложения"
            case .notDetermined:
                break
            @unknown default:
                break
            }
        }
    }

    nonisolated func locationManager(_ manager: CLLocationManager, didFailWithError error: Error) {
        Task { @MainActor in
            self.errorMessage = error.localizedDescription
        }
    }
}
