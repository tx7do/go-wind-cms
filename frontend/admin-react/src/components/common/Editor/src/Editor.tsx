import { lazy, Suspense, useMemo } from 'react';

import { EditorType } from '../types';
import type { EditorProps } from '../types';

const LazyTiptapEditor = lazy(() => import('./TiptapEditor'));
const LazyMarkdownEditor = lazy(() => import('./MarkdownEditor'));
const LazyJsonEditor = lazy(() => import('./JsonEditor'));
const LazyPlainTextEditor = lazy(() => import('./PlainTextEditor'));
const LazyCodeEditor = lazy(() => import('./CodeEditor'));

const Editor: React.FC<EditorProps> = ({
  value,
  editorType = EditorType.MARKDOWN,
  height = '100%',
  disabled = false,
  placeholder,
  uploadImage,
  markdownOptions,
  jsonOptions,
  codeOptions,
  onChange,
  onReady,
}) => {
  const currentEditorComponent = useMemo(() => {
    switch (editorType) {
      case EditorType.CODE:
      case EditorType.VISUAL_BUILDER:
        return LazyCodeEditor;
      case EditorType.JSON:
        return LazyJsonEditor;
      case EditorType.MARKDOWN:
        return LazyMarkdownEditor;
      case EditorType.PLAIN_TEXT:
        return LazyPlainTextEditor;
      case EditorType.RICH_TEXT:
        return LazyTiptapEditor;
      default:
        return LazyMarkdownEditor;
    }
  }, [editorType]);

  const currentOptions = useMemo(() => {
    switch (editorType) {
      case EditorType.CODE:
        return codeOptions;
      case EditorType.JSON:
        return jsonOptions;
      case EditorType.MARKDOWN:
        return markdownOptions;
      default:
        return undefined;
    }
  }, [editorType, codeOptions, jsonOptions, markdownOptions]);

  const EditorComponent = currentEditorComponent;

  return (
    // height 必须落在这一层：编辑器各实现的根节点都以 height:'100%' 向上取高，
    // 若本层高度为 auto（默认无 CSS 规则），百分比失去参照，会退化为
    // md-editor/tiptap 的库内默认高度（约 500px），编辑页撑不满剩余空间。
    <div className="editor-container" style={{ height }}>
      <Suspense fallback={<div style={{ padding: 16, textAlign: 'center', color: 'var(--ant-color-text-secondary)' }}>Loading editor...</div>}>
        <EditorComponent
          value={value ?? ''}
          height={height}
          disabled={disabled}
          placeholder={placeholder}
          uploadImage={uploadImage}
          autoDetectLanguage={codeOptions?.autoDetectLanguage}
          // eslint-disable-next-line @typescript-eslint/no-explicit-any
          options={currentOptions as any}
          onChange={onChange}
          onReady={onReady}
        />
      </Suspense>
    </div>
  );
};

export default Editor;
