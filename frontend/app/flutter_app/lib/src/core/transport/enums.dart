/// 连接状态
enum ConnectionState {
  /// 断开连接中
  disconnecting,

  /// 已断开连接
  disconnected,

  /// 连接中
  connecting,

  /// 已连接
  connected,

  /// 发生错误
  faulted
}
