export interface CategoryEditProps {
  id?: number;
  name: string;
  slug: string;
  description: string;
  lang: string;
  parentId?: number;
  icon?: string;
  isNav?: boolean;
  sortOrder?: number;
  status?: string;
  contentModelId?: number;
  customFields?: Record<string, string>;
}
