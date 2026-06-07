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

  static String m0(count) => "已收藏 ${count} 篇文章";

  static String m1(count) => "评论 (${count})";

  static String m2(days) => "${days} 天前";

  static String m3(month, day) => "${month} 月 ${day} 日";

  static String m4(query) => "没有找到「${query}」相关内容";

  static String m5(count) => "${count} 篇";

  static String m6(count) => "${count} 篇文章";

  static String m7(count) => "相关文章 (${count})";

  static String m8(count) => "${count} 篇相关文章";

  static String m9(weeks) => "${weeks} 周前";

  static String m10(year, month, day) => "${year} 年 ${month} 月 ${day} 日";

  final messages = _notInlinedMessages(_notInlinedMessages);
  static Map<String, Function> _notInlinedMessages(_) => <String, Function>{
    "about": MessageLookupByLibrary.simpleMessage("关于"),
    "aboutFeature1Desc": MessageLookupByLibrary.simpleMessage(
      "使用直观且强大的编辑器创建、编辑和发布内容",
    ),
    "aboutFeature1Title": MessageLookupByLibrary.simpleMessage("内容管理"),
    "aboutFeature2Desc": MessageLookupByLibrary.simpleMessage(
      "内置国际化能力，轻松服务全球受众",
    ),
    "aboutFeature2Title": MessageLookupByLibrary.simpleMessage("多语言支持"),
    "aboutFeature3Desc": MessageLookupByLibrary.simpleMessage(
      "在 Web、iOS、Android 和桌面平台享受一致的体验",
    ),
    "aboutFeature3Title": MessageLookupByLibrary.simpleMessage("跨平台"),
    "aboutSubtitle": MessageLookupByLibrary.simpleMessage(
      "一个基于 Go 和 Flutter 构建的现代化内容管理系统",
    ),
    "aboutTechStack": MessageLookupByLibrary.simpleMessage("技术栈"),
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
    "cancel": MessageLookupByLibrary.simpleMessage("取消"),
    "comments": MessageLookupByLibrary.simpleMessage("评论"),
    "commentsCount": m1,
    "confirm": MessageLookupByLibrary.simpleMessage("确定"),
    "contactCommunity": MessageLookupByLibrary.simpleMessage("开发者社区"),
    "contactCommunityDesc": MessageLookupByLibrary.simpleMessage(
      "加入我们的 GitHub 开发者社区，报告问题、分享想法、参与项目贡献。",
    ),
    "contactEmail": MessageLookupByLibrary.simpleMessage("电子邮件"),
    "contactEmailDesc": MessageLookupByLibrary.simpleMessage(
      "您可以通过 support@gowind.dev 联系我们，咨询任何问题、建议或反馈。我们通常在 1-2 个工作日内回复。",
    ),
    "contactUs": MessageLookupByLibrary.simpleMessage("联系我们"),
    "contactWebsite": MessageLookupByLibrary.simpleMessage("官方网站"),
    "contactWebsiteDesc": MessageLookupByLibrary.simpleMessage(
      "访问我们的官方网站 gowind.dev，获取最新动态、文档和社区资源。",
    ),
    "dark": MessageLookupByLibrary.simpleMessage("深色"),
    "darkMode": MessageLookupByLibrary.simpleMessage("深色模式"),
    "daysAgo": m2,
    "disclaimer": MessageLookupByLibrary.simpleMessage("免责条款"),
    "disclaimerContent1Desc": MessageLookupByLibrary.simpleMessage(
      "本平台提供的信息仅供参考。我们对内容的完整性、准确性或可靠性不作任何保证。您根据本平台信息采取的任何行动均由您自行承担风险。",
    ),
    "disclaimerContent1Title": MessageLookupByLibrary.simpleMessage("内容准确性"),
    "disclaimerContent2Desc": MessageLookupByLibrary.simpleMessage(
      "本平台可能包含指向外部网站的链接。我们无法控制这些网站的内容和性质，对因浏览或使用这些网站造成的任何损害不承担责任。",
    ),
    "disclaimerContent2Title": MessageLookupByLibrary.simpleMessage("外部链接"),
    "disclaimerContent3Desc": MessageLookupByLibrary.simpleMessage(
      "在任何情况下，我们均不对因使用本平台而产生的任何直接、间接、附带、后果性或特殊性损害承担责任。",
    ),
    "disclaimerContent3Title": MessageLookupByLibrary.simpleMessage("责任限制"),
    "discover": MessageLookupByLibrary.simpleMessage("发现"),
    "errorOccurred": MessageLookupByLibrary.simpleMessage("发生错误！"),
    "featureNotAvailable": MessageLookupByLibrary.simpleMessage("该功能即将上线，敬请期待"),
    "followSystem": MessageLookupByLibrary.simpleMessage("跟随系统"),
    "footerText": MessageLookupByLibrary.simpleMessage(
      "© 2026 GoWind CMS  ·  Powered by Flutter",
    ),
    "guestUser": MessageLookupByLibrary.simpleMessage("访客用户"),
    "home": MessageLookupByLibrary.simpleMessage("首页"),
    "hotSearch": MessageLookupByLibrary.simpleMessage("热门搜索"),
    "hotTags": MessageLookupByLibrary.simpleMessage("热门标签"),
    "language": MessageLookupByLibrary.simpleMessage("语言"),
    "latestPosts": MessageLookupByLibrary.simpleMessage("最新文章"),
    "light": MessageLookupByLibrary.simpleMessage("浅色"),
    "likes": MessageLookupByLibrary.simpleMessage("点赞"),
    "login": MessageLookupByLibrary.simpleMessage("登录"),
    "loginButton": MessageLookupByLibrary.simpleMessage("登录"),
    "loginFailed": MessageLookupByLibrary.simpleMessage("登录失败，请检查用户名和密码"),
    "loginForMore": MessageLookupByLibrary.simpleMessage("登录后享受更多功能"),
    "loginSuccess": MessageLookupByLibrary.simpleMessage("登录成功"),
    "logout": MessageLookupByLibrary.simpleMessage("退出登录"),
    "logoutConfirm": MessageLookupByLibrary.simpleMessage("确定要退出登录吗？"),
    "manageComments": MessageLookupByLibrary.simpleMessage("管理发表的评论"),
    "me": MessageLookupByLibrary.simpleMessage("我的"),
    "monthDay": m3,
    "myBookmarks": MessageLookupByLibrary.simpleMessage("我的收藏"),
    "myComments": MessageLookupByLibrary.simpleMessage("我的评论"),
    "noBookmarks": MessageLookupByLibrary.simpleMessage("还没有收藏的文章"),
    "noCommentsYet": MessageLookupByLibrary.simpleMessage("暂无评论"),
    "noNewMessages": MessageLookupByLibrary.simpleMessage("暂无新消息"),
    "noRelatedPosts": MessageLookupByLibrary.simpleMessage("暂无相关文章"),
    "noSearchResults": m4,
    "notifications": MessageLookupByLibrary.simpleMessage("消息通知"),
    "pageNotFound": MessageLookupByLibrary.simpleMessage("页面未找到"),
    "pageNotFoundDesc": MessageLookupByLibrary.simpleMessage(
      "抱歉，您访问的页面不存在或已被移动。",
    ),
    "password": MessageLookupByLibrary.simpleMessage("密码"),
    "passwordHint": MessageLookupByLibrary.simpleMessage("请输入密码"),
    "postsCount": m5,
    "postsCountFull": m6,
    "privacyContent1Desc": MessageLookupByLibrary.simpleMessage(
      "我们仅收集提供服务所需的最少个人信息，可能包括您的电子邮箱、用户名和使用偏好。我们不会出售或与第三方共享您的个人数据。",
    ),
    "privacyContent1Title": MessageLookupByLibrary.simpleMessage("信息收集"),
    "privacyContent2Desc": MessageLookupByLibrary.simpleMessage(
      "您的数据安全地存储在我们的服务器上，采用行业标准的加密技术。我们仅在提供服务所必需的期限或法律要求的期限内保留您的数据。",
    ),
    "privacyContent2Title": MessageLookupByLibrary.simpleMessage("数据存储"),
    "privacyContent3Desc": MessageLookupByLibrary.simpleMessage(
      "我们使用必要的 Cookie 以确保平台正常运行。可能会使用分析 Cookie 以改善用户体验，您可以在浏览器设置中禁用这些 Cookie。",
    ),
    "privacyContent3Title": MessageLookupByLibrary.simpleMessage("Cookie 与追踪"),
    "privacyContent4Desc": MessageLookupByLibrary.simpleMessage(
      "您有权随时访问、更正或删除您的个人数据。如有任何隐私相关请求，请联系我们的支持团队。",
    ),
    "privacyContent4Title": MessageLookupByLibrary.simpleMessage("您的权利"),
    "privacyPolicy": MessageLookupByLibrary.simpleMessage("隐私协议"),
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
    "termsContent1Desc": MessageLookupByLibrary.simpleMessage(
      "访问和使用本平台即表示您同意受本服务条款的约束。如果您不同意这些条款的任何部分，请勿使用本平台。",
    ),
    "termsContent1Title": MessageLookupByLibrary.simpleMessage("接受条款"),
    "termsContent2Desc": MessageLookupByLibrary.simpleMessage(
      "您有责任保管好您的账户信息。您同意不发布任何违法、有害、威胁、辱骂或其他不当内容。",
    ),
    "termsContent2Title": MessageLookupByLibrary.simpleMessage("用户责任"),
    "termsContent3Desc": MessageLookupByLibrary.simpleMessage(
      "用户不得试图未经授权访问我们的系统、干扰平台运营或使用自动化工具未经许可抓取或收集数据。",
    ),
    "termsContent3Title": MessageLookupByLibrary.simpleMessage("禁止行为"),
    "termsContent4Desc": MessageLookupByLibrary.simpleMessage(
      "我们保留随时修改本条款的权利。在条款变更后继续使用本平台，即表示您接受修改后的条款。",
    ),
    "termsContent4Title": MessageLookupByLibrary.simpleMessage("条款修改"),
    "termsOfService": MessageLookupByLibrary.simpleMessage("服务条款"),
    "themeColor": MessageLookupByLibrary.simpleMessage("主题色"),
    "themeLanguagePrefs": MessageLookupByLibrary.simpleMessage("主题、语言等偏好"),
    "today": MessageLookupByLibrary.simpleMessage("今天"),
    "username": MessageLookupByLibrary.simpleMessage("用户名"),
    "usernameHint": MessageLookupByLibrary.simpleMessage("请输入用户名"),
    "versionInfo": MessageLookupByLibrary.simpleMessage("版本信息和帮助"),
    "viewReadingHistory": MessageLookupByLibrary.simpleMessage("查看阅读记录"),
    "views": MessageLookupByLibrary.simpleMessage("浏览"),
    "weeksAgo": m9,
    "welcomeBack": MessageLookupByLibrary.simpleMessage("欢迎回来"),
    "writeComment": MessageLookupByLibrary.simpleMessage("写下你的评论..."),
    "yearMonthDay": m10,
    "yesterday": MessageLookupByLibrary.simpleMessage("昨天"),
  };
}
