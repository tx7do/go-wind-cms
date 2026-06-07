import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';

import 'package:flutter_app/generated/l10n.dart';
import 'package:flutter_app/generated/api/models/content_service_v1_list_tag_response.dart';
import 'package:flutter_app/src/features/cms/services/tag_service.dart';
import 'package:flutter_app/src/core/constants/breakpoints.dart';
import 'package:flutter_app/src/core/widgets/responsive_layout.dart';
import 'package:flutter_app/src/features/cms/widgets/tag_chip.dart';


/// 标签列表页
///
/// 展示所有标签，点击可跳转到该标签下的文章列表。
class TagListPage extends StatefulWidget {
  const TagListPage({super.key});

  @override
  State<TagListPage> createState() => _TagListPageState();
}

class _TagListPageState extends State<TagListPage> {
  final _tagService = TagService();

  List<Tag> _tags = [];
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    _loadData();
  }

  Future<void> _loadData() async {
    final result = await _tagService.list();
    if (!mounted) return;

    setState(() {
      _tags = (result as ContentServiceV1ListTagResponse?)?.items ?? [];
      _isLoading = false;
    });
  }

  @override
  Widget build(BuildContext context) {
    if (_isLoading) {
      return const Scaffold(body: Center(child: CircularProgressIndicator()));
    }

    return ResponsiveLayout(
      mobileBody: _buildView(isMobile: true),
      webBody: _buildView(isMobile: false),
    );
  }

  Widget _buildView({required bool isMobile}) {
    final theme = Theme.of(context);
    final loc = S.of(context);

    return Scaffold(
      backgroundColor: theme.scaffoldBackgroundColor,
      appBar: AppBar(
        backgroundColor: theme.colorScheme.surface,
        surfaceTintColor: Colors.transparent,
        elevation: 0,
        leading: IconButton(
          icon: const Icon(Icons.arrow_back, size: 22),
          onPressed: () => Navigator.of(context).pop(),
        ),
        title: Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(
              Icons.label_outlined,
              size: isMobile ? 20.sp : 20,
              color: theme.colorScheme.primary,
            ),
            SizedBox(width: isMobile ? 6.w : 6),
            Text(
              loc.hotTags,
              style: TextStyle(
                fontSize: isMobile ? 18.sp : 18,
                fontWeight: FontWeight.bold,
              ),
            ),
          ],
        ),
        centerTitle: true,
      ),
      body: _tags.isEmpty
          ? Center(
              child: Text(
                loc.noRelatedPosts,
                style: TextStyle(
                  fontSize: 15,
                  color: theme.colorScheme.onSurface.withAlpha(120),
                ),
              ),
            )
          : _buildBody(isMobile),
    );
  }

  Widget _buildBody(bool isMobile) {
    if (!isMobile) {
      return Center(
        child: ConstrainedBox(
          constraints: const BoxConstraints(
            maxWidth: Breakpoints.webContentMaxWidth,
          ),
          child: Padding(
            padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 20),
            child: Wrap(
              spacing: 10,
              runSpacing: 10,
              children: _tags.map((tag) {
                return TagChip(
                  tag: tag,
                  isMobile: false,
                  showPostCount: true,
                );
              }).toList(),
            ),
          ),
        ),
      );
    }

    return Padding(
      padding: EdgeInsets.symmetric(horizontal: 16.w, vertical: 16.h),
      child: Wrap(
        spacing: 8.w,
        runSpacing: 8.h,
        children: _tags.map((tag) {
          return TagChip(
            tag: tag,
            isMobile: true,
            showPostCount: true,
          );
        }).toList(),
      ),
    );
  }
}
