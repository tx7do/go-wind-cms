import 'package:get_it/get_it.dart';

import 'package:flutter_app/generated/api/app/service/v1/index.dart'
    show ContentServiceV1Category, ContentServiceV1CategoryTranslation,
        ContentServiceV1Post, ContentServiceV1PostTranslation,
        ContentServiceV1Tag, ContentServiceV1TagTranslation;
import 'package:flutter_app/src/core/preference/user_preference_cache.dart';

typedef Category = ContentServiceV1Category;
typedef CategoryTranslation = ContentServiceV1CategoryTranslation;
typedef Post = ContentServiceV1Post;
typedef PostTranslation = ContentServiceV1PostTranslation;
typedef Tag = ContentServiceV1Tag;
typedef TagTranslation = ContentServiceV1TagTranslation;

// ==============================
// 当前语言偏好
// ==============================

/// 获取当前用户语言偏好，转换为 API 格式 (zh_CN → zh-CN)
String? getCurrentLocale() {
  try {
    final lang = GetIt.instance<UserPreferenceCache>().language;
    if (lang.isEmpty) return null;
    return lang.replaceAll('_', '-');
  } catch (_) {
    return null;
  }
}

// ==============================
// Category 翻译辅助
// ==============================

/// 获取分类的匹配翻译（优先当前语言，fallback 到第一个翻译）
CategoryTranslation? getCategoryTranslation(Category? category) {
  if (category?.translations == null || category!.translations!.isEmpty) {
    return null;
  }
  final locale = getCurrentLocale();
  if (locale != null) {
    final match = category.translations!.whereType<CategoryTranslation?>().firstWhere(
          (t) => t?.languageCode == locale,
          orElse: () => null,
        );
    if (match != null) return match;
  }
  return category.translations!.first;
}

/// 获取分类名称
String getCategoryName(Category? category, {String fallback = ''}) {
  return getCategoryTranslation(category)?.name ?? fallback;
}

/// 获取分类描述
String getCategoryDescription(Category? category) {
  return getCategoryTranslation(category)?.description ?? '';
}

/// 获取分类缩略图
String getCategoryThumbnail(Category? category, {String fallback = ''}) {
  return getCategoryTranslation(category)?.thumbnail ?? fallback;
}

// ==============================
// Post 翻译辅助
// ==============================

/// 获取文章的匹配翻译（优先当前语言，fallback 到第一个翻译）
PostTranslation? getPostTranslation(Post? post) {
  if (post?.translations == null || post!.translations!.isEmpty) return null;
  final locale = getCurrentLocale();
  if (locale != null) {
    final match = post.translations!.whereType<PostTranslation?>().firstWhere(
          (t) => t?.languageCode == locale,
          orElse: () => null,
        );
    if (match != null) return match;
  }
  return post.translations!.first;
}

/// 获取文章标题
String getPostTitle(Post? post, {String fallback = ''}) {
  return getPostTranslation(post)?.title ?? fallback;
}

/// 获取文章摘要
String getPostSummary(Post? post, {String fallback = ''}) {
  return getPostTranslation(post)?.summary ?? fallback;
}

/// 获取文章内容
String getPostContent(Post? post, {String fallback = ''}) {
  return getPostTranslation(post)?.content ?? fallback;
}

/// 获取文章缩略图
String? getPostThumbnail(Post? post) {
  return getPostTranslation(post)?.thumbnail;
}

// ==============================
// Tag 翻译辅助
// ==============================

/// 获取标签的匹配翻译（优先当前语言，fallback 到第一个翻译）
TagTranslation? getTagTranslation(Tag? tag) {
  if (tag?.translations == null || tag!.translations!.isEmpty) return null;
  final locale = getCurrentLocale();
  if (locale != null) {
    final match = tag.translations!.whereType<TagTranslation?>().firstWhere(
          (t) => t?.languageCode == locale,
          orElse: () => null,
        );
    if (match != null) return match;
  }
  return tag.translations!.first;
}

/// 获取标签名称
String getTagName(Tag? tag, {String fallback = ''}) {
  return getTagTranslation(tag)?.name ?? fallback;
}
