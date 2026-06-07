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

  /// `Home`
  String get home {
    return Intl.message('Home', name: 'home', desc: '', args: []);
  }

  /// `Discover`
  String get discover {
    return Intl.message('Discover', name: 'discover', desc: '', args: []);
  }

  /// `Bookmarks`
  String get bookmarks {
    return Intl.message('Bookmarks', name: 'bookmarks', desc: '', args: []);
  }

  /// `Me`
  String get me {
    return Intl.message('Me', name: 'me', desc: '', args: []);
  }

  /// `Search`
  String get search {
    return Intl.message('Search', name: 'search', desc: '', args: []);
  }

  /// `Settings`
  String get settings {
    return Intl.message('Settings', name: 'settings', desc: '', args: []);
  }

  /// `Login`
  String get login {
    return Intl.message('Login', name: 'login', desc: '', args: []);
  }

  /// `Latest Posts`
  String get latestPosts {
    return Intl.message(
      'Latest Posts',
      name: 'latestPosts',
      desc: '',
      args: [],
    );
  }

  /// `Related Articles`
  String get relatedArticles {
    return Intl.message(
      'Related Articles',
      name: 'relatedArticles',
      desc: '',
      args: [],
    );
  }

  /// `All Posts`
  String get allPosts {
    return Intl.message('All Posts', name: 'allPosts', desc: '', args: []);
  }

  /// `— All Loaded —`
  String get allLoaded {
    return Intl.message(
      '— All Loaded —',
      name: 'allLoaded',
      desc: '',
      args: [],
    );
  }

  /// `Browse Categories`
  String get browseCategories {
    return Intl.message(
      'Browse Categories',
      name: 'browseCategories',
      desc: '',
      args: [],
    );
  }

  /// `Hot Tags`
  String get hotTags {
    return Intl.message('Hot Tags', name: 'hotTags', desc: '', args: []);
  }

  /// `Recommend`
  String get recommend {
    return Intl.message('Recommend', name: 'recommend', desc: '', args: []);
  }

  /// `{count} posts`
  String postsCount(int count) {
    return Intl.message(
      '$count posts',
      name: 'postsCount',
      desc: '',
      args: [count],
    );
  }

  /// `{count} articles`
  String postsCountFull(int count) {
    return Intl.message(
      '$count articles',
      name: 'postsCountFull',
      desc: '',
      args: [count],
    );
  }

  /// `My Bookmarks`
  String get myBookmarks {
    return Intl.message(
      'My Bookmarks',
      name: 'myBookmarks',
      desc: '',
      args: [],
    );
  }

  /// `Bookmarked {count} articles`
  String bookmarkedCount(int count) {
    return Intl.message(
      'Bookmarked $count articles',
      name: 'bookmarkedCount',
      desc: '',
      args: [count],
    );
  }

  /// `No bookmarked articles yet`
  String get noBookmarks {
    return Intl.message(
      'No bookmarked articles yet',
      name: 'noBookmarks',
      desc: '',
      args: [],
    );
  }

  /// `Tap the bookmark button while reading to save`
  String get bookmarkHint {
    return Intl.message(
      'Tap the bookmark button while reading to save',
      name: 'bookmarkHint',
      desc: '',
      args: [],
    );
  }

  /// `Search articles, tags...`
  String get searchHint {
    return Intl.message(
      'Search articles, tags...',
      name: 'searchHint',
      desc: '',
      args: [],
    );
  }

  /// `Hot Searches`
  String get hotSearch {
    return Intl.message('Hot Searches', name: 'hotSearch', desc: '', args: []);
  }

  /// `Recommended Reading`
  String get recommendedReading {
    return Intl.message(
      'Recommended Reading',
      name: 'recommendedReading',
      desc: '',
      args: [],
    );
  }

  /// `No results found for "{query}"`
  String noSearchResults(String query) {
    return Intl.message(
      'No results found for "$query"',
      name: 'noSearchResults',
      desc: '',
      args: [query],
    );
  }

  /// `Related Tags`
  String get relatedTags {
    return Intl.message(
      'Related Tags',
      name: 'relatedTags',
      desc: '',
      args: [],
    );
  }

  /// `Related Articles ({count})`
  String relatedPostsCount(int count) {
    return Intl.message(
      'Related Articles ($count)',
      name: 'relatedPostsCount',
      desc: '',
      args: [count],
    );
  }

  /// `Comments ({count})`
  String commentsCount(int count) {
    return Intl.message(
      'Comments ($count)',
      name: 'commentsCount',
      desc: '',
      args: [count],
    );
  }

  /// `Views`
  String get views {
    return Intl.message('Views', name: 'views', desc: '', args: []);
  }

  /// `Likes`
  String get likes {
    return Intl.message('Likes', name: 'likes', desc: '', args: []);
  }

  /// `Comments`
  String get comments {
    return Intl.message('Comments', name: 'comments', desc: '', args: []);
  }

  /// `Share`
  String get share {
    return Intl.message('Share', name: 'share', desc: '', args: []);
  }

  /// `Reply`
  String get reply {
    return Intl.message('Reply', name: 'reply', desc: '', args: []);
  }

  /// `Write your comment...`
  String get writeComment {
    return Intl.message(
      'Write your comment...',
      name: 'writeComment',
      desc: '',
      args: [],
    );
  }

  /// `Today`
  String get today {
    return Intl.message('Today', name: 'today', desc: '', args: []);
  }

  /// `Yesterday`
  String get yesterday {
    return Intl.message('Yesterday', name: 'yesterday', desc: '', args: []);
  }

  /// `{days} days ago`
  String daysAgo(int days) {
    return Intl.message(
      '$days days ago',
      name: 'daysAgo',
      desc: '',
      args: [days],
    );
  }

  /// `{month}/{day}`
  String monthDay(int month, int day) {
    return Intl.message(
      '$month/$day',
      name: 'monthDay',
      desc: '',
      args: [month, day],
    );
  }

  /// `{year}/{month}/{day}`
  String yearMonthDay(int year, int month, int day) {
    return Intl.message(
      '$year/$month/$day',
      name: 'yearMonthDay',
      desc: '',
      args: [year, month, day],
    );
  }

  /// `{count} related articles`
  String relatedPostsCountFull(int count) {
    return Intl.message(
      '$count related articles',
      name: 'relatedPostsCountFull',
      desc: '',
      args: [count],
    );
  }

  /// `No related articles yet`
  String get noRelatedPosts {
    return Intl.message(
      'No related articles yet',
      name: 'noRelatedPosts',
      desc: '',
      args: [],
    );
  }

  /// `Guest`
  String get guestUser {
    return Intl.message('Guest', name: 'guestUser', desc: '', args: []);
  }

  /// `Login for more features`
  String get loginForMore {
    return Intl.message(
      'Login for more features',
      name: 'loginForMore',
      desc: '',
      args: [],
    );
  }

  /// `Appearance`
  String get appearance {
    return Intl.message('Appearance', name: 'appearance', desc: '', args: []);
  }

  /// `Theme Color`
  String get themeColor {
    return Intl.message('Theme Color', name: 'themeColor', desc: '', args: []);
  }

  /// `Dark Mode`
  String get darkMode {
    return Intl.message('Dark Mode', name: 'darkMode', desc: '', args: []);
  }

  /// `Light`
  String get light {
    return Intl.message('Light', name: 'light', desc: '', args: []);
  }

  /// `System`
  String get followSystem {
    return Intl.message('System', name: 'followSystem', desc: '', args: []);
  }

  /// `Dark`
  String get dark {
    return Intl.message('Dark', name: 'dark', desc: '', args: []);
  }

  /// `Reading Stats`
  String get readingStats {
    return Intl.message(
      'Reading Stats',
      name: 'readingStats',
      desc: '',
      args: [],
    );
  }

  /// `Browse History`
  String get browseHistory {
    return Intl.message(
      'Browse History',
      name: 'browseHistory',
      desc: '',
      args: [],
    );
  }

  /// `View reading history`
  String get viewReadingHistory {
    return Intl.message(
      'View reading history',
      name: 'viewReadingHistory',
      desc: '',
      args: [],
    );
  }

  /// `My Comments`
  String get myComments {
    return Intl.message('My Comments', name: 'myComments', desc: '', args: []);
  }

  /// `Manage your comments`
  String get manageComments {
    return Intl.message(
      'Manage your comments',
      name: 'manageComments',
      desc: '',
      args: [],
    );
  }

  /// `Notifications`
  String get notifications {
    return Intl.message(
      'Notifications',
      name: 'notifications',
      desc: '',
      args: [],
    );
  }

  /// `No new messages`
  String get noNewMessages {
    return Intl.message(
      'No new messages',
      name: 'noNewMessages',
      desc: '',
      args: [],
    );
  }

  /// `Theme, language & preferences`
  String get themeLanguagePrefs {
    return Intl.message(
      'Theme, language & preferences',
      name: 'themeLanguagePrefs',
      desc: '',
      args: [],
    );
  }

  /// `About`
  String get about {
    return Intl.message('About', name: 'about', desc: '', args: []);
  }

  /// `Version info & help`
  String get versionInfo {
    return Intl.message(
      'Version info & help',
      name: 'versionInfo',
      desc: '',
      args: [],
    );
  }

  /// `Articles Read`
  String get readPosts {
    return Intl.message('Articles Read', name: 'readPosts', desc: '', args: []);
  }

  /// `Bookmarked`
  String get bookmarkedPostsLabel {
    return Intl.message(
      'Bookmarked',
      name: 'bookmarkedPostsLabel',
      desc: '',
      args: [],
    );
  }

  /// `Reading Time`
  String get readingTime {
    return Intl.message(
      'Reading Time',
      name: 'readingTime',
      desc: '',
      args: [],
    );
  }

  /// `Error Occurred!`
  String get errorOccurred {
    return Intl.message(
      'Error Occurred!',
      name: 'errorOccurred',
      desc: '',
      args: [],
    );
  }

  /// `Page Not Found`
  String get pageNotFound {
    return Intl.message(
      'Page Not Found',
      name: 'pageNotFound',
      desc: '',
      args: [],
    );
  }

  /// `Sorry, the page you are looking for does not exist or has been moved.`
  String get pageNotFoundDesc {
    return Intl.message(
      'Sorry, the page you are looking for does not exist or has been moved.',
      name: 'pageNotFoundDesc',
      desc: '',
      args: [],
    );
  }

  /// `Back to Home`
  String get backToHome {
    return Intl.message('Back to Home', name: 'backToHome', desc: '', args: []);
  }
}

class AppLocalizationDelegate extends LocalizationsDelegate<S> {
  const AppLocalizationDelegate();

  List<Locale> get supportedLocales {
    return const <Locale>[
      Locale.fromSubtags(languageCode: 'en'),
      Locale.fromSubtags(languageCode: 'zh', countryCode: 'CN'),
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
