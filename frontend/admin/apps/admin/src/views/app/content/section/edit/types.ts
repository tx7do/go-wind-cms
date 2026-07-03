export interface SectionEditProps {
  id?: number;
  pageId?: number;
  type?: string;
  name: string;
  sortOrder?: number;
  config?: Record<string, string>;
  // 当前语言下的正文（映射到 translations[].content.body）
  content: string;
  lang: string;
}
