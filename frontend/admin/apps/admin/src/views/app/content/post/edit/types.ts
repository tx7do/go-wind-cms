import type {
  contentservicev1_Post,
  contentservicev1_SeoMeta,
} from '#/api/generated/admin/service/v1';

import { EditorType } from '#/adapter/component/Editor';

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
  editorType: EditorType;
  status?: contentservicev1_Post['status'];

  code?: string;
  categoryIds?: number[];
  isFeatured?: boolean;
  sortOrder?: number;
  disallowComment?: boolean;
  publishTime?: string;
  thumbnail?: string;

  summary?: string;
  slug?: string;
  seo?: contentservicev1_SeoMeta;
}
