import 'package:jose/jose.dart';

import 'package:flutter_app/src/core/utilities/logger.dart' show fatal;

/// JWT 令牌工具类
class JwtUtils {
  JwtUtils._(); // Private constructor to prevent instantiation

  /// 验证token是否过期
  static bool isTokenExpired(String? token) {
    try {
      // 直接解析JWT（不验证签名）
      var jwt = JsonWebToken.unverified(token ?? '');

      // 没有exp字段，视为未过期
      if (jwt.claims.expiry == null) {
        return false;
      }

      return jwt.claims.expiry!.isBefore(DateTime.now());
    } catch (e) {
      fatal('parse jwt payload failed: $e');
      return true;
    }
  }
}
