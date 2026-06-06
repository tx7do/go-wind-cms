import 'package:uuid/uuid.dart';
import 'package:xid/xid.dart';

/// UUID工具类
class UuidUtils {
  UuidUtils._(); // Private constructor to prevent instantiation

  /// Mongo Object ID
  static String get xid {
    return Xid().toString();
  }

  /// UUID V1
  static String get v1 {
    return const Uuid().v1();
  }

  static String get v1WithoutHyphens {
    return const Uuid().v1().replaceAll('-', '');
  }

  /// UUID V4
  static String get v4 {
    return const Uuid().v4();
  }

  static String get v4WithoutHyphens {
    return const Uuid().v4().replaceAll('-', '');
  }

  /// UUID V5
  static String v5(String? namespace, String? name) {
    return const Uuid().v5(namespace, name);
  }

  static String v5WithoutHyphens(String? namespace, String? name) {
    return const Uuid().v5(namespace, name).replaceAll('-', '');
  }

  /// UUID V6
  static String get v6 {
    return const Uuid().v6();
  }

  static String get v6WithoutHyphens {
    return const Uuid().v6().replaceAll('-', '');
  }

  /// UUID V7
  static String get v7 {
    return const Uuid().v7();
  }

  static String get v7WithoutHyphens {
    return const Uuid().v7().replaceAll('-', '');
  }

  /// UUID V8
  static String get v8 {
    return const Uuid().v8();
  }

  static String get v8WithoutHyphens {
    return const Uuid().v8().replaceAll('-', '');
  }
}
