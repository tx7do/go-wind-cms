export enum EditorType {
  CODE = 'EDITOR_TYPE_CODE',
  JSON = 'EDITOR_TYPE_JSON_BLOCK',
  MARKDOWN = 'EDITOR_TYPE_MARKDOWN',
  PLAIN_TEXT = 'EDITOR_TYPE_PLAIN_TEXT',
  RICH_TEXT = 'EDITOR_TYPE_RICH_TEXT',
  VISUAL_BUILDER = 'EDITOR_TYPE_VISUAL_BUILDER',
}

export interface EditorProps {
  modelValue: string;
  editorType?: EditorType | string;
  height?: number | string;
  disabled?: boolean;
  placeholder?: string;
  // UEditor specific
  ueditorConfig?: Record<string, any>;
  // Markdown specific
  markdownOptions?: {
    hideModeSwitch?: boolean;
    initialEditType?: 'markdown' | 'wysiwyg';
    previewStyle?: 'tab' | 'vertical';
    toolbarItems?: string[][];
  };
  // JSON Editor specific
  jsonOptions?: {
    mode?: 'code' | 'form' | 'text' | 'tree' | 'view';
    modes?: string[];
    search?: boolean;
  };
}

export interface EditorEmits {
  (e: 'update:modelValue', value: string): void;
  (e: 'change', value: string): void;
  (e: 'ready'): void;
}
