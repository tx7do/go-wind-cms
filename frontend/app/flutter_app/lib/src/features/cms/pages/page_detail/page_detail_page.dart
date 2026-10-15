import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';

import 'package:flutter_app/generated/api/app/service/v1/index.dart'
    show ContentServiceV1Page;
import 'package:flutter_app/generated/l10n.dart';
import 'package:flutter_app/src/core/constants/breakpoints.dart';
import 'package:flutter_app/src/core/utils/translation_helpers.dart';
import 'package:flutter_app/src/core/widgets/app_back_button.dart';
import 'package:flutter_app/src/core/widgets/responsive_layout.dart';
import 'package:flutter_app/src/features/cms/pages/page_detail/widgets/section_renderer.dart';
import 'package:flutter_app/src/features/cms/services/page_service.dart';

typedef PageModel = ContentServiceV1Page;

/// CMS 页面详情页（Page + Section 区块渲染）
///
/// 由路由 `/page/:id` 进入；后台搭建的 Section 页面在此渲染为原生视图。
/// 仅已发布页面可达（app BFF Get 过滤），草稿/归档按未找到处理。
class PageDetailPage extends StatefulWidget {
  final int pageId;
  final String? slug;

  const PageDetailPage({super.key, this.pageId = 0, this.slug});

  @override
  State<PageDetailPage> createState() => _PageDetailPageState();
}

class _PageDetailPageState extends State<PageDetailPage> {
  final _pageService = PageService();

  PageModel? _page;
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    _loadData();
  }

  Future<void> _loadData() async {
    final result = await _pageService.get(
      widget.pageId,
      slug: widget.slug,
      locale: getCurrentLocale(),
    );
    if (!mounted) return;
    setState(() {
      // get 失败时返回 Status（见 BaseService.handleDioError），按未找到处理
      _page = result is PageModel ? result : null;
      _isLoading = false;
    });
  }

  @override
  Widget build(BuildContext context) {
    if (_isLoading) {
      return const Center(child: CircularProgressIndicator());
    }

    final page = _page;
    if (page == null) {
      final theme = Theme.of(context);
      return Scaffold(
        backgroundColor: theme.scaffoldBackgroundColor,
        appBar: AppBar(
          backgroundColor: theme.colorScheme.surface,
          surfaceTintColor: Colors.transparent,
          elevation: 0,
          leading: const AppBackButton(),
        ),
        body: Center(
          child: Text(
            S.of(context).pageNotFound,
            style: TextStyle(
              fontSize: 16,
              color: theme.colorScheme.onSurface.withAlpha(150),
            ),
          ),
        ),
      );
    }

    return ResponsiveLayout(
      mobileBody: _buildView(context, isMobile: true),
      webBody: _buildView(context, isMobile: false),
    );
  }

  Widget _buildView(BuildContext context, {required bool isMobile}) {
    final theme = Theme.of(context);
    final page = _page!;
    final title = getPageTitle(page);
    final cover = page.thumbnail;

    final content = Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        if (!isMobile) ...[
          // Web 端轻量返回栏
          Container(
            decoration: BoxDecoration(
              color: theme.colorScheme.surface,
              border: Border(
                bottom: BorderSide(
                  color: theme.colorScheme.onSurface.withAlpha(15),
                  width: 1,
                ),
              ),
            ),
            child: Center(
              child: ConstrainedBox(
                constraints: const BoxConstraints(
                  maxWidth: Breakpoints.webContentMaxWidth,
                ),
                child: Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 24),
                  child: Row(
                    children: [
                      AppBackButton(color: theme.colorScheme.onSurface),
                    ],
                  ),
                ),
              ),
            ),
          ),
          SizedBox(height: isMobile ? 0 : 16),
        ],
        Center(
          child: ConstrainedBox(
            constraints: BoxConstraints(
              maxWidth: isMobile ? double.infinity : 900,
            ),
            child: Padding(
              padding: EdgeInsets.symmetric(horizontal: isMobile ? 20.w : 24),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  if (cover != null && cover.isNotEmpty) ...[
                    ClipRRect(
                      borderRadius: BorderRadius.circular(isMobile ? 12.r : 12),
                      child: Image.network(
                        cover,
                        fit: BoxFit.cover,
                        width: double.infinity,
                        height: isMobile ? 200.h : 300,
                        errorBuilder: (_, _, _) => const SizedBox.shrink(),
                      ),
                    ),
                    SizedBox(height: isMobile ? 16.h : 24),
                  ],
                  if (title.isNotEmpty) ...[
                    Text(
                      title,
                      style: theme.textTheme.headlineSmall?.copyWith(
                        fontSize: isMobile ? 22.sp : 28,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    SizedBox(height: isMobile ? 16.h : 24),
                  ],
                  SectionRenderer(sections: page.sections, isMobile: isMobile),
                  SizedBox(height: isMobile ? 24.h : 32),
                ],
              ),
            ),
          ),
        ),
      ],
    );

    if (isMobile) {
      return Scaffold(
        backgroundColor: theme.scaffoldBackgroundColor,
        appBar: AppBar(
          backgroundColor: theme.colorScheme.surface,
          surfaceTintColor: Colors.transparent,
          elevation: 0,
          leading: const AppBackButton(),
        ),
        body: SingleChildScrollView(child: content),
      );
    }

    return Scaffold(
      backgroundColor: theme.scaffoldBackgroundColor,
      body: SingleChildScrollView(child: content),
    );
  }
}
