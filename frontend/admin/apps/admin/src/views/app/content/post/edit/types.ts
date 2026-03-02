import { EditorType } from '#/adapter/component/Editor';

/**
 * 文章编辑表单数据接口
 */
export interface PostEditProps {
  id?: number;
  title: string;
  content: string;
  lang: string;
  editorType: EditorType;
}
