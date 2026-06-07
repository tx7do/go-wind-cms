import 'dart:ui';

import 'package:flutter_app/generated/l10n.dart';
import 'package:language_code/language_code.dart';

class LanguageInfo {
  String code;
  String nativeName;

  LanguageInfo({required this.code, required this.nativeName});
}

class LanguageListRepository {
  List<Locale> get languageList => S.delegate.supportedLocales;

  List<LanguageInfo> getLanguageList() {
    List<LanguageInfo> infoList = [];
    for (int i = 0; i < languageList.length; i++) {
      final locale = languageList[i];
      final codes = localeToNativeName(locale);
      infoList.add(
        LanguageInfo(code: codes.code, nativeName: codes.nativeName),
      );
    }
    return infoList;
  }

  /// 获取语言的原生名称
  static LanguageCodes localeToNativeName(Locale locale) {
    final codes = LanguageCodes.fromLocale(locale);
    // logger.d('localeToNativeName ${codes.code}');
    return codes;
  }
}
