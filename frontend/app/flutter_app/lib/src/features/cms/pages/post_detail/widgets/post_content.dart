import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';

import 'package:flutter_app/generated/api/models/content_service_v1_post.dart';
import 'package:flutter_app/src/core/utils/translation_helpers.dart';
import 'package:flutter_app/src/features/cms/widgets/content_viewer.dart';

typedef Post = ContentServiceV1Post;

class PostContent extends StatelessWidget {
  final Post post;
  final bool isMobile;

  const PostContent({
    super.key,
    required this.post,
    required this.isMobile,
  });

  String get _summary => getPostSummary(post);

  String get _content => getPostContent(post);

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Padding(
      padding: EdgeInsets.fromLTRB(
        isMobile ? 20.w : 0,
        0,
        isMobile ? 20.w : 0,
        8,
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // 摘要
          Container(
            padding: EdgeInsets.all(isMobile ? 14.w : 14),
            decoration: BoxDecoration(
              color: theme.colorScheme.primaryContainer.withAlpha(
                (0.2 * 255).round(),
              ),
              borderRadius: BorderRadius.circular(isMobile ? 12.r : 12),
              border: Border(
                left: BorderSide(color: theme.colorScheme.primary, width: 3),
              ),
            ),
            child: Text(
              _summary,
              style: TextStyle(
                fontSize: isMobile ? 14.sp : 14,
                color: theme.colorScheme.onSurface.withAlpha(180),
                height: 1.6,
                fontStyle: FontStyle.italic,
              ),
            ),
          ),
          SizedBox(height: isMobile ? 16.h : 16),
          // 正文 — 使用 ContentViewer 根据 editorType 自适应渲染
          ContentViewer(
            content: _content,
            editorType: post.editorType,
            isMobile: isMobile,
          ),
        ],
      ),
    );
  }
}
