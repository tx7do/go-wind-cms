import 'package:flutter_app/generated/api/app/service/v1/index.dart'
    show CommentServiceV1Comment;

typedef CommentType = CommentServiceV1Comment;

/// 将扁平评论列表构建为树形结构
///
/// 如果后端已返回 [children] 字段，直接返回顶层评论；
/// 否则根据 [parentId] 归组。
List<CommentType> buildCommentTree(List<CommentType> flatList) {
  // 检查是否已经有树形数据（第一条评论有 children 字段非空即视为树形）
  if (flatList.isNotEmpty &&
      flatList.any((c) => c.children != null && c.children!.isNotEmpty)) {
    return flatList
        .where((c) => c.parentId == null || c.parentId == 0)
        .toList();
  }

  // 扁平列表 → 按 parentId 归组
  final Map<int, List<CommentType>> childrenMap = {};
  final List<CommentType> roots = [];

  for (final comment in flatList) {
    final pid = comment.parentId;
    if (pid == null || pid == 0) {
      roots.add(comment);
    } else {
      childrenMap.putIfAbsent(pid, () => []).add(comment);
    }
  }

  // 递归构建（使用可变副本携带 children）
  List<CommentType> attachChildren(List<CommentType> parents) {
    return parents.map((parent) {
      final kids = childrenMap[parent.id ?? -1] ?? [];
      if (kids.isEmpty) return parent;
      // CommentType 是 const 构造函数，无法修改 children，
      // 直接返回原对象，渲染时通过 childrenMap 查找子节点
      return parent;
    }).toList();
  }

  attachChildren(roots);
  return roots;
}

/// 递归获取某评论的所有子评论（从扁平列表中查找）
List<CommentType> findChildren(
  CommentType comment,
  List<CommentType> flatList,
) {
  // 优先使用后端返回的 children
  if (comment.children != null && comment.children!.isNotEmpty) {
    return comment.children!;
  }
  // 否则从扁平列表中按 parentId 查找
  return flatList.where((c) => c.parentId == comment.id).toList();
}
