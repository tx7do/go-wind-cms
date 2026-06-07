import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';

import 'package:flutter_app/generated/l10n.dart';
import 'package:flutter_app/generated/api/models/comment_service_v1_comment.dart';
import 'package:flutter_app/src/core/constants/breakpoints.dart';

class CommentInputBar extends StatefulWidget {
  final bool isMobile;
  final CommentServiceV1Comment? replyTo;
  final void Function(CommentServiceV1Comment? comment)? onReplyChanged;
  final Future<void> Function(String content, {int? parentId, int? replyToId})
      onSend;

  const CommentInputBar({
    super.key,
    required this.isMobile,
    this.replyTo,
    this.onReplyChanged,
    required this.onSend,
  });

  @override
  State<CommentInputBar> createState() => _CommentInputBarState();
}

class _CommentInputBarState extends State<CommentInputBar> {
  final _controller = TextEditingController();
  final _focusNode = FocusNode();
  bool _isSending = false;

  @override
  void initState() {
    super.initState();
    // 当回复目标变化时，自动聚焦输入框
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (widget.replyTo != null) {
        _focusNode.requestFocus();
      }
    });
  }

  @override
  void didUpdateWidget(covariant CommentInputBar oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (widget.replyTo != oldWidget.replyTo && widget.replyTo != null) {
      WidgetsBinding.instance.addPostFrameCallback((_) {
        _focusNode.requestFocus();
      });
    }
  }

  @override
  void dispose() {
    _controller.dispose();
    _focusNode.dispose();
    super.dispose();
  }

  Future<void> _handleSend() async {
    final text = _controller.text.trim();
    if (text.isEmpty || _isSending) return;

    setState(() => _isSending = true);
    try {
      await widget.onSend(
        text,
        parentId: widget.replyTo?.id,
        replyToId: widget.replyTo?.id,
      );
      _controller.clear();
      widget.onReplyChanged?.call(null);
    } finally {
      if (mounted) setState(() => _isSending = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final replyTo = widget.replyTo;

    return Container(
      padding: EdgeInsets.fromLTRB(
        widget.isMobile ? 16.w : 16,
        10,
        widget.isMobile ? 16.w : 16,
        10,
      ),
      decoration: BoxDecoration(
        color: theme.colorScheme.surface,
        border: Border(
          top: BorderSide(
            color:
                theme.colorScheme.onSurface.withAlpha((0.06 * 255).round()),
          ),
        ),
      ),
      child: SafeArea(
        top: false,
        child: Center(
          child: ConstrainedBox(
            constraints: const BoxConstraints(
              maxWidth: Breakpoints.webContentMaxWidth,
            ),
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                // 回复提示条
                if (replyTo != null)
                  Container(
                    padding: const EdgeInsets.symmetric(
                        horizontal: 10, vertical: 4),
                    margin: const EdgeInsets.only(bottom: 8),
                    decoration: BoxDecoration(
                      color: theme.colorScheme.primaryContainer
                          .withAlpha(80),
                      borderRadius: BorderRadius.circular(8),
                    ),
                    child: Row(
                      children: [
                        Icon(Icons.reply,
                            size: 16, color: theme.colorScheme.primary),
                        const SizedBox(width: 4),
                        Expanded(
                          child: Text(
                            '${S.of(context).reply} ${replyTo.authorName ?? ''}',
                            style: TextStyle(
                              fontSize: 12,
                              color: theme.colorScheme.primary,
                            ),
                          ),
                        ),
                        GestureDetector(
                          onTap: () =>
                              widget.onReplyChanged?.call(null),
                          child: Icon(
                            Icons.close,
                            size: 16,
                            color: theme.colorScheme.onSurface
                                .withAlpha(120),
                          ),
                        ),
                      ],
                    ),
                  ),
                Row(
                  children: [
                    Expanded(
                      child: Container(
                        padding:
                            const EdgeInsets.symmetric(horizontal: 14),
                        decoration: BoxDecoration(
                          color: theme
                              .colorScheme.surfaceContainerHighest
                              .withAlpha((0.5 * 255).round()),
                          borderRadius: BorderRadius.circular(20),
                        ),
                        child: TextField(
                          controller: _controller,
                          focusNode: _focusNode,
                          enabled: !_isSending,
                          style: TextStyle(
                            fontSize: widget.isMobile ? 14.sp : 14,
                            color: theme.colorScheme.onSurface,
                          ),
                          decoration: InputDecoration(
                            border: InputBorder.none,
                            hintText: S.of(context).writeComment,
                            hintStyle: TextStyle(
                              fontSize:
                                  widget.isMobile ? 14.sp : 14,
                              color: theme.colorScheme.onSurface
                                  .withAlpha(100),
                            ),
                            isDense: true,
                            contentPadding:
                                const EdgeInsets.symmetric(vertical: 8),
                          ),
                          textInputAction: TextInputAction.send,
                          onSubmitted: (_) => _handleSend(),
                        ),
                      ),
                    ),
                    const SizedBox(width: 10),
                    IconButton(
                      icon: _isSending
                          ? SizedBox(
                              width: 18,
                              height: 18,
                              child: CircularProgressIndicator(
                                strokeWidth: 2,
                                color: theme.colorScheme.primary,
                              ),
                            )
                          : Icon(
                              Icons.send,
                              size: 22,
                              color: theme.colorScheme.primary,
                            ),
                      onPressed: _isSending ? null : _handleSend,
                    ),
                  ],
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
