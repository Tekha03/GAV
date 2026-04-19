import CoreLocation
import Foundation

@available(macOS 12.0, *)
class UserLocation: NSObject, ObservableObject, CLLocationManagerDelegate {
    @Published var location: CLLocation?
    @Published var locationStatus: LocationStatus = .inactive
    @Published var locationVisibility: LocationVisibility = .everyone
    @Published var isLocationEnabled: Bool = false

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

    func locationManager(_ manager: CLLocationManager, didUpdateLocations locations: [CLLocation]) {
        guard let location = locations.last else { return }
        DispatchQueue.main.async {
            self.location = location
        }
    }

    func updateLocationStatus(_ status: LocationStatus) {
        locationStatus = status
    }

    func updateLocationVisibility(_ visibility: LocationVisibility) {
        locationVisibility = visibility
    }
}