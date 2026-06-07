// GENERATED CODE - DO NOT MODIFY BY HAND
import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import 'intl/messages_all.dart';

// **************************************************************************
// Generator: Flutter Intl IDE plugin
// Made by Localizely
// **************************************************************************

// ignore_for_file: non_constant_identifier_names, lines_longer_than_80_chars
// ignore_for_file: join_return_with_assignment, prefer_final_in_for_each
// ignore_for_file: avoid_redundant_argument_values, avoid_escaping_inner_quotes

class S {
  S();

  static S? _current;

  static S get current {
    assert(
      _current != null,
      'No instance of S was loaded. Try to initialize the S delegate before accessing S.current.',
    );
    return _current!;
  }

  static const AppLocalizationDelegate delegate = AppLocalizationDelegate();

  static Future<S> load(Locale locale) {
    final name = (locale.countryCode?.isEmpty ?? false)
        ? locale.languageCode
        : locale.toString();
    final localeName = Intl.canonicalizedLocale(name);
    return initializeMessages(localeName).then((_) {
      Intl.defaultLocale = localeName;
      final instance = S();
      S._current = instance;

      return instance;
    });
  }

  static S of(BuildContext context) {
    final instance = S.maybeOf(context);
    assert(
      instance != null,
      'No instance of S present in the widget tree. Did you add S.delegate in localizationsDelegates?',
    );
    return instance!;
  }

  static S? maybeOf(BuildContext context) {
    return Localizations.of<S>(context, S);
  }

  /// `GoWind CMS`
  String get appName {
    return Intl.message('GoWind CMS', name: 'appName', desc: '', args: []);
  }

  /// `首页`
  String get home {
    return Intl.message('首页', name: 'home', desc: '', args: []);
  }

  /// `发现`
  String get discover {
    return Intl.message('发现', name: 'discover', desc: '', args: []);
  }

  /// `收藏`
  String get bookmarks {
    return Intl.message('收藏', name: 'bookmarks', desc: '', args: []);
  }

  /// `我的`
  String get me {
    return Intl.message('我的', name: 'me', desc: '', args: []);
  }

  /// `搜索`
  String get search {
    return Intl.message('搜索', name: 'search', desc: '', args: []);
  }

  /// `设置`
  String get settings {
    return Intl.message('设置', name: 'settings', desc: '', args: []);
  }

  /// `登录`
  String get login {
    return Intl.message('登录', name: 'login', desc: '', args: []);
  }

  /// `最新文章`
  String get latestPosts {
    return Intl.message('最新文章', name: 'latestPosts', desc: '', args: []);
  }

  /// `相关文章`
  String get relatedArticles {
    return Intl.message('相关文章', name: 'relatedArticles', desc: '', args: []);
  }

  /// `全部文章`
  String get allPosts {
    return Intl.message('全部文章', name: 'allPosts', desc: '', args: []);
  }

  /// `— 已加载全部 —`
  String get allLoaded {
    return Intl.message('— 已加载全部 —', name: 'allLoaded', desc: '', args: []);
  }

  /// `浏览分类`
  String get browseCategories {
    return Intl.message('浏览分类', name: 'browseCategories', desc: '', args: []);
  }

  /// `热门标签`
  String get hotTags {
    return Intl.message('热门标签', name: 'hotTags', desc: '', args: []);
  }

  /// `推荐`
  String get recommend {
    return Intl.message('推荐', name: 'recommend', desc: '', args: []);
  }

  /// `{count} 篇`
  String postsCount(int count) {
    return Intl.message(
      '$count 篇',
      name: 'postsCount',
      desc: '',
      args: [count],
    );
  }

  /// `{count} 篇文章`
  String postsCountFull(int count) {
    return Intl.message(
      '$count 篇文章',
      name: 'postsCountFull',
      desc: '',
      args: [count],
    );
  }

  /// `我的收藏`
  String get myBookmarks {
    return Intl.message('我的收藏', name: 'myBookmarks', desc: '', args: []);
  }

  /// `已收藏 {count} 篇文章`
  String bookmarkedCount(int count) {
    return Intl.message(
      '已收藏 $count 篇文章',
      name: 'bookmarkedCount',
      desc: '',
      args: [count],
    );
  }

  /// `还没有收藏的文章`
  String get noBookmarks {
    return Intl.message('还没有收藏的文章', name: 'noBookmarks', desc: '', args: []);
  }

  /// `浏览文章时点击收藏按钮即可保存`
  String get bookmarkHint {
    return Intl.message(
      '浏览文章时点击收藏按钮即可保存',
      name: 'bookmarkHint',
      desc: '',
      args: [],
    );
  }

  /// `搜索文章、标签...`
  String get searchHint {
    return Intl.message('搜索文章、标签...', name: 'searchHint', desc: '', args: []);
  }

  /// `热门搜索`
  String get hotSearch {
    return Intl.message('热门搜索', name: 'hotSearch', desc: '', args: []);
  }

  /// `推荐阅读`
  String get recommendedReading {
    return Intl.message('推荐阅读', name: 'recommendedReading', desc: '', args: []);
  }

  /// `没有找到「{query}」相关内容`
  String noSearchResults(String query) {
    return Intl.message(
      '没有找到「$query」相关内容',
      name: 'noSearchResults',
      desc: '',
      args: [query],
    );
  }

  /// `相关标签`
  String get relatedTags {
    return Intl.message('相关标签', name: 'relatedTags', desc: '', args: []);
  }

  /// `相关文章 ({count})`
  String relatedPostsCount(int count) {
    return Intl.message(
      '相关文章 ($count)',
      name: 'relatedPostsCount',
      desc: '',
      args: [count],
    );
  }

  /// `评论 ({count})`
  String commentsCount(int count) {
    return Intl.message(
      '评论 ($count)',
      name: 'commentsCount',
      desc: '',
      args: [count],
    );
  }

  /// `浏览`
  String get views {
    return Intl.message('浏览', name: 'views', desc: '', args: []);
  }

  /// `点赞`
  String get likes {
    return Intl.message('点赞', name: 'likes', desc: '', args: []);
  }

  /// `评论`
  String get comments {
    return Intl.message('评论', name: 'comments', desc: '', args: []);
  }

  /// `分享`
  String get share {
    return Intl.message('分享', name: 'share', desc: '', args: []);
  }

  /// `回复`
  String get reply {
    return Intl.message('回复', name: 'reply', desc: '', args: []);
  }

  /// `写下你的评论...`
  String get writeComment {
    return Intl.message('写下你的评论...', name: 'writeComment', desc: '', args: []);
  }

  /// `今天`
  String get today {
    return Intl.message('今天', name: 'today', desc: '', args: []);
  }

  /// `昨天`
  String get yesterday {
    return Intl.message('昨天', name: 'yesterday', desc: '', args: []);
  }

  /// `{days} 天前`
  String daysAgo(int days) {
    return Intl.message('$days 天前', name: 'daysAgo', desc: '', args: [days]);
  }

  /// `{weeks} 周前`
  String weeksAgo(int weeks) {
    return Intl.message('$weeks 周前', name: 'weeksAgo', desc: '', args: [weeks]);
  }

  /// `{month} 月 {day} 日`
  String monthDay(int month, int day) {
    return Intl.message(
      '$month 月 $day 日',
      name: 'monthDay',
      desc: '',
      args: [month, day],
    );
  }

  /// `{year} 年 {month} 月 {day} 日`
  String yearMonthDay(int year, int month, int day) {
    return Intl.message(
      '$year 年 $month 月 $day 日',
      name: 'yearMonthDay',
      desc: '',
      args: [year, month, day],
    );
  }

  /// `{count} 篇相关文章`
  String relatedPostsCountFull(int count) {
    return Intl.message(
      '$count 篇相关文章',
      name: 'relatedPostsCountFull',
      desc: '',
      args: [count],
    );
  }

  /// `暂无相关文章`
  String get noRelatedPosts {
    return Intl.message('暂无相关文章', name: 'noRelatedPosts', desc: '', args: []);
  }

  /// `访客用户`
  String get guestUser {
    return Intl.message('访客用户', name: 'guestUser', desc: '', args: []);
  }

  /// `登录后享受更多功能`
  String get loginForMore {
    return Intl.message('登录后享受更多功能', name: 'loginForMore', desc: '', args: []);
  }

  /// `外观设置`
  String get appearance {
    return Intl.message('外观设置', name: 'appearance', desc: '', args: []);
  }

  /// `语言`
  String get language {
    return Intl.message('语言', name: 'language', desc: '', args: []);
  }

  /// `主题色`
  String get themeColor {
    return Intl.message('主题色', name: 'themeColor', desc: '', args: []);
  }

  /// `深色模式`
  String get darkMode {
    return Intl.message('深色模式', name: 'darkMode', desc: '', args: []);
  }

  /// `浅色`
  String get light {
    return Intl.message('浅色', name: 'light', desc: '', args: []);
  }

  /// `跟随系统`
  String get followSystem {
    return Intl.message('跟随系统', name: 'followSystem', desc: '', args: []);
  }

  /// `深色`
  String get dark {
    return Intl.message('深色', name: 'dark', desc: '', args: []);
  }

  /// `阅读统计`
  String get readingStats {
    return Intl.message('阅读统计', name: 'readingStats', desc: '', args: []);
  }

  /// `浏览历史`
  String get browseHistory {
    return Intl.message('浏览历史', name: 'browseHistory', desc: '', args: []);
  }

  /// `查看阅读记录`
  String get viewReadingHistory {
    return Intl.message(
      '查看阅读记录',
      name: 'viewReadingHistory',
      desc: '',
      args: [],
    );
  }

  /// `我的评论`
  String get myComments {
    return Intl.message('我的评论', name: 'myComments', desc: '', args: []);
  }

  /// `管理发表的评论`
  String get manageComments {
    return Intl.message('管理发表的评论', name: 'manageComments', desc: '', args: []);
  }

  /// `消息通知`
  String get notifications {
    return Intl.message('消息通知', name: 'notifications', desc: '', args: []);
  }

  /// `暂无新消息`
  String get noNewMessages {
    return Intl.message('暂无新消息', name: 'noNewMessages', desc: '', args: []);
  }

  /// `主题、语言等偏好`
  String get themeLanguagePrefs {
    return Intl.message(
      '主题、语言等偏好',
      name: 'themeLanguagePrefs',
      desc: '',
      args: [],
    );
  }

  /// `关于`
  String get about {
    return Intl.message('关于', name: 'about', desc: '', args: []);
  }

  /// `版本信息和帮助`
  String get versionInfo {
    return Intl.message('版本信息和帮助', name: 'versionInfo', desc: '', args: []);
  }

  /// `已读文章`
  String get readPosts {
    return Intl.message('已读文章', name: 'readPosts', desc: '', args: []);
  }

  /// `收藏文章`
  String get bookmarkedPostsLabel {
    return Intl.message(
      '收藏文章',
      name: 'bookmarkedPostsLabel',
      desc: '',
      args: [],
    );
  }

  /// `阅读时长`
  String get readingTime {
    return Intl.message('阅读时长', name: 'readingTime', desc: '', args: []);
  }

  /// `暂无评论`
  String get noCommentsYet {
    return Intl.message('暂无评论', name: 'noCommentsYet', desc: '', args: []);
  }

  /// `发生错误！`
  String get errorOccurred {
    return Intl.message('发生错误！', name: 'errorOccurred', desc: '', args: []);
  }

  /// `页面未找到`
  String get pageNotFound {
    return Intl.message('页面未找到', name: 'pageNotFound', desc: '', args: []);
  }

  /// `抱歉，您访问的页面不存在或已被移动。`
  String get pageNotFoundDesc {
    return Intl.message(
      '抱歉，您访问的页面不存在或已被移动。',
      name: 'pageNotFoundDesc',
      desc: '',
      args: [],
    );
  }

  /// `返回首页`
  String get backToHome {
    return Intl.message('返回首页', name: 'backToHome', desc: '', args: []);
  }

  /// `用户名`
  String get username {
    return Intl.message('用户名', name: 'username', desc: '', args: []);
  }

  /// `密码`
  String get password {
    return Intl.message('密码', name: 'password', desc: '', args: []);
  }

  /// `请输入用户名`
  String get usernameHint {
    return Intl.message('请输入用户名', name: 'usernameHint', desc: '', args: []);
  }

  /// `请输入密码`
  String get passwordHint {
    return Intl.message('请输入密码', name: 'passwordHint', desc: '', args: []);
  }

  /// `登录`
  String get loginButton {
    return Intl.message('登录', name: 'loginButton', desc: '', args: []);
  }

  /// `登录成功`
  String get loginSuccess {
    return Intl.message('登录成功', name: 'loginSuccess', desc: '', args: []);
  }

  /// `登录失败，请检查用户名和密码`
  String get loginFailed {
    return Intl.message(
      '登录失败，请检查用户名和密码',
      name: 'loginFailed',
      desc: '',
      args: [],
    );
  }

  /// `退出登录`
  String get logout {
    return Intl.message('退出登录', name: 'logout', desc: '', args: []);
  }

  /// `确定要退出登录吗？`
  String get logoutConfirm {
    return Intl.message('确定要退出登录吗？', name: 'logoutConfirm', desc: '', args: []);
  }

  /// `取消`
  String get cancel {
    return Intl.message('取消', name: 'cancel', desc: '', args: []);
  }

  /// `确定`
  String get confirm {
    return Intl.message('确定', name: 'confirm', desc: '', args: []);
  }

  /// `该功能即将上线，敬请期待`
  String get featureNotAvailable {
    return Intl.message(
      '该功能即将上线，敬请期待',
      name: 'featureNotAvailable',
      desc: '',
      args: [],
    );
  }

  /// `欢迎回来`
  String get welcomeBack {
    return Intl.message('欢迎回来', name: 'welcomeBack', desc: '', args: []);
  }

  /// `© 2026 GoWind CMS  ·  Powered by Flutter`
  String get footerText {
    return Intl.message(
      '© 2026 GoWind CMS  ·  Powered by Flutter',
      name: 'footerText',
      desc: '',
      args: [],
    );
  }

  /// `一个基于 Go 和 Flutter 构建的现代化内容管理系统`
  String get aboutSubtitle {
    return Intl.message(
      '一个基于 Go 和 Flutter 构建的现代化内容管理系统',
      name: 'aboutSubtitle',
      desc: '',
      args: [],
    );
  }

  /// `内容管理`
  String get aboutFeature1Title {
    return Intl.message('内容管理', name: 'aboutFeature1Title', desc: '', args: []);
  }

  /// `使用直观且强大的编辑器创建、编辑和发布内容`
  String get aboutFeature1Desc {
    return Intl.message(
      '使用直观且强大的编辑器创建、编辑和发布内容',
      name: 'aboutFeature1Desc',
      desc: '',
      args: [],
    );
  }

  /// `多语言支持`
  String get aboutFeature2Title {
    return Intl.message(
      '多语言支持',
      name: 'aboutFeature2Title',
      desc: '',
      args: [],
    );
  }

  /// `内置国际化能力，轻松服务全球受众`
  String get aboutFeature2Desc {
    return Intl.message(
      '内置国际化能力，轻松服务全球受众',
      name: 'aboutFeature2Desc',
      desc: '',
      args: [],
    );
  }

  /// `跨平台`
  String get aboutFeature3Title {
    return Intl.message('跨平台', name: 'aboutFeature3Title', desc: '', args: []);
  }

  /// `在 Web、iOS、Android 和桌面平台享受一致的体验`
  String get aboutFeature3Desc {
    return Intl.message(
      '在 Web、iOS、Android 和桌面平台享受一致的体验',
      name: 'aboutFeature3Desc',
      desc: '',
      args: [],
    );
  }

  /// `技术栈`
  String get aboutTechStack {
    return Intl.message('技术栈', name: 'aboutTechStack', desc: '', args: []);
  }

  /// `联系我们`
  String get contactUs {
    return Intl.message('联系我们', name: 'contactUs', desc: '', args: []);
  }

  /// `电子邮件`
  String get contactEmail {
    return Intl.message('电子邮件', name: 'contactEmail', desc: '', args: []);
  }

  /// `您可以通过 support@gowind.dev 联系我们，咨询任何问题、建议或反馈。我们通常在 1-2 个工作日内回复。`
  String get contactEmailDesc {
    return Intl.message(
      '您可以通过 support@gowind.dev 联系我们，咨询任何问题、建议或反馈。我们通常在 1-2 个工作日内回复。',
      name: 'contactEmailDesc',
      desc: '',
      args: [],
    );
  }

  /// `官方网站`
  String get contactWebsite {
    return Intl.message('官方网站', name: 'contactWebsite', desc: '', args: []);
  }

  /// `访问我们的官方网站 gowind.dev，获取最新动态、文档和社区资源。`
  String get contactWebsiteDesc {
    return Intl.message(
      '访问我们的官方网站 gowind.dev，获取最新动态、文档和社区资源。',
      name: 'contactWebsiteDesc',
      desc: '',
      args: [],
    );
  }

  /// `开发者社区`
  String get contactCommunity {
    return Intl.message('开发者社区', name: 'contactCommunity', desc: '', args: []);
  }

  /// `加入我们的 GitHub 开发者社区，报告问题、分享想法、参与项目贡献。`
  String get contactCommunityDesc {
    return Intl.message(
      '加入我们的 GitHub 开发者社区，报告问题、分享想法、参与项目贡献。',
      name: 'contactCommunityDesc',
      desc: '',
      args: [],
    );
  }

  /// `免责条款`
  String get disclaimer {
    return Intl.message('免责条款', name: 'disclaimer', desc: '', args: []);
  }

  /// `内容准确性`
  String get disclaimerContent1Title {
    return Intl.message(
      '内容准确性',
      name: 'disclaimerContent1Title',
      desc: '',
      args: [],
    );
  }

  /// `本平台提供的信息仅供参考。我们对内容的完整性、准确性或可靠性不作任何保证。您根据本平台信息采取的任何行动均由您自行承担风险。`
  String get disclaimerContent1Desc {
    return Intl.message(
      '本平台提供的信息仅供参考。我们对内容的完整性、准确性或可靠性不作任何保证。您根据本平台信息采取的任何行动均由您自行承担风险。',
      name: 'disclaimerContent1Desc',
      desc: '',
      args: [],
    );
  }

  /// `外部链接`
  String get disclaimerContent2Title {
    return Intl.message(
      '外部链接',
      name: 'disclaimerContent2Title',
      desc: '',
      args: [],
    );
  }

  /// `本平台可能包含指向外部网站的链接。我们无法控制这些网站的内容和性质，对因浏览或使用这些网站造成的任何损害不承担责任。`
  String get disclaimerContent2Desc {
    return Intl.message(
      '本平台可能包含指向外部网站的链接。我们无法控制这些网站的内容和性质，对因浏览或使用这些网站造成的任何损害不承担责任。',
      name: 'disclaimerContent2Desc',
      desc: '',
      args: [],
    );
  }

  /// `责任限制`
  String get disclaimerContent3Title {
    return Intl.message(
      '责任限制',
      name: 'disclaimerContent3Title',
      desc: '',
      args: [],
    );
  }

  /// `在任何情况下，我们均不对因使用本平台而产生的任何直接、间接、附带、后果性或特殊性损害承担责任。`
  String get disclaimerContent3Desc {
    return Intl.message(
      '在任何情况下，我们均不对因使用本平台而产生的任何直接、间接、附带、后果性或特殊性损害承担责任。',
      name: 'disclaimerContent3Desc',
      desc: '',
      args: [],
    );
  }

  /// `隐私协议`
  String get privacyPolicy {
    return Intl.message('隐私协议', name: 'privacyPolicy', desc: '', args: []);
  }

  /// `信息收集`
  String get privacyContent1Title {
    return Intl.message(
      '信息收集',
      name: 'privacyContent1Title',
      desc: '',
      args: [],
    );
  }

  /// `我们仅收集提供服务所需的最少个人信息，可能包括您的电子邮箱、用户名和使用偏好。我们不会出售或与第三方共享您的个人数据。`
  String get privacyContent1Desc {
    return Intl.message(
      '我们仅收集提供服务所需的最少个人信息，可能包括您的电子邮箱、用户名和使用偏好。我们不会出售或与第三方共享您的个人数据。',
      name: 'privacyContent1Desc',
      desc: '',
      args: [],
    );
  }

  /// `数据存储`
  String get privacyContent2Title {
    return Intl.message(
      '数据存储',
      name: 'privacyContent2Title',
      desc: '',
      args: [],
    );
  }

  /// `您的数据安全地存储在我们的服务器上，采用行业标准的加密技术。我们仅在提供服务所必需的期限或法律要求的期限内保留您的数据。`
  String get privacyContent2Desc {
    return Intl.message(
      '您的数据安全地存储在我们的服务器上，采用行业标准的加密技术。我们仅在提供服务所必需的期限或法律要求的期限内保留您的数据。',
      name: 'privacyContent2Desc',
      desc: '',
      args: [],
    );
  }

