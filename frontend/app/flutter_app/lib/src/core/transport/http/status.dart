/// RPC返回状态
class Status {
  final int? code;
  final String? reason;
  final String? message;
  final Map<String, String>? metadata;

  get getMessage => message ?? '';

  get getReason => reason ?? '';

  Status({
    this.code,
    this.reason,
    this.message,
    this.metadata,
  });

  Status.fromJson(Map<String, dynamic> json)
      : code = (json['code'] != null) ? json['code'] : null,
        reason = (json['reason'] != null) ? json['reason'] : null,
        message = (json['message'] != null) ? json['message'] : null,
        metadata =
            (json['metadata'] != null) ? Map.from(json['metadata']) : null;
}
