import type { EditorType } from '@/components/common/Editor';

/**
 * 文章编辑表单数据接口。
 * code/categoryIds/isFeatured/sortOrder/disallowComment/publishTime/thumbnail 为 Post 级字段；
 * summary/slug/seo 为翻译级字段（跟随当前编辑语言）。
 */
export interface PostEditProps {
  id?: number;
  title: string;
  content: string;
  lang: string;
  editorType: EditorType | string;
  status?: string;

  code?: string;
  categoryIds?: number[];
  // 自定义字段（Post 级，全语言共用）：组件值为 JSON 字符串
  customFields?: Record<string, string>;
  isFeatured?: boolean;
  sortOrder?: number;
  disallowComment?: boolean;
  publishTime?: string;
  thumbnail?: string;

  summary?: string;
  slug?: string;
  seo?: SeoMetaForm;
}

export interface SeoMetaForm {
  seoTitle?: string;
  metaKeywords?: string;
  metaDescription?: string;
  ogTitle?: string;
  ogDescription?: string;
  ogImage?: string;
  canonicalUrl?: string;
}