  /// `Cookie 与追踪`
  String get privacyContent3Title {
    return Intl.message(
      'Cookie 与追踪',
      name: 'privacyContent3Title',
      desc: '',
      args: [],
    );
  }

  /// `我们使用必要的 Cookie 以确保平台正常运行。可能会使用分析 Cookie 以改善用户体验，您可以在浏览器设置中禁用这些 Cookie。`
  String get privacyContent3Desc {
    return Intl.message(
      '我们使用必要的 Cookie 以确保平台正常运行。可能会使用分析 Cookie 以改善用户体验，您可以在浏览器设置中禁用这些 Cookie。',
      name: 'privacyContent3Desc',
      desc: '',
      args: [],
    );
  }

  /// `您的权利`
  String get privacyContent4Title {
    return Intl.message(
      '您的权利',
      name: 'privacyContent4Title',
      desc: '',
      args: [],
    );
  }

  /// `您有权随时访问、更正或删除您的个人数据。如有任何隐私相关请求，请联系我们的支持团队。`
  String get privacyContent4Desc {
    return Intl.message(
      '您有权随时访问、更正或删除您的个人数据。如有任何隐私相关请求，请联系我们的支持团队。',
      name: 'privacyContent4Desc',
      desc: '',
      args: [],
    );
  }

  /// `服务条款`
  String get termsOfService {
    return Intl.message('服务条款', name: 'termsOfService', desc: '', args: []);
  }

  /// `接受条款`
  String get termsContent1Title {
    return Intl.message('接受条款', name: 'termsContent1Title', desc: '', args: []);
  }

  /// `访问和使用本平台即表示您同意受本服务条款的约束。如果您不同意这些条款的任何部分，请勿使用本平台。`
  String get termsContent1Desc {
    return Intl.message(
      '访问和使用本平台即表示您同意受本服务条款的约束。如果您不同意这些条款的任何部分，请勿使用本平台。',
      name: 'termsContent1Desc',
      desc: '',
      args: [],
    );
  }

  /// `用户责任`
  String get termsContent2Title {
    return Intl.message('用户责任', name: 'termsContent2Title', desc: '', args: []);
  }

  /// `您有责任保管好您的账户信息。您同意不发布任何违法、有害、威胁、辱骂或其他不当内容。`
  String get termsContent2Desc {
    return Intl.message(
      '您有责任保管好您的账户信息。您同意不发布任何违法、有害、威胁、辱骂或其他不当内容。',
      name: 'termsContent2Desc',
      desc: '',
      args: [],
    );
  }

  /// `禁止行为`
  String get termsContent3Title {
    return Intl.message('禁止行为', name: 'termsContent3Title', desc: '', args: []);
  }

  /// `用户不得试图未经授权访问我们的系统、干扰平台运营或使用自动化工具未经许可抓取或收集数据。`
  String get termsContent3Desc {
    return Intl.message(
      '用户不得试图未经授权访问我们的系统、干扰平台运营或使用自动化工具未经许可抓取或收集数据。',
      name: 'termsContent3Desc',
      desc: '',
      args: [],
    );
  }

  /// `条款修改`
  String get termsContent4Title {
    return Intl.message('条款修改', name: 'termsContent4Title', desc: '', args: []);
  }

  /// `我们保留随时修改本条款的权利。在条款变更后继续使用本平台，即表示您接受修改后的条款。`
  String get termsContent4Desc {
    return Intl.message(
      '我们保留随时修改本条款的权利。在条款变更后继续使用本平台，即表示您接受修改后的条款。',
      name: 'termsContent4Desc',
      desc: '',
      args: [],
    );
  }
}

class AppLocalizationDelegate extends LocalizationsDelegate<S> {
  const AppLocalizationDelegate();

  List<Locale> get supportedLocales {
    return const <Locale>[
      Locale.fromSubtags(languageCode: 'zh', countryCode: 'CN'),
      Locale.fromSubtags(languageCode: 'en', countryCode: 'US'),
    ];
  }

  @override
  bool isSupported(Locale locale) => _isSupported(locale);
  @override
  Future<S> load(Locale locale) => S.load(locale);
  @override
  bool shouldReload(AppLocalizationDelegate old) => false;

  bool _isSupported(Locale locale) {
    for (var supportedLocale in supportedLocales) {
      if (supportedLocale.languageCode == locale.languageCode) {
        return true;
      }
    }
    return false;
  }
}
