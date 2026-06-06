import 'dart:convert';
import 'dart:io';
import 'dart:math' as math;

import 'package:flutter/widgets.dart';
import 'package:intl/intl.dart';

/// 字符串工具类
class StringUtils {
  StringUtils._(); // Private constructor to prevent instantiation

  static String formatDouble3(double value) {
    return NumberFormat('0.000').format(value);
  }

  static String getCurrencySymbol(BuildContext context) {
    Locale locale = Localizations.localeOf(context);
    return NumberFormat.simpleCurrency(
      locale: locale.toString(),
    ).currencySymbol;
  }

  /// 生成随机字符串
  static String getRandomString(int length) {
    const chars =
        'AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz1234567890';
    math.Random rnd = math.Random();
    return String.fromCharCodes(
      Iterable.generate(
        length,
        (_) => chars.codeUnitAt(rnd.nextInt(chars.length)),
      ),
    );
  }

  ///
  /// 返回给定字符串，如果给定字符串为空则返回默认字符串
  ///
  static String defaultString(String? str, {String defaultStr = ''}) {
    return str ?? defaultStr;
  }

  ///
  /// 检查给定的字符串[s]是空还是空
  ///
  static bool isNullOrEmpty(String? s) =>
      (s == null || s.isEmpty) ? true : false;

  ///
  /// 检查给定的字符串[s]是否为null或空
  ///
  static bool isNotNullOrEmpty(String? s) => !isNullOrEmpty(s);

  /// 字符串是否为空
  static bool isStringEmpty(String? val) {
    if (val == null) {
      return true;
    }
    return val.trim().isEmpty;
  }

  /// 是否字符串不为空
  static bool isStringNotEmpty(String? val) {
    return !StringUtils.isStringEmpty(val);
  }

  ///
  /// 将给定的字符串[s]从 大写大小写(camelCase) 转换为 大写下划线(upperCaseUnderscore)
  /// 示例 : helloWorld => HELLO_WORLD
  ///
  static String camelCaseToUpperUnderscore(String s) {
    var sb = StringBuffer();
    var first = true;
    for (var rune in s.runes) {
      var char = String.fromCharCode(rune);
      if (isUpperCase(char) && !first) {
        sb.write('_');
        sb.write(char.toUpperCase());
      } else {
        first = false;
        sb.write(char.toUpperCase());
      }
    }
    return sb.toString();
  }

  ///
  /// 将给定的字符串[s]从 大写大小写(camelCase) 转换为 小写下划线(lowerCaseUnderscore)
  /// 示例 : helloWorld => hello_world
  ///
  static String camelCaseToLowerUnderscore(String s) {
    var sb = StringBuffer();
    var first = true;
    for (var rune in s.runes) {
      var char = String.fromCharCode(rune);
      if (isUpperCase(char) && !first) {
        if (char != '_') {
          sb.write('_');
        }
        sb.write(char.toLowerCase());
      } else {
        first = false;
        sb.write(char.toLowerCase());
      }
    }
    return sb.toString();
  }

  ///
  /// 检查给定字符串[s]是否为小写
  ///
  static bool isLowerCase(String s) {
    return s == s.toLowerCase();
  }

  ///
  /// 检查给定字符串[s]是否为大写
  ///
  static bool isUpperCase(String s) {
    return s == s.toUpperCase();
  }

  ///
  /// 检查给定字符串[s]是否只包含ascii字符
  ///
  static bool isAscii(String s) {
    try {
      asciiCodec.decode(s.codeUnits);
    } catch (e) {
      return false;
    }
    return true;
  }

  ///
  /// 将给定字符串[s]大写。如果[allWords]设置为true，则将给定字符串[s]中的所有单词大写。
  ///
  /// 因此，字符串[s]被“”(空格)分隔。
  ///
  /// 示例 :
  ///
  /// * [s] = "world" => World
  /// * [s] = "WORLD" => World
  /// * [s] = "the quick lazy fox"  => The quick lazy fox
  /// * [s] = "the quick lazy fox" and [allWords] = true => The Quick Lazy Fox
  ///
  static String capitalize(String s, {bool allWords = false}) {
    if (s.isEmpty) {
      return '';
    }
    s = s.trim();
    if (allWords) {
      var words = s.split(' ');
      var capitalized = [];
      for (var w in words) {
        capitalized.add(capitalize(w));
      }
      return capitalized.join(' ');
    } else {
      return s.substring(0, 1).toUpperCase() + s.substring(1).toLowerCase();
    }
  }

  ///
  /// 反转给定的字符串[s]
  /// 示例 : hello => olleh
  ///
  static String reverse(String s) {
    return String.fromCharCodes(s.runes.toList().reversed);
  }

