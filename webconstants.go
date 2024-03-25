package gweb

type MessageType int

const TextMessage MessageType = 1
const BinaryMessage MessageType = 2

const NotFound = "Not Found"
const InternalServerError = "Internal Server Error"
const InvalidData = "Invalid data"
const Authorization = "Authorization"
const MsgInvalidToken = "Invalid Token"
const MsgExpiredToken = "Expired Token"
const WebNotInitialized = "Web Instance not initialized yet"
const InvalidWebGroup = "Invalid web group"
const InvalidPath = "Invalid path, mising /"
const NoWebSocket = "No active websocket connection"
