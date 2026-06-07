// DO NOT EDIT. This is code generated via package:intl/generate_localized.dart
// This is a library that provides messages for a en locale. All the
// messages from the main program should be duplicated here with the same
// function name.

// Ignore issues from commonly used lints in this file.
// ignore_for_file:unnecessary_brace_in_string_interps, unnecessary_new
// ignore_for_file:prefer_single_quotes,comment_references, directives_ordering
// ignore_for_file:annotate_overrides,prefer_generic_function_type_aliases
// ignore_for_file:unused_import, file_names, avoid_escaping_inner_quotes
// ignore_for_file:unnecessary_string_interpolations, unnecessary_string_escapes

import 'package:intl/intl.dart';
import 'package:intl/message_lookup_by_library.dart';

final messages = new MessageLookup();

typedef String MessageIfAbsent(String messageStr, List<dynamic> args);

class MessageLookup extends MessageLookupByLibrary {
  String get localeName => 'en';

  static String m0(count) => "Bookmarked ${count} articles";

  static String m1(count) => "Comments (${count})";

  static String m2(days) => "${days} days ago";

  static String m3(month, day) => "${month}/${day}";

  static String m4(query) => "No results found for \"${query}\"";

  static String m5(count) => "${count} posts";

  static String m6(count) => "${count} articles";

  static String m7(count) => "Related Articles (${count})";

  static String m8(count) => "${count} related articles";

  static String m9(year, month, day) => "${year}/${month}/${day}";

  final messages = _notInlinedMessages(_notInlinedMessages);
  static Map<String, Function> _notInlinedMessages(_) => <String, Function>{
    "about": MessageLookupByLibrary.simpleMessage("About"),
    "allLoaded": MessageLookupByLibrary.simpleMessage("— All Loaded —"),
    "allPosts": MessageLookupByLibrary.simpleMessage("All Posts"),
    "appName": MessageLookupByLibrary.simpleMessage("GoWind CMS"),
    "appearance": MessageLookupByLibrary.simpleMessage("Appearance"),
    "backToHome": MessageLookupByLibrary.simpleMessage("Back to Home"),
    "bookmarkHint": MessageLookupByLibrary.simpleMessage(
      "Tap the bookmark button while reading to save",
    ),
    "bookmarkedCount": m0,
    "bookmarkedPostsLabel": MessageLookupByLibrary.simpleMessage("Bookmarked"),
    "bookmarks": MessageLookupByLibrary.simpleMessage("Bookmarks"),
    "browseCategories": MessageLookupByLibrary.simpleMessage(
      "Browse Categories",
    ),
    "browseHistory": MessageLookupByLibrary.simpleMessage("Browse History"),
    "comments": MessageLookupByLibrary.simpleMessage("Comments"),
    "commentsCount": m1,
    "dark": MessageLookupByLibrary.simpleMessage("Dark"),
    "darkMode": MessageLookupByLibrary.simpleMessage("Dark Mode"),
    "daysAgo": m2,
    "discover": MessageLookupByLibrary.simpleMessage("Discover"),
    "errorOccurred": MessageLookupByLibrary.simpleMessage("Error Occurred!"),
    "followSystem": MessageLookupByLibrary.simpleMessage("System"),
    "guestUser": MessageLookupByLibrary.simpleMessage("Guest"),
    "home": MessageLookupByLibrary.simpleMessage("Home"),
    "hotSearch": MessageLookupByLibrary.simpleMessage("Hot Searches"),
    "hotTags": MessageLookupByLibrary.simpleMessage("Hot Tags"),
    "latestPosts": MessageLookupByLibrary.simpleMessage("Latest Posts"),
    "light": MessageLookupByLibrary.simpleMessage("Light"),
    "likes": MessageLookupByLibrary.simpleMessage("Likes"),
    "login": MessageLookupByLibrary.simpleMessage("Login"),
    "loginForMore": MessageLookupByLibrary.simpleMessage(
      "Login for more features",
    ),
    "manageComments": MessageLookupByLibrary.simpleMessage(
      "Manage your comments",
    ),
    "me": MessageLookupByLibrary.simpleMessage("Me"),
    "monthDay": m3,
    "myBookmarks": MessageLookupByLibrary.simpleMessage("My Bookmarks"),
    "myComments": MessageLookupByLibrary.simpleMessage("My Comments"),
    "noBookmarks": MessageLookupByLibrary.simpleMessage(
      "No bookmarked articles yet",
    ),
    "noNewMessages": MessageLookupByLibrary.simpleMessage("No new messages"),
    "noRelatedPosts": MessageLookupByLibrary.simpleMessage(
      "No related articles yet",
    ),
    "noSearchResults": m4,
    "notifications": MessageLookupByLibrary.simpleMessage("Notifications"),
    "pageNotFound": MessageLookupByLibrary.simpleMessage("Page Not Found"),
    "pageNotFoundDesc": MessageLookupByLibrary.simpleMessage(
      "Sorry, the page you are looking for does not exist or has been moved.",
    ),
    "postsCount": m5,
    "postsCountFull": m6,
    "readPosts": MessageLookupByLibrary.simpleMessage("Articles Read"),
    "readingStats": MessageLookupByLibrary.simpleMessage("Reading Stats"),
    "readingTime": MessageLookupByLibrary.simpleMessage("Reading Time"),
    "recommend": MessageLookupByLibrary.simpleMessage("Recommend"),
    "recommendedReading": MessageLookupByLibrary.simpleMessage(
      "Recommended Reading",
    ),
    "relatedArticles": MessageLookupByLibrary.simpleMessage("Related Articles"),
    "relatedPostsCount": m7,
    "relatedPostsCountFull": m8,
    "relatedTags": MessageLookupByLibrary.simpleMessage("Related Tags"),
    "reply": MessageLookupByLibrary.simpleMessage("Reply"),
    "search": MessageLookupByLibrary.simpleMessage("Search"),
    "searchHint": MessageLookupByLibrary.simpleMessage(
      "Search articles, tags...",
    ),
    "settings": MessageLookupByLibrary.simpleMessage("Settings"),
    "share": MessageLookupByLibrary.simpleMessage("Share"),
    "themeColor": MessageLookupByLibrary.simpleMessage("Theme Color"),
    "themeLanguagePrefs": MessageLookupByLibrary.simpleMessage(
      "Theme, language & preferences",
    ),
    "today": MessageLookupByLibrary.simpleMessage("Today"),
    "versionInfo": MessageLookupByLibrary.simpleMessage("Version info & help"),
    "viewReadingHistory": MessageLookupByLibrary.simpleMessage(
      "View reading history",
    ),
    "views": MessageLookupByLibrary.simpleMessage("Views"),
    "writeComment": MessageLookupByLibrary.simpleMessage(
      "Write your comment...",
    ),
    "yearMonthDay": m9,
    "yesterday": MessageLookupByLibrary.simpleMessage("Yesterday"),
  };
}
