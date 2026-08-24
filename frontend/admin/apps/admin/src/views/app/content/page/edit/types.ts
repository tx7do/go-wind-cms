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
}
