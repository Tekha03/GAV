import Foundation

public struct SignInRequestModel: Encodable {
    public let email: String
    public let password: String
}

public struct SignUpRequestModel: Encodable {
    public let email: String
    public let password: String
}