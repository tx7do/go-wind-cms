// DO NOT EDIT. This is code generated via package:intl/generate_localized.dart
// This is a library that provides messages for a zh_CN locale. All the
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
  String get localeName => 'zh_CN';

  static String m0(count) => "已收藏 \$count 篇文章";

  static String m1(count) => "评论 (\$count)";

  static String m2(days) => "\$days 天前";

  static String m3(month, day) => "\$month 月 \$day 日";

  static String m4(query) => "没有找到「\$query」相关内容";

  static String m5(count) => "\$count 篇";

  static String m6(count) => "\$count 篇文章";

  static String m7(count) => "相关文章 (\$count)";

  static String m8(count) => "\$count 篇相关文章";

  static String m9(year, month, day) => "\$year 年 \$month 月 \$day 日";

  final messages = _notInlinedMessages(_notInlinedMessages);
  static Map<String, Function> _notInlinedMessages(_) => <String, Function>{
    "about": MessageLookupByLibrary.simpleMessage("关于"),
    "allLoaded": MessageLookupByLibrary.simpleMessage("— 已加载全部 —"),
    "allPosts": MessageLookupByLibrary.simpleMessage("全部文章"),
    "appName": MessageLookupByLibrary.simpleMessage("GoWind CMS"),
    "appearance": MessageLookupByLibrary.simpleMessage("外观设置"),
    "backToHome": MessageLookupByLibrary.simpleMessage("返回首页"),
    "bookmarkHint": MessageLookupByLibrary.simpleMessage("浏览文章时点击收藏按钮即可保存"),
    "bookmarkedCount": m0,
    "bookmarkedPostsLabel": MessageLookupByLibrary.simpleMessage("收藏文章"),
    "bookmarks": MessageLookupByLibrary.simpleMessage("收藏"),
    "browseCategories": MessageLookupByLibrary.simpleMessage("浏览分类"),
    "browseHistory": MessageLookupByLibrary.simpleMessage("浏览历史"),
    "comments": MessageLookupByLibrary.simpleMessage("评论"),
    "commentsCount": m1,
    "dark": MessageLookupByLibrary.simpleMessage("深色"),
    "darkMode": MessageLookupByLibrary.simpleMessage("深色模式"),
    "daysAgo": m2,
    "discover": MessageLookupByLibrary.simpleMessage("发现"),
    "errorOccurred": MessageLookupByLibrary.simpleMessage("发生错误！"),
    "followSystem": MessageLookupByLibrary.simpleMessage("跟随系统"),
    "guestUser": MessageLookupByLibrary.simpleMessage("访客用户"),
    "home": MessageLookupByLibrary.simpleMessage("首页"),
    "hotSearch": MessageLookupByLibrary.simpleMessage("热门搜索"),
    "hotTags": MessageLookupByLibrary.simpleMessage("热门标签"),
    "latestPosts": MessageLookupByLibrary.simpleMessage("最新文章"),
    "light": MessageLookupByLibrary.simpleMessage("浅色"),
    "likes": MessageLookupByLibrary.simpleMessage("点赞"),
    "login": MessageLookupByLibrary.simpleMessage("登录"),
    "loginForMore": MessageLookupByLibrary.simpleMessage("登录后享受更多功能"),
    "manageComments": MessageLookupByLibrary.simpleMessage("管理发表的评论"),
    "me": MessageLookupByLibrary.simpleMessage("我的"),
    "monthDay": m3,
    "myBookmarks": MessageLookupByLibrary.simpleMessage("我的收藏"),
    "myComments": MessageLookupByLibrary.simpleMessage("我的评论"),
    "noBookmarks": MessageLookupByLibrary.simpleMessage("还没有收藏的文章"),
    "noNewMessages": MessageLookupByLibrary.simpleMessage("暂无新消息"),
    "noRelatedPosts": MessageLookupByLibrary.simpleMessage("暂无相关文章"),
    "noSearchResults": m4,
    "notifications": MessageLookupByLibrary.simpleMessage("消息通知"),
    "pageNotFound": MessageLookupByLibrary.simpleMessage("页面未找到"),
    "pageNotFoundDesc": MessageLookupByLibrary.simpleMessage(
      "抱歉，您访问的页面不存在或已被移动。",
    ),
    "postsCount": m5,
    "postsCountFull": m6,
    "readPosts": MessageLookupByLibrary.simpleMessage("已读文章"),
    "readingStats": MessageLookupByLibrary.simpleMessage("阅读统计"),
    "readingTime": MessageLookupByLibrary.simpleMessage("阅读时长"),
    "recommend": MessageLookupByLibrary.simpleMessage("推荐"),
    "recommendedReading": MessageLookupByLibrary.simpleMessage("推荐阅读"),
    "relatedArticles": MessageLookupByLibrary.simpleMessage("相关文章"),
    "relatedPostsCount": m7,
    "relatedPostsCountFull": m8,
    "relatedTags": MessageLookupByLibrary.simpleMessage("相关标签"),
    "reply": MessageLookupByLibrary.simpleMessage("回复"),
    "search": MessageLookupByLibrary.simpleMessage("搜索"),
    "searchHint": MessageLookupByLibrary.simpleMessage("搜索文章、标签..."),
    "settings": MessageLookupByLibrary.simpleMessage("设置"),
    "share": MessageLookupByLibrary.simpleMessage("分享"),
    "themeColor": MessageLookupByLibrary.simpleMessage("主题色"),
    "themeLanguagePrefs": MessageLookupByLibrary.simpleMessage("主题、语言等偏好"),
    "today": MessageLookupByLibrary.simpleMessage("今天"),
    "versionInfo": MessageLookupByLibrary.simpleMessage("版本信息和帮助"),
    "viewReadingHistory": MessageLookupByLibrary.simpleMessage("查看阅读记录"),
    "views": MessageLookupByLibrary.simpleMessage("浏览"),
    "writeComment": MessageLookupByLibrary.simpleMessage("写下你的评论..."),
    "yearMonthDay": m9,
    "yesterday": MessageLookupByLibrary.simpleMessage("昨天"),
  };
}