  ///
  /// 计算给定[char]在给定字符串[s]中出现的频率。
  /// 值[caseSensitive]控制是否只查找给定的[char]
  /// 或者同样的大小写版本。
  ///
  /// 示例: Hello and char l => 2
  ///
  static int countChars(String s, String char, {bool caseSensitive = true}) {
    var count = 0;
    s.codeUnits.toList().forEach((i) {
      if (caseSensitive) {
        if (i == char.runes.first) {
          count++;
        }
      } else {
        if (i == char.toLowerCase().runes.first ||
            i == char.toUpperCase().runes.first) {
          count++;
        }
      }
    });
    return count;
  }

  ///
  /// 检查给定字符串[s]是否为数字。
  ///
  /// 如果给定的字符串[s]为空，则返回false。
  ///
  static bool isDigit(String s) {
    if (s.isEmpty) {
      return false;
    }
    if (s.length > 1) {
      for (var r in s.runes) {
        if (r ^ 0x30 > 9) {
          return false;
        }
      }
      return true;
    } else {
      return s.runes.first ^ 0x30 <= 9;
    }
  }

  ///
  /// 比较给定字符串[a]和[b]。
  ///
  static bool equalsIgnoreCase(String a, String b) =>
      a.toLowerCase() == b.toLowerCase();

  ///
  /// 检查给定的[list]是否包含字符串[s]
  ///
  static bool inList(String s, List<String> list, {bool ignoreCase = false}) {
    for (var l in list) {
      if (ignoreCase) {
        if (equalsIgnoreCase(s, l)) {
          return true;
        }
      } else {
        if (s == l) {
          return true;
        }
      }
    }
    return false;
  }

  ///
  /// 检查给定字符串[s]是否为回文
  /// 示例:
  /// aha => true
  /// hello => false
  ///
  static bool isPalindrome(String s) {
    for (var i = 0; i < s.length / 2; i++) {
      if (s[i] != s[s.length - 1 - i]) return false;
    }
    return true;
  }

  ///
  /// Replaces chars of the given String [s] with [replace].
  ///
  /// The default value of [replace] is *.
  /// [begin] determines the start of the 'replacing'. If [begin] is null, it starts from index 0.
  /// [end] defines the end of the 'replacing'. If [end] is null, it ends at [s] length divided by 2.
  /// If [s] is empty or consists of only 1 char, the method returns null.
  ///
  /// Example :
  /// 1234567890 => *****67890
  /// 1234567890 with begin 2 and end 6 => 12****7890
  /// 1234567890 with begin 1 => 1****67890
  ///
  static String? hidePartial(
    String s, {
    int begin = 0,
    int? end,
    String replace = '*',
  }) {
    var buffer = StringBuffer();
    if (s.length <= 1) {
      return null;
    }
    if (end == null) {
      end = (s.length / 2).round();
    } else {
      if (end > s.length) {
        end = s.length;
      }
    }
    for (var i = 0; i < s.length; i++) {
      if (i >= end) {
        buffer.write(String.fromCharCode(s.runes.elementAt(i)));
        continue;
      }
      if (i >= begin) {
        buffer.write(replace);
        continue;
      }
      buffer.write(String.fromCharCode(s.runes.elementAt(i)));
    }
    return buffer.toString();
  }

  ///
  /// Add a [char] at a [position] with the given String [s].
  ///
  /// The boolean [repeat] defines whether to add the [char] at every [position].
  /// If [position] is greater than the length of [s], it will return [s].
  /// If [repeat] is true and [position] is 0, it will return [s].
  ///
  /// Example :
  /// 1234567890 , '-', 3 => 123-4567890
  /// 1234567890 , '-', 3, true => 123-456-789-0
  ///
  static String addCharAtPosition(
    String s,
    String char,
    int position, {
    bool repeat = false,
  }) {
    if (!repeat) {
      if (s.length < position) {
        return s;
      }
      var before = s.substring(0, position);
      var after = s.substring(position, s.length);
      return before + char + after;
    } else {
      if (position == 0) {
        return s;
      }
      var buffer = StringBuffer();
      for (var i = 0; i < s.length; i++) {
        if (i != 0 && i % position == 0) {
          buffer.write(char);
        }
        buffer.write(String.fromCharCode(s.runes.elementAt(i)));
      }
      return buffer.toString();
    }
  }

  ///
  /// Splits the given String [s] in chunks with the given [chunkSize].
  ///
  static List<String> chunk(String s, int chunkSize) {
    var chunked = <String>[];
    for (var i = 0; i < s.length; i += chunkSize) {
      var end = (i + chunkSize < s.length) ? i + chunkSize : s.length;
      chunked.add(s.substring(i, end));
    }
    return chunked;
  }

