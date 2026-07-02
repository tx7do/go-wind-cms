import { EditorType } from '#/adapter/component/Editor';

/**
 * 分类选项
 */
export interface CategoryOption {
  label: string;
  value: number;
}

/**
 * 文章编辑表单数据接口
 */
export interface PostEditProps {
  id?: number;
  title: string;
  content: string;
  lang: string;
  editorType: EditorType;
  categoryIds?: number[];
}
