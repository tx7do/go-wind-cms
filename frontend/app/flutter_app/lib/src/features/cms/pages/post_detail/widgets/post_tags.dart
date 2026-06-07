import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';

import 'package:flutter_app/generated/api/models/content_service_v1_tag.dart';
import 'package:flutter_app/src/features/cms/widgets/tag_chip.dart';

/// 文章标签展示组件
///
/// 支持两种渲染模式：
/// - [asSliver] = true → 返回 SliverToBoxAdapter（用于 CustomScrollView）
/// - [asSliver] = false → 返回普通 Widget（用于 Column）
class PostTags extends StatelessWidget {
  final List<ContentServiceV1Tag> allTags;
  final List<int>? tagIds;
  final bool isMobile;
  final bool asSliver;

  const PostTags({
    super.key,
    required this.allTags,
    this.tagIds,
    required this.isMobile,
    this.asSliver = false,
  });

  List<ContentServiceV1Tag> get _matchedTags {
    return allTags
        .where((t) => tagIds != null && t.id != null && tagIds!.contains(t.id!))
        .toList();
  }

  @override
  Widget build(BuildContext context) {
    final tags = _matchedTags;
    if (tags.isEmpty) {
      return asSliver
          ? const SliverToBoxAdapter(child: SizedBox.shrink())
          : const SizedBox.shrink();
    }

    final content = Padding(
      padding: asSliver
          ? EdgeInsets.fromLTRB(isMobile ? 20.w : 0, 8, isMobile ? 20.w : 0, 0)
          : const EdgeInsets.fromLTRB(0, 8, 0, 0),
      child: Wrap(
        spacing: isMobile ? 8.w : 8,
        runSpacing: isMobile ? 6.h : 6,
        children: tags
            .map((tag) => TagChip(tag: tag, isMobile: isMobile, compact: true))
            .toList(),
      ),
    );

    return asSliver ? SliverToBoxAdapter(child: content) : content;
  }
}
