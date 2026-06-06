import 'package:jiffy/jiffy.dart';
// import 'package:ntp/ntp.dart';

// import 'package:flutter_app/src/core/constants/index.dart';

/// 日期时间工具集
class DateTimeUtils {
  DateTimeUtils._(); // Private constructor to prevent instantiation

  static int ntpOffset = 0;

  static init() async {
    // NTP对时
    // DateTimeUtil.refreshNTPOffset();

    // 设置日期的本地化
    // await Jiffy.setLocale(LocalStorage.currentLanguage);
  }

  /// 刷新NTF时间偏移量
  static refreshNTPOffset() async {
    // final int offset = await NTP.getNtpOffset(localTime: DateTime.now(), lookUpAddress: Environments.ntpHost);
    // ntpOffset = offset;
  }

  /// 校准过的当前时间
  static DateTime now({bool forceSync = false}) {
    if (forceSync) {
      refreshNTPOffset();
    }
    DateTime internetTime = DateTime.now().add(Duration(milliseconds: ntpOffset));
    return internetTime;
  }

  /// 当前时间的毫秒时间戳
  static int currentTimeMillis() {
    return now().millisecondsSinceEpoch;
  }

  /// 当前utc时间戳
  static int utc() {
    return now().toUtc().millisecondsSinceEpoch;
  }

  /// UTC时间的秒时间戳
  static int utcSecond() {
    return utc() ~/ 1000;
  }

  static String formatLastTime(DateTime dateTime) {
    int diff = Jiffy.now().diff(
      Jiffy.parseFromDateTime(dateTime),
      unit: Unit.day,
    ) as int;
    if (diff > 6) {
      // 2022-01-22
      return Jiffy.parseFromDateTime(dateTime).format(pattern: 'y-MM-dd');
    } else if (diff > 2) {
      // 星期二 09:18
      return Jiffy.parseFromDateTime(dateTime).format(pattern: 'EEEE HH:mm');
    } else {
      return Jiffy.parseFromDateTime(dateTime).startOf(Unit.minute).fromNow();
    }
  }

  static String formatCustomDateHeader(DateTime dt) {
    int diff = Jiffy.now().diff(
      Jiffy.parseFromDateTime(dt),
      unit: Unit.day,
    ) as int;
    if (diff > 6) {
      // 2022-01-22 11:58
      return Jiffy.parseFromDateTime(dt).format(pattern: 'y-MM-dd HH:mm');
    } else if (diff > 2) {
      // 星期二 09:18
      return Jiffy.parseFromDateTime(dt).format(pattern: 'EEEE HH:mm');
    } else {
      return Jiffy.parseFromDateTime(dt).startOf(Unit.minute).fromNow();
    }
  }
}
