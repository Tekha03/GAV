import Foundation

struct APIErrorCode: RawRepresentable, Codable, Hashable, Sendable {
    let rawValue: String

    init(rawValue: String) {
        self.rawValue = rawValue
    }

    static let authCredentialsInvalid = Self(rawValue: "AUTH_CREDENTIALS_INVALID")
    static let authTokenMissing = Self(rawValue: "AUTH_TOKEN_MISSING")
    static let authTokenInvalid = Self(rawValue: "AUTH_TOKEN_INVALID")
    static let authTokenExpired = Self(rawValue: "AUTH_TOKEN_EXPIRED")
    static let authRefreshInvalid = Self(rawValue: "AUTH_REFRESH_TOKEN_INVALID")
    static let authForbidden = Self(rawValue: "AUTH_FORBIDDEN")
    static let authEmailExists = Self(rawValue: "AUTH_EMAIL_ALREADY_EXISTS")
    static let profileNotFound = Self(rawValue: "PROFILE_NOT_FOUND")
    static let profileAlreadyExists = Self(rawValue: "PROFILE_ALREADY_EXISTS")
    static let mediaFileTooLarge = Self(rawValue: "MEDIA_FILE_TOO_LARGE")
    static let mediaTypeInvalid = Self(rawValue: "MEDIA_TYPE_INVALID")
    static let validation = Self(rawValue: "VALIDATION_ERROR")
    static let serviceUnavailable = Self(rawValue: "SERVICE_UNAVAILABLE")
    static let requestTimeout = Self(rawValue: "REQUEST_TIMEOUT")
    static let internalError = Self(rawValue: "INTERNAL_ERROR")
}

struct APIErrorCategory: RawRepresentable, Codable, Hashable, Sendable {
    let rawValue: String

    init(rawValue: String) {
        self.rawValue = rawValue
    }

    static let validation = Self(rawValue: "validation")
    static let unauthenticated = Self(rawValue: "unauthenticated")
    static let permissionDenied = Self(rawValue: "permission_denied")
    static let notFound = Self(rawValue: "not_found")
    static let conflict = Self(rawValue: "conflict")
    static let unavailable = Self(rawValue: "unavailable")
    static let unsupported = Self(rawValue: "unsupported")
    static let cancelled = Self(rawValue: "cancelled")
    static let internalError = Self(rawValue: "internal")
}

enum JSONValue: Codable, Equatable, Sendable {
    case string(String)
    case number(Double)
    case bool(Bool)
    case object([String: JSONValue])
    case array([JSONValue])
    case null

    init(from decoder: Decoder) throws {
        let container = try decoder.singleValueContainer()
        if container.decodeNil() {
            self = .null
        } else if let value = try? container.decode(Bool.self) {
            self = .bool(value)
        } else if let value = try? container.decode(Double.self) {
            self = .number(value)
        } else if let value = try? container.decode(String.self) {
            self = .string(value)
        } else if let value = try? container.decode([String: JSONValue].self) {
            self = .object(value)
        } else if let value = try? container.decode([JSONValue].self) {
            self = .array(value)
        } else {
            throw DecodingError.dataCorruptedError(
                in: container,
                debugDescription: "Unsupported JSON value"
            )
        }
    }

    func encode(to encoder: Encoder) throws {
        var container = encoder.singleValueContainer()
        switch self {
        case .string(let value): try container.encode(value)
        case .number(let value): try container.encode(value)
        case .bool(let value): try container.encode(value)
        case .object(let value): try container.encode(value)
        case .array(let value): try container.encode(value)
        case .null: try container.encodeNil()
        }
    }
}

struct APIErrorBody: Codable, Equatable, Sendable {
    let code: APIErrorCode
    let category: APIErrorCategory
    let message: String
    let details: [String: JSONValue]?
}

struct APIErrorResponse: Codable, Equatable, Sendable {
    let error: APIErrorBody
}

enum APIError: Error, LocalizedError {
    case invalidURL
    case server(statusCode: Int, body: APIErrorBody)
    case invalidResponse(statusCode: Int)
    case decodingError(Error)
    case networkError(Error)

    var code: APIErrorCode? {
        guard case .server(_, let body) = self else { return nil }
        return body.code
    }

    var category: APIErrorCategory? {
        guard case .server(_, let body) = self else { return nil }
        return body.category
    }

    var errorDescription: String? {
        switch self {
        case .invalidURL:
            return "Не удалось сформировать адрес запроса"
        case .server(_, let body):
            return Self.userMessage(for: body)
        case .invalidResponse:
            return "Сервер вернул некорректный ответ"
        case .decodingError:
            return "Не удалось обработать ответ сервера"
        case .networkError(let error):
            if let urlError = error as? URLError, urlError.code == .notConnectedToInternet {
                return "Нет подключения к интернету"
            }
            return "Не удалось подключиться к серверу"
        }
    }

    private static func userMessage(for body: APIErrorBody) -> String {
        switch body.code {
        case .authCredentialsInvalid:
            return "Неверная почта или пароль"
        case .authTokenMissing, .authTokenInvalid, .authTokenExpired, .authRefreshInvalid:
            return "Сессия истекла. Войдите снова"
        case .authForbidden:
            return "Недостаточно прав для этого действия"
        case .authEmailExists:
            return "Пользователь с такой почтой уже существует"
        case .profileAlreadyExists:
            return "Этот никнейм уже занят"
        case .profileNotFound:
            return "Профиль не найден"
        case .mediaFileTooLarge:
            return "Размер файла превышает допустимый"
        case .mediaTypeInvalid:
            return "Этот формат файла не поддерживается"
        case .requestTimeout:
            return "Сервер не успел ответить. Попробуйте ещё раз"
        case .serviceUnavailable:
            return "Сервис временно недоступен"
        default:
            return fallbackMessage(for: body.category)
        }
    }

    private static func fallbackMessage(for category: APIErrorCategory) -> String {
        switch category {
        case .validation: return "Проверьте введённые данные"
        case .unauthenticated: return "Необходимо войти в аккаунт"
        case .permissionDenied: return "Недостаточно прав для этого действия"
        case .notFound: return "Запрошенные данные не найдены"
        case .conflict: return "Данные уже существуют или были изменены"
        case .unavailable: return "Сервис временно недоступен"
        case .unsupported: return "Операция не поддерживается"
        case .cancelled: return "Операция отменена"
        default: return "Что-то пошло не так"
        }
    }
}
