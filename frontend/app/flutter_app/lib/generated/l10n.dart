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

  /// `{weeks} weeks ago`
  String weeksAgo(int weeks) {
    return Intl.message(
      '$weeks weeks ago',
      name: 'weeksAgo',
      desc: '',
      args: [weeks],
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

  /// `Username`
  String get username {
    return Intl.message('Username', name: 'username', desc: '', args: []);
  }

  /// `Password`
  String get password {
    return Intl.message('Password', name: 'password', desc: '', args: []);
  }

  /// `Enter username`
  String get usernameHint {
    return Intl.message(
      'Enter username',
      name: 'usernameHint',
      desc: '',
      args: [],
    );
  }

  /// `Enter password`
  String get passwordHint {
    return Intl.message(
      'Enter password',
      name: 'passwordHint',
      desc: '',
      args: [],
    );
  }

  /// `Login`
  String get loginButton {
    return Intl.message('Login', name: 'loginButton', desc: '', args: []);
  }

  /// `Login successful`
  String get loginSuccess {
    return Intl.message(
      'Login successful',
      name: 'loginSuccess',
      desc: '',
      args: [],
    );
  }

  /// `Login failed, please check username and password`
  String get loginFailed {
    return Intl.message(
      'Login failed, please check username and password',
      name: 'loginFailed',
      desc: '',
      args: [],
    );
  }

  /// `Logout`
  String get logout {
    return Intl.message('Logout', name: 'logout', desc: '', args: []);
  }

  /// `Are you sure you want to logout?`
  String get logoutConfirm {
    return Intl.message(
      'Are you sure you want to logout?',
      name: 'logoutConfirm',
      desc: '',
      args: [],
    );
  }

  /// `Cancel`
  String get cancel {
    return Intl.message('Cancel', name: 'cancel', desc: '', args: []);
  }

  /// `Confirm`
  String get confirm {
    return Intl.message('Confirm', name: 'confirm', desc: '', args: []);
  }

  /// `This feature is coming soon`
  String get featureNotAvailable {
    return Intl.message(
      'This feature is coming soon',
      name: 'featureNotAvailable',
      desc: '',
      args: [],
    );
  }

  /// `Welcome back`
  String get welcomeBack {
    return Intl.message(
      'Welcome back',
      name: 'welcomeBack',
      desc: '',
      args: [],
    );
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

  /// `A modern content management system powered by Go and Flutter`
  String get aboutSubtitle {
    return Intl.message(
      'A modern content management system powered by Go and Flutter',
      name: 'aboutSubtitle',
      desc: '',
      args: [],
    );
  }

  /// `Content Management`
  String get aboutFeature1Title {
    return Intl.message(
      'Content Management',
      name: 'aboutFeature1Title',
      desc: '',
      args: [],
    );
  }

  /// `Create, edit and publish content with an intuitive and powerful editor`
  String get aboutFeature1Desc {
    return Intl.message(
      'Create, edit and publish content with an intuitive and powerful editor',
      name: 'aboutFeature1Desc',
      desc: '',
      args: [],
    );
  }

  /// `Multi-language Support`
  String get aboutFeature2Title {
    return Intl.message(
      'Multi-language Support',
      name: 'aboutFeature2Title',
      desc: '',
      args: [],
    );
  }

  /// `Built-in internationalization to serve a global audience effortlessly`
  String get aboutFeature2Desc {
    return Intl.message(
      'Built-in internationalization to serve a global audience effortlessly',
      name: 'aboutFeature2Desc',
      desc: '',
      args: [],
    );
  }

  /// `Cross-platform`
  String get aboutFeature3Title {
    return Intl.message(
      'Cross-platform',
      name: 'aboutFeature3Title',
      desc: '',
      args: [],
    );
  }

  /// `Seamless experience across Web, iOS, Android and desktop platforms`
  String get aboutFeature3Desc {
    return Intl.message(
      'Seamless experience across Web, iOS, Android and desktop platforms',
      name: 'aboutFeature3Desc',
      desc: '',
      args: [],
    );
  }

  /// `Built with`
  String get aboutTechStack {
    return Intl.message(
      'Built with',
      name: 'aboutTechStack',
      desc: '',
      args: [],
    );
  }

  /// `Contact Us`
  String get contactUs {
    return Intl.message('Contact Us', name: 'contactUs', desc: '', args: []);
  }

  /// `Email`
  String get contactEmail {
    return Intl.message('Email', name: 'contactEmail', desc: '', args: []);
  }

  /// `You can reach us via email at support@gowind.dev for any questions, suggestions or feedback. We typically respond within 1-2 business days.`
  String get contactEmailDesc {
    return Intl.message(
      'You can reach us via email at support@gowind.dev for any questions, suggestions or feedback. We typically respond within 1-2 business days.',
      name: 'contactEmailDesc',
      desc: '',
      args: [],
    );
  }

  /// `Website`
  String get contactWebsite {
    return Intl.message('Website', name: 'contactWebsite', desc: '', args: []);
  }

  /// `Visit our official website gowind.dev for the latest updates, documentation, and community resources.`
  String get contactWebsiteDesc {
    return Intl.message(
      'Visit our official website gowind.dev for the latest updates, documentation, and community resources.',
      name: 'contactWebsiteDesc',
      desc: '',
      args: [],
    );
  }

  /// `Community`
  String get contactCommunity {
    return Intl.message(
      'Community',
      name: 'contactCommunity',
      desc: '',
      args: [],
    );
  }

  /// `Join our developer community on GitHub to report issues, share ideas, and contribute to the project.`
  String get contactCommunityDesc {
    return Intl.message(
      'Join our developer community on GitHub to report issues, share ideas, and contribute to the project.',
      name: 'contactCommunityDesc',
      desc: '',
      args: [],
    );
  }

  /// `Disclaimer`
  String get disclaimer {
    return Intl.message('Disclaimer', name: 'disclaimer', desc: '', args: []);
  }

  /// `Content Accuracy`
  String get disclaimerContent1Title {
    return Intl.message(
      'Content Accuracy',
      name: 'disclaimerContent1Title',
      desc: '',
      args: [],
    );
  }

  /// `The information provided on this platform is for general informational purposes only. We make no warranties about the completeness, accuracy, or reliability of the content. Any action you take upon the information is strictly at your own risk.`
  String get disclaimerContent1Desc {
    return Intl.message(
      'The information provided on this platform is for general informational purposes only. We make no warranties about the completeness, accuracy, or reliability of the content. Any action you take upon the information is strictly at your own risk.',
      name: 'disclaimerContent1Desc',
      desc: '',
      args: [],
    );
  }

  /// `External Links`
  String get disclaimerContent2Title {
    return Intl.message(
      'External Links',
      name: 'disclaimerContent2Title',
      desc: '',
      args: [],
    );
  }

  /// `This platform may contain links to external websites. We have no control over the content and nature of these sites and are not responsible for any damages from browsing or using them.`
  String get disclaimerContent2Desc {
    return Intl.message(
      'This platform may contain links to external websites. We have no control over the content and nature of these sites and are not responsible for any damages from browsing or using them.',
      name: 'disclaimerContent2Desc',
      desc: '',
      args: [],
    );
  }

  /// `Limitation of Liability`
  String get disclaimerContent3Title {
    return Intl.message(
      'Limitation of Liability',
      name: 'disclaimerContent3Title',
      desc: '',
      args: [],
    );
  }

  /// `In no event shall we be liable for any direct, indirect, incidental, consequential, or special damages arising out of or in connection with the use of this platform.`
  String get disclaimerContent3Desc {
    return Intl.message(
      'In no event shall we be liable for any direct, indirect, incidental, consequential, or special damages arising out of or in connection with the use of this platform.',
      name: 'disclaimerContent3Desc',
      desc: '',
      args: [],
    );
  }

  /// `Privacy Policy`
  String get privacyPolicy {
    return Intl.message(
      'Privacy Policy',
      name: 'privacyPolicy',
      desc: '',
      args: [],
    );
  }

  /// `Information Collection`
  String get privacyContent1Title {
    return Intl.message(
      'Information Collection',
      name: 'privacyContent1Title',
      desc: '',
      args: [],
    );
  }

  /// `We collect minimal personal information necessary to provide our services. This may include your email address, username, and usage preferences. We do not sell or share your personal data with third parties.`
  String get privacyContent1Desc {
    return Intl.message(
      'We collect minimal personal information necessary to provide our services. This may include your email address, username, and usage preferences. We do not sell or share your personal data with third parties.',
      name: 'privacyContent1Desc',
      desc: '',
      args: [],
    );
  }

  /// `Data Storage`
  String get privacyContent2Title {
    return Intl.message(
      'Data Storage',
      name: 'privacyContent2Title',
      desc: '',
      args: [],
    );
  }

  /// `Your data is stored securely on our servers with industry-standard encryption. We retain your data only for as long as necessary to provide the services or as required by law.`
  String get privacyContent2Desc {
    return Intl.message(
      'Your data is stored securely on our servers with industry-standard encryption. We retain your data only for as long as necessary to provide the services or as required by law.',
      name: 'privacyContent2Desc',
      desc: '',
      args: [],
    );
  }

  /// `Cookies & Tracking`
  String get privacyContent3Title {
    return Intl.message(
      'Cookies & Tracking',
      name: 'privacyContent3Title',
      desc: '',
      args: [],
    );
  }

  /// `We use essential cookies to ensure the proper functioning of the platform. Analytics cookies may be used to improve user experience, which can be disabled in your browser settings.`
  String get privacyContent3Desc {
    return Intl.message(
      'We use essential cookies to ensure the proper functioning of the platform. Analytics cookies may be used to improve user experience, which can be disabled in your browser settings.',
      name: 'privacyContent3Desc',
      desc: '',
      args: [],
    );
  }

  /// `Your Rights`
  String get privacyContent4Title {
    return Intl.message(
      'Your Rights',
      name: 'privacyContent4Title',
      desc: '',
      args: [],
    );
  }

  /// `You have the right to access, correct, or delete your personal data at any time. Contact our support team for any privacy-related requests.`
  String get privacyContent4Desc {
    return Intl.message(
      'You have the right to access, correct, or delete your personal data at any time. Contact our support team for any privacy-related requests.',
      name: 'privacyContent4Desc',
      desc: '',
      args: [],
    );
  }

  /// `Terms of Service`
  String get termsOfService {
    return Intl.message(
      'Terms of Service',
      name: 'termsOfService',
      desc: '',
      args: [],
    );
  }

  /// `Acceptance of Terms`
  String get termsContent1Title {
    return Intl.message(
      'Acceptance of Terms',
      name: 'termsContent1Title',
      desc: '',
      args: [],
    );
  }

  /// `By accessing and using this platform, you agree to be bound by these Terms of Service. If you do not agree with any part of these terms, you must not use the platform.`
  String get termsContent1Desc {
    return Intl.message(
      'By accessing and using this platform, you agree to be bound by these Terms of Service. If you do not agree with any part of these terms, you must not use the platform.',
      name: 'termsContent1Desc',
      desc: '',
      args: [],
    );
  }

  /// `User Responsibilities`
  String get termsContent2Title {
    return Intl.message(
      'User Responsibilities',
      name: 'termsContent2Title',
      desc: '',
      args: [],
    );
  }

  /// `You are responsible for maintaining the confidentiality of your account. You agree not to post any content that is unlawful, harmful, threatening, abusive, or otherwise objectionable.`
  String get termsContent2Desc {
    return Intl.message(
      'You are responsible for maintaining the confidentiality of your account. You agree not to post any content that is unlawful, harmful, threatening, abusive, or otherwise objectionable.',
      name: 'termsContent2Desc',
      desc: '',
      args: [],
    );
  }

  /// `Prohibited Activities`
  String get termsContent3Title {
    return Intl.message(
      'Prohibited Activities',
      name: 'termsContent3Title',
      desc: '',
      args: [],
    );
  }

  /// `Users must not attempt to gain unauthorized access to our systems, interfere with the platform's operation, or use automated tools to scrape or collect data without permission.`
  String get termsContent3Desc {
    return Intl.message(
      'Users must not attempt to gain unauthorized access to our systems, interfere with the platform\'s operation, or use automated tools to scrape or collect data without permission.',
      name: 'termsContent3Desc',
      desc: '',
      args: [],
    );
  }

  /// `Modifications`
  String get termsContent4Title {
    return Intl.message(
      'Modifications',
      name: 'termsContent4Title',
      desc: '',
      args: [],
    );
  }

  /// `We reserve the right to modify these terms at any time. Continued use of the platform after changes constitutes acceptance of the updated terms.`
  String get termsContent4Desc {
    return Intl.message(
      'We reserve the right to modify these terms at any time. Continued use of the platform after changes constitutes acceptance of the updated terms.',
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
