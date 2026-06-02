/**
 * API Hooks 索引文件
 * 导出所有业务模块的 hooks
 */

// 认证相关
export * from './auth';

// 用户资料
export * from './user-profile';

// 分类管理
export {
  useListCategories,
  fetchListCategories,
  useGetCategory,
  fetchCategory,
  useCreateCategory,
  useUpdateCategory,
  useDeleteCategory,
  getTranslation as getCategoryTranslation,
  getCategoryName,
  getCategoryDescription,
  getCategoryThumbnail,
} from './category';

// 评论管理
export * from './comment';

// 文件传输
export * from './file-transfer';

// 导航管理
export * from './navigation';

// 页面管理
export * from './page';

// 文章管理
export {
  useListPosts,
  fetchListPosts,
  useGetPost,
  fetchPost,
  useCreatePost,
  useUpdatePost,
  useDeletePost,
  getTranslation as getPostTranslation,
  getPostTitle,
  getPostSummary,
  getPostThumbnail,
  getPostContent,
} from './post';

// 标签管理
export {
  useListTags,
  fetchListTags,
  useGetTag,
  fetchTag,
  useCreateTag,
  useUpdateTag,
  useDeleteTag,
  getTranslation as getTagTranslation,
} from './tag';
