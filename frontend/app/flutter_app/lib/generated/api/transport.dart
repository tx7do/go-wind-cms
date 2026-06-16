// Shared transport infrastructure for protoc-gen-dart-http.
// Auto-generated. DO NOT EDIT.

import 'dart:async';

/// Metadata for an RPC call, passed to the transport for routing and diagnostics.
class TransportMeta {
  final String service;
  final String method;
  const TransportMeta({required this.service, required this.method});

  @override
  String toString() => 'TransportMeta(service: $service, method: $method)';

  @override
  bool operator ==(Object other) =>
    identical(this, other) ||
    other is TransportMeta &&
      runtimeType == other.runtimeType &&
      service == other.service &&
      method == other.method;

  @override
  int get hashCode => Object.hash(service, method);
}

/// Abstract transport interface for making HTTP requests.
///
/// Implement this with your preferred HTTP client (package:http, dio, etc.).
abstract class ClientTransport {
  /// Performs a unary (request/response) RPC.
  Future<dynamic> unary(
    String path,
    String method,
    String? body,
    TransportMeta meta, {
    Map<String, String>? headers,
  });

  /// Opens a server-streaming connection (e.g. SSE).
  /// Returns a stream of JSON-decoded event payloads.
  Stream<Map<String, dynamic>> serverStream(
    String path,
    TransportMeta meta, {
    Map<String, String>? headers,
  });

  /// Opens a bidirectional streaming connection (e.g. WebSocket).
  DuplexConnection duplexStream(
    String path,
    TransportMeta meta, {
    Map<String, String>? headers,
  });
}

/// Safely joins a base URL with a request path.
///
/// Handles all combinations of trailing/leading slashes so that
/// the result always contains exactly one slash between base and path:
///
/// ```dart
/// joinPath('https://api.example.com', '/v1/users');   // https://api.example.com/v1/users
/// joinPath('https://api.example.com/', '/v1/users'); // https://api.example.com/v1/users
/// joinPath('https://api.example.com/', 'v1/users');   // https://api.example.com/v1/users
/// joinPath('https://api.example.com', 'v1/users');    // https://api.example.com/v1/users
/// ```
///
/// Use this in your [ClientTransport] implementation to build the full URL.
String joinPath(String baseUrl, String path) {
  if (path.isEmpty) return baseUrl;
  final base = baseUrl.endsWith('/') ? baseUrl.substring(0, baseUrl.length - 1) : baseUrl;
  final p = path.startsWith('/') ? path : '/$path';
  return base + p;
}

/// Abstract bidirectional connection for duplex streaming.
abstract class DuplexConnection {
  /// Stream of incoming JSON messages from the server.
  Stream<Map<String, dynamic>> get incoming;
  /// Sends a JSON message to the server.
  void send(Map<String, dynamic> data);
  /// Closes the connection and releases resources.
  Future<void> close();
}

/// Type-safe wrapper around [DuplexConnection] that handles JSON (de)serialization.
class TypedDuplexConnection<TIn, TOut> {
  final DuplexConnection _conn;
  final TOut Function(Map<String, dynamic>) _fromJson;
  final Map<String, dynamic> Function(TIn) _toJson;

  TypedDuplexConnection(this._conn, this._fromJson, this._toJson);

  /// Typed stream of incoming messages.
  Stream<TOut> get stream => _conn.incoming.map(_fromJson);

  /// Sends a typed message, serialized to JSON.
  void send(TIn data) => _conn.send(_toJson(data));

  /// Closes the connection and releases resources.
  Future<void> close() => _conn.close();
}
