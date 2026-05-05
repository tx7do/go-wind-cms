export type EditorType =
  | 'markdown'
  | 'richtext'
  | 'code'
  | 'json'
  | 'text'

export interface EditorProps {
  modelValue: string;
  editorType?: EditorType | string;
  height?: number | string;
  disabled?: boolean;
  placeholder?: string;
  uploadImage?: (file: File) => Promise<string>;
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
  // Code Editor specific
  codeOptions?: {
    language?: string;
    lineNumbers?: boolean;
    tabSize?: number;
  };
}

export interface EditorEmits {
  (e: 'update:modelValue', value: string): void
  (e: 'change', value: string): void
  (e: 'typeChange', type: EditorType): void
  (e: 'ready'): void;
}
