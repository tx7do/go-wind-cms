import 'dart:convert' show json;

class UserPreference {
  UserPreference({
    this.language = 'zh_CN',
    this.themeMode = 'system',  // 默认跟随系统
    this.fontSize = 12,
    this.allowNotifications = true,
    this.soundEnabled = true,
    this.vibrationEnabled = true,
    this.badgeEnabled = true,
    this.anonymousMode = true,
    this.shareAnalytics = true,
    this.dataCollection = true,
    this.locationSharing = true,
    this.autoPlayMedia = true,
    this.showPreviews = true,
    this.useBiometric = true,
    this.hapticFeedback = true,
    this.rememberLogin = true,
  });

  String language;
  String themeMode;
  int fontSize;
  bool allowNotifications;
  bool soundEnabled;
  bool vibrationEnabled;
  bool badgeEnabled;
  bool anonymousMode;
  bool shareAnalytics;
  bool dataCollection;
  bool locationSharing;
  bool autoPlayMedia;
  bool showPreviews;
  bool useBiometric;
  bool hapticFeedback;
  bool rememberLogin;

  factory UserPreference.fromJson(dynamic source) {
    final Map<String, dynamic> map = source is String
        ? (json.decode(source) as Map<String, dynamic>)
        : (source as Map<String, dynamic>);

    return UserPreference(
      language: map['language'] as String? ?? 'zh_CN',
      themeMode: map['themeMode'] as String? ?? 'system',  // 默认跟随系统
      fontSize: (map['fontSize'] as num?)?.toInt() ?? 12,
      allowNotifications: map['allowNotifications'] as bool? ?? true,
      soundEnabled: map['soundEnabled'] as bool? ?? true,
      vibrationEnabled: map['vibrationEnabled'] as bool? ?? true,
      badgeEnabled: map['badgeEnabled'] as bool? ?? true,
      anonymousMode: map['anonymousMode'] as bool? ?? true,
      shareAnalytics: map['shareAnalytics'] as bool? ?? true,
      dataCollection: map['dataCollection'] as bool? ?? true,
      locationSharing: map['locationSharing'] as bool? ?? true,
      autoPlayMedia: map['autoPlayMedia'] as bool? ?? true,
      showPreviews: map['showPreviews'] as bool? ?? true,
      useBiometric: map['useBiometric'] as bool? ?? true,
      hapticFeedback: map['hapticFeedback'] as bool? ?? true,
      rememberLogin: map['rememberLogin'] as bool? ?? true,
    );
  }

  Map<String, dynamic> toProto3Json() {
    return <String, dynamic>{
      'language': language,
      'themeMode': themeMode,
      'fontSize': fontSize,
      'allowNotifications': allowNotifications,
      'soundEnabled': soundEnabled,
      'vibrationEnabled': vibrationEnabled,
      'badgeEnabled': badgeEnabled,
      'anonymousMode': anonymousMode,
      'shareAnalytics': shareAnalytics,
      'dataCollection': dataCollection,
      'locationSharing': locationSharing,
      'autoPlayMedia': autoPlayMedia,
      'showPreviews': showPreviews,
      'useBiometric': useBiometric,
      'hapticFeedback': hapticFeedback,
      'rememberLogin': rememberLogin,
    };
  }
}
