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
    <div className="editor-container">
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
