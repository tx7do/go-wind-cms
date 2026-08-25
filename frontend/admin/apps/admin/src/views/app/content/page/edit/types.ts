export interface PageEditProps {
  id?: number;
  title: string;
  slug: string;
  content: string;
  lang: string;
  editorType: string;
  parentId?: number;
  type?: string;
  status?: string;
  showInNavigation?: boolean;
  disallowComment?: boolean;
  authorId?: number;
  template?: string;
  isCustomTemplate?: boolean;
  customHead?: string;
  customFoot?: string;
  sortOrder?: number;
  contentModelId?: number;
  customFields?: Record<string, string>;
  /**
   * 页面嵌套区块列表。作为页面子部件随页面整体读写，
   * 对齐后端 Page.Sections 整体替换语义。
   */
  sections?: SectionFormItem[];
}

/**
 * 区块表单项：区块本身字段（语言无关）+ 各语言翻译内容。
 */
export interface SectionFormItem {
  id?: number;
  type?: string;
  name?: string;
  sortOrder?: number;
  config?: Record<string, string>;
  /**
   * 各语言翻译内容。键为语言代码，值为该语言的内容键值对。
   */
  translations?: SectionFormItemTranslation[];
}

export interface SectionFormItemTranslation {
  languageCode?: string;
  content?: Record<string, string>;
}