  ///
  /// Picks only required string[value] starting [from] and ending at [to]
  ///
  /// Example :
  /// pickOnly('123456789',from:3,to:7);
  /// returns '34567'
  ///
  static String pickOnly(value, {int from = 1, int to = -1}) {
    try {
      return value.substring(
        from == 0 ? 0 : from - 1,
        to == -1 ? value.length : to,
      );
    } catch (e) {
      return value;
    }
  }

  ///
  /// Removes character with [index] from a String [value]
  ///
  /// Example:
  /// removeCharAtPosition('flutterr', 8);
  /// returns 'flutter'
  static String removeCharAtPosition(String value, int index) {
    try {
      return value.substring(0, -1 + index) +
          value.substring(index, value.length);
    } catch (e) {
      return value;
    }
  }

  ///
  ///Remove String[value] with [pattern]
  ///
  ///[repeat]:boolean => if(true) removes all occurence
  ///
  ///[casensitive]:boolean => if(true) a != A
  ///
  ///Example: removeExp('Hello This World', 'This'); returns 'Hello World'
  ///
  static String removeExp(
    String value,
    String pattern, {
    bool repeat = true,
    bool caseSensitive = true,
    bool multiLine = false,
    bool dotAll = false,
    bool unicode = false,
  }) {
    var result = value;
    if (repeat) {
      result = value
          .replaceAll(
            RegExp(
              pattern,
              caseSensitive: caseSensitive,
              multiLine: multiLine,
              dotAll: dotAll,
              unicode: unicode,
            ),
            '',
          )
          .replaceAll(RegExp(' +'), ' ')
          .trim();
    } else {
      result = value
          .replaceFirst(
            RegExp(
              pattern,
              caseSensitive: caseSensitive,
              multiLine: multiLine,
              dotAll: dotAll,
              unicode: unicode,
            ),
            '',
          )
          .replaceAll(RegExp(' +'), ' ')
          .trim();
    }
    return result;
  }

  ///
  /// Takes in a String[value] and truncates it with [length]
  /// [symbol] default is '...'
  ///truncate('This is a Dart Utility Library', 26)
  /// returns 'This is a Dart Utility Lib...'
  static String truncate(String value, int length, {String symbol = '...'}) {
    var result = value;

    try {
      result = value.substring(0, length) + symbol;
    } catch (e) {
      print(e.toString());
    }
    return result;
  }

  ///
  /// Generates a Random string
  ///
  /// * [length] = length of string
  /// * [alphabet] = add alphabet to string
  /// * [uppercase] = adds lowercase alphabet to string
  /// * [lowercase] = adds lowercase alphabet to string
  /// * [numeric] = add integers to string
  /// * [special] = add special characters
  /// * [from] = a string that contains the allowed signs to be used for generating the random string
  ///
  static String generateRandomString(
    int length, {
    alphabet = true,
    numeric = true,
    special = true,
    uppercase = true,
    lowercase = true,
    String from = '',
  }) {
    var res = '';

    do {
      res += _randomizer(
        alphabet,
        numeric,
        lowercase,
        uppercase,
        special,
        from,
      );
    } while (res.length < length);

    var possible = res.split('');
    possible.shuffle();
    var result = [];

    for (var i = 0; i < length; i++) {
      var randomNumber = math.Random().nextInt(length);
      result.add(possible[randomNumber]);
    }

    return result.join();
  }

  static String _randomizer(
    bool alphabet,
    bool numeric,
    bool lowercase,
    bool uppercase,
    bool special,
    String from,
  ) {
    var a = 'ABCDEFGHIJKLMNOPQRXYZ';
    var la = 'abcdefghijklmnopqrxyz';
    var b = '0123456789';
    var c = '~^!@#\$%^&*;`(=?]:[.)_+-|\{}';
    var result = '';

    if (alphabet) {
      if (lowercase) {
        result += la;
      }
      if (uppercase) {
        result += a;
      }

      if (!uppercase && !lowercase) {
        result += a;
        result += la;
      }
    }
    if (numeric) {
      result += b;
    }

    if (special) {
      result += c;
    }

    if (from != '') {
      //if set return it
      result = from;
    }

    return result;
  }

  ///
  /// Converts the given String [s] to PascalCase
  /// Example: your name => YourName
  ///
  static String toPascalCase(String s) {
    final separatedWords = s.split(RegExp(r'[!@#<>?":`~;[\]\\|=+)(*&^%-\s_]+'));
    var newString = '';

    for (final word in separatedWords) {
      newString += word[0].toUpperCase() + word.substring(1).toLowerCase();
    }

    return newString;
  }

