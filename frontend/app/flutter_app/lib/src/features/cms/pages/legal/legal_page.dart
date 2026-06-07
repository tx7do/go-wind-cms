import 'package:flutter/material.dart';
import 'package:flutter_app/generated/l10n.dart';
import 'package:flutter_app/src/core/widgets/app_back_button.dart';
import 'package:flutter_app/src/core/widgets/responsive_layout.dart';

/// 法律/信息页面类型
enum LegalPageType { contact, disclaimer, privacy, terms }

/// 通用法律/信息页面
///
/// 用于展示联系我们、免责条款、隐私协议、服务条款等内容。
/// 所有文本通过 i18n 管理，方便多语言。
class LegalPage extends StatelessWidget {
  final LegalPageType type;

  const LegalPage({super.key, required this.type});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return ResponsiveLayout(
      mobileBody: _buildView(context, theme, isMobile: true),
      webBody: _buildView(context, theme, isMobile: false),
    );
  }

  Widget _buildView(
    BuildContext ctx,
    ThemeData theme, {
    required bool isMobile,
  }) {
    final s = S.of(ctx);
    final body = Center(
      child: ConstrainedBox(
        constraints: BoxConstraints(
          maxWidth: isMobile ? double.infinity : 720,
        ),
        child: SingleChildScrollView(
          padding: EdgeInsets.symmetric(
            horizontal: isMobile ? 20 : 48,
            vertical: 32,
          ),
          child: _buildContent(theme, s),
        ),
      ),
    );

    // Web 端由 WebShellLayout 提供 Scaffold
    if (!isMobile) return body;

    return Scaffold(
      backgroundColor: theme.scaffoldBackgroundColor,
      appBar: AppBar(
        backgroundColor: theme.colorScheme.surface,
        surfaceTintColor: Colors.transparent,
        elevation: 0,
        leading: const AppBackButton(),
        title: Text(_getTitle(s)),
      ),
      body: body,
    );
  }

  String _getTitle(S s) => switch (type) {
    LegalPageType.contact => s.contactUs,
    LegalPageType.disclaimer => s.disclaimer,
    LegalPageType.privacy => s.privacyPolicy,
    LegalPageType.terms => s.termsOfService,
  };

  Widget _buildContent(ThemeData theme, S s) {
    return switch (type) {
      LegalPageType.contact => _buildContact(theme, s),
      LegalPageType.disclaimer => _buildDisclaimer(theme, s),
      LegalPageType.privacy => _buildPrivacy(theme, s),
      LegalPageType.terms => _buildTerms(theme, s),
    };
  }

  // ─── 联系我们 ────────────────────────────────────────

  Widget _buildContact(ThemeData theme, S s) {
    return _SectionList(
      theme: theme,
      sections: [
        _Section(
          icon: Icons.email_outlined,
          title: s.contactEmail,
          content: s.contactEmailDesc,
        ),
        _Section(
          icon: Icons.language_outlined,
          title: s.contactWebsite,
          content: s.contactWebsiteDesc,
        ),
        _Section(
          icon: Icons.forum_outlined,
          title: s.contactCommunity,
          content: s.contactCommunityDesc,
        ),
      ],
    );
  }

  // ─── 免责条款 ────────────────────────────────────────

  Widget _buildDisclaimer(ThemeData theme, S s) {
    return _SectionList(
      theme: theme,
      sections: [
        _Section(
          icon: Icons.info_outline,
          title: s.disclaimerContent1Title,
          content: s.disclaimerContent1Desc,
        ),
        _Section(
          icon: Icons.link_off_outlined,
          title: s.disclaimerContent2Title,
          content: s.disclaimerContent2Desc,
        ),
        _Section(
          icon: Icons.gavel_outlined,
          title: s.disclaimerContent3Title,
          content: s.disclaimerContent3Desc,
        ),
      ],
    );
  }

  // ─── 隐私协议 ────────────────────────────────────────

  Widget _buildPrivacy(ThemeData theme, S s) {
    return _SectionList(
      theme: theme,
      sections: [
        _Section(
          icon: Icons.shield_outlined,
          title: s.privacyContent1Title,
          content: s.privacyContent1Desc,
        ),
        _Section(
          icon: Icons.storage_outlined,
          title: s.privacyContent2Title,
          content: s.privacyContent2Desc,
        ),
        _Section(
          icon: Icons.security_outlined,
          title: s.privacyContent3Title,
          content: s.privacyContent3Desc,
        ),
        _Section(
          icon: Icons.contact_mail_outlined,
          title: s.privacyContent4Title,
          content: s.privacyContent4Desc,
        ),
      ],
    );
  }

  // ─── 服务条款 ────────────────────────────────────────

  Widget _buildTerms(ThemeData theme, S s) {
    return _SectionList(
      theme: theme,
      sections: [
        _Section(
          icon: Icons.description_outlined,
          title: s.termsContent1Title,
          content: s.termsContent1Desc,
        ),
        _Section(
          icon: Icons.person_outline,
          title: s.termsContent2Title,
          content: s.termsContent2Desc,
        ),
        _Section(
          icon: Icons.block_outlined,
          title: s.termsContent3Title,
          content: s.termsContent3Desc,
        ),
        _Section(
          icon: Icons.update_outlined,
          title: s.termsContent4Title,
          content: s.termsContent4Desc,
        ),
      ],
    );
  }
}

/// 段落数据
class _Section {
  final IconData icon;
  final String title;
  final String content;

  const _Section({
    required this.icon,
    required this.title,
    required this.content,
  });
}

/// 段落列表
class _SectionList extends StatelessWidget {
  final ThemeData theme;
  final List<_Section> sections;

  const _SectionList({required this.theme, required this.sections});

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        for (int i = 0; i < sections.length; i++) ...[
          _SectionCard(section: sections[i], theme: theme),
          if (i < sections.length - 1) const SizedBox(height: 16),
        ],
      ],
    );
  }
}

/// 段落卡片
class _SectionCard extends StatelessWidget {
  final _Section section;
  final ThemeData theme;

  const _SectionCard({required this.section, required this.theme});

  @override
  Widget build(BuildContext context) {
    return Card(
      elevation: 0,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(14),
        side: BorderSide(
          color: theme.colorScheme.onSurface.withAlpha((0.06 * 255).round()),
        ),
      ),
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Row(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Container(
              width: 40,
              height: 40,
              decoration: BoxDecoration(
                color: theme.colorScheme.primaryContainer.withAlpha(120),
                borderRadius: BorderRadius.circular(10),
              ),
              child: Icon(
                section.icon,
                size: 20,
                color: theme.colorScheme.primary,
              ),
            ),
            const SizedBox(width: 16),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    section.title,
                    style: TextStyle(
                      fontSize: 15,
                      fontWeight: FontWeight.w600,
                      color: theme.colorScheme.onSurface,
                    ),
                  ),
                  const SizedBox(height: 6),
                  Text(
                    section.content,
                    style: TextStyle(
                      fontSize: 13,
                      color: theme.colorScheme.onSurface.withAlpha(180),
                      height: 1.6,
                    ),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}
