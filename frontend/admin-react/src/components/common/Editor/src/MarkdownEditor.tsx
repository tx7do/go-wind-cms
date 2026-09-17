import { useCallback, useEffect, useMemo, useRef, useState } from 'react';

import { MdEditor } from 'md-editor-rt';
import 'md-editor-rt/lib/style.css';
import './markdown-editor.css';

import { hideUploadProgress, notifyUploadError, showUploadProgress } from '@/utils/upload-feedback';

import { isDarkMode } from './utils';

export interface MarkdownEditorProps {
  value: string;
  height?: number | string;
  disabled?: boolean;
  placeholder?: string;
  options?: Record<string, any>;
  enableExport?: boolean;
  uploadImage?: (
    file: File,
    onProgress?: (percent: number) => void,
  ) => Promise<string>;
  onChange?: (value: string) => void;
  onReady?: () => void;
}

const MarkdownEditor: React.FC<MarkdownEditorProps> = ({
  value,
  height = '100%',
  disabled = false,
  placeholder,
  options = {},
  enableExport = false,
  uploadImage,
  onChange,
  onReady,
}) => {
  const [localValue, setLocalValue] = useState(value);
  const [isDark, setIsDark] = useState(isDarkMode());
  const wrapperRef = useRef<HTMLDivElement>(null);
  const resizeObserverRef = useRef<ResizeObserver | null>(null);

  const theme = useMemo(() => (isDark ? 'dark' : 'light'), [isDark]);

  const toolbars = useMemo(() => {
    const base = [
      'bold',
      'underline',
      'italic',
      'strikeThrough',
      '-',
      'title',
      'sub',
      'sup',
      'alignLeft',
      'alignCenter',
      'alignRight',
      'alignJustify',
      '-',
      'quote',
      'unorderedList',
      'orderedList',
      'task',
      'indent',
      'outdent',
      '-',
      'codeRow',
      'code',
      'link',
      'image',
      'table',
      'horizontalRule',
      'emoji',
      'footnote',
      '-',
      'mermaid',
      'katex',
      '-',
      'revoke',
      'next',
      'clear',
      'save',
      '=',
      'pageFullscreen',
      'fullscreen',
      'preview',
      'htmlPreview',
      'catalog',
      'help',
    ];
    if (enableExport) {
      const saveIndex = base.indexOf('save');
      if (saveIndex !== -1) {
        base.splice(saveIndex + 1, 0, 'exportPdf', 'exportHtml', '-');
      }
    }
    return base as any[];
  }, [enableExport]);

  const editorProps = useMemo(
    () => ({
      preview: true,
      showCodeRowNumber: true,
      noMermaid: false,
      noKatex: false,
      toolbars,
      ...options,
    }),
    [toolbars, options],
  );

  const handleChange = useCallback(
    (val: string) => {
      setLocalValue(val);
      onChange?.(val);
    },
    [onChange],
  );

  const handleUploadImages = useCallback(
    async (files: File[], callback: (urls: string[]) => void) => {
      if (!uploadImage) {
        callback([]);
        return;
      }
      showUploadProgress(-1);
      try {
        const urls = await Promise.all(
          files.map((file) => uploadImage(file, showUploadProgress)),
        );
        callback(urls);
      } catch (error) {
        notifyUploadError(error);
      } finally {
        hideUploadProgress();
      }
    },
    [uploadImage],
  );

  const handleSave = useCallback((val: string) => {
    const blob = new Blob([val], { type: 'text/markdown;charset=utf-8' });
    const timestamp = new Date().toISOString().slice(0, 19).replaceAll(':', '-');
    const filename = `document-${timestamp}.md`;
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = filename;
    document.body.append(link);
    link.click();
    link.remove();
    URL.revokeObjectURL(url);
  }, []);

  // Sync external value
  useEffect(() => {
    if (value !== localValue) {
      setLocalValue(value);
    }
  }, [value]);

  // Dark mode tracking
  useEffect(() => {
    const checkDark = () => setIsDark(isDarkMode());
    checkDark();

    const observer = new MutationObserver(checkDark);
    observer.observe(document.documentElement, {
      attributes: true,
      attributeFilter: ['class', 'data-theme'],
    });
    return () => observer.disconnect();
  }, []);

  // Resize observer
  useEffect(() => {
    onReady?.();
    if (wrapperRef.current) {
      resizeObserverRef.current = new ResizeObserver(() => {
        // Editor will auto-resize due to automaticLayout
      });
      resizeObserverRef.current.observe(wrapperRef.current);
    }
    return () => {
      resizeObserverRef.current?.disconnect();
      resizeObserverRef.current = null;
    };
  }, []);

  const wrapperStyle = useMemo(
    () => ({
      height: height === '100vh' ? '100vh' : '100%',
    }),
    [height],
  );

  return (
    <div ref={wrapperRef} className="md-editor-wrapper" style={wrapperStyle}>
      <MdEditor
        value={localValue}
        onChange={handleChange}
        theme={theme}
        // 代码高亮主题由 md-editor-rt 运行时加载，暗色用深色主题
        codeTheme={isDark ? 'github-dark' : 'github'}
        placeholder={placeholder}
        disabled={disabled}
        onSave={handleSave}
        onUploadImg={handleUploadImages}
        className="md-editor-inner"
        // md-editor-rt 根元素默认按内容收缩（约 500px），必须显式给满高，
        // 否则 wrapper 已撑满而编辑器本体仍在内容高度处截断。
        style={{ height: '100%' }}
        toolbars={editorProps.toolbars}
        preview={editorProps.preview}
        showCodeRowNumber={editorProps.showCodeRowNumber}
        noMermaid={editorProps.noMermaid}
        noKatex={editorProps.noKatex}
      />
    </div>
  );
};

export default MarkdownEditor;