  ///
  /// Generates [amount] random strings
  ///
  /// * [length] = length of string
  /// * [alphabet] = add alphabet to string
  /// * [uppercase] = adds lowercase alphabet to string
  /// * [lowercase] = adds lowercase alphabet to string
  /// * [numeric] = add integers to string
  /// * [special] = add special characters
  /// * [from] = a string that contains the allowed signs to be used for generating the random string
  ///
  static List<String> generateRandomStrings(
    int amount,
    int length, {
    alphabet = true,
    numeric = true,
    special = true,
    uppercase = true,
    lowercase = true,
    String from = '',
  }) {
    var l = <String>[];
    for (var i = 0; i < amount; i++) {
      var s = generateRandomString(
        length,
        alphabet: alphabet,
        numeric: numeric,
        special: special,
        uppercase: uppercase,
        lowercase: lowercase,
        from: from,
      );
      l.add(s);
    }
    return l;
  }

  /// Returns text representation of a provided bytes value (e.g. 1kB, 1GB).
  static String formatBytes(
    int size, {
    int fractionDigits = 2,
    int num = 1024,
  }) {
    if (size <= 0) return '0 B';
    final multiple = (math.log(size) / math.log(num)).floor();
    return '${(size / math.pow(num, multiple)).toStringAsFixed(fractionDigits)} ${['B', 'kB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'][multiple]}';
  }

  ///
  /// Checks whether the given String [s] is an IPv4 or IPv6 address.
  ///
  static bool isIP(String s, {InternetAddressType? ipType}) {
    if (ipType == null || ipType == InternetAddressType.any) {
      return isIP(s, ipType: InternetAddressType.IPv4) ||
          isIP(s, ipType: InternetAddressType.IPv6);
    } else if (ipType == InternetAddressType.IPv4) {
      if (!_ipv4Maybe.hasMatch(s)) {
        return false;
      }
      var parts = s.split('.');
      parts.sort((a, b) => int.parse(a) - int.parse(b));
      return int.parse(parts[3]) <= 255;
    } else if (ipType == InternetAddressType.IPv6) {
      return _ipv6.hasMatch(s);
    }
    return false;
  }

  /// 校验身份证
  static bool isValidIdCard(String value) {
    return _idCardRegex.hasMatch(value);
  }

  /// 校验中文
  static bool isChinese(String value) {
    return _chineseRegex.hasMatch(value);
  }

  /// 是否含有的是英文字母和数字
  static bool isAlphanumeric(String input) {
    return _alphanumericRegex.hasMatch(input);
  }

  /// 校验支付宝名称
  static bool isValidAliPayName(String value) {
    return _aliPayRegex.hasMatch(value);
  }

  /// 校验手机号码
  static bool isValidPhone(String? value) {
    if (StringUtils.isStringEmpty(value) || value!.length != 11) {
      return false;
    }
    return _mobilRegex.hasMatch(value);
  }

  /// 验证网页URL
  static bool isValidUrl(String value) {
    return _urlRegex.hasMatch(value);
  }

  /// 校验邮箱地址
  static bool isValidEmail(String? email) {
    if (StringUtils.isStringEmpty(email)) {
      return false;
    }
    return _emailRegex.hasMatch(email!);
  }

  /// 从Email地址中解析出域名来
  static String? extractDomainFromEmail(String email) {
    RegExp regExp = RegExp(r'@[a-zA-Z0-9]+\\\\.([a-zA-Z]+)');
    var match = regExp.firstMatch(email);
    return match != null ? match.group(1) : '';
  }

  /// 查找第一个单词
  static String? findFirstCapitalWord(String text) {
    RegExp regExp = RegExp(r'\\\\b[A-Z][a-z]*\\\\b');
    Match? match = regExp.firstMatch(text);
    return match != null ? match.group(0) : '';
  }

  static AsciiCodec asciiCodec = const AsciiCodec();

  static final RegExp _ipv4Maybe = RegExp(
    r'^(\d?\d?\d)\.(\d?\d?\d)\.(\d?\d?\d)\.(\d?\d?\d)$',
  );

  static final RegExp _ipv6 = RegExp(
    r'^::|^::1|^([a-fA-F0-9]{1,4}::?){1,7}([a-fA-F0-9]{1,4})$',
  );

  static final RegExp _emailRegex = RegExp(
    r"[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?",
    caseSensitive: false,
  );
  static final RegExp _urlRegex = RegExp(
    r"^((https|http|ftp|rtsp|mms)?://)\S+",
    caseSensitive: false,
  );
  static final RegExp _mobilRegex = RegExp(
    r"^(\+\d{1,2}\s)?\(?\d{3}\)?[\s.-]?\d{3}[\s.-]?\d{4}\$",
    caseSensitive: false,
  );
  static final RegExp _alphanumericRegex = RegExp(r'^[a-zA-Z0-9]+$');
  static final RegExp _idCardRegex = RegExp(r"\d{17}[\d|x]|\d{15}");
  static final RegExp _chineseRegex = RegExp(r"[\u4e00-\u9fa5]");
  static final RegExp _aliPayRegex = RegExp(r"[\u4e00-\u9fa5_a-zA-Z]");
}
