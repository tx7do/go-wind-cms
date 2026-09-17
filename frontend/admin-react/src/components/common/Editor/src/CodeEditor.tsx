import { useCallback, useEffect, useRef } from 'react';

import hljs from 'highlight.js';
import * as monaco from 'monaco-editor';

import { initMonacoWorkers } from './monaco-loader';
import { isDarkMode } from './utils';

export interface CodeEditorProps {
  value: string;
  height?: number | string;
  disabled?: boolean;
  placeholder?: string;
  autoDetectLanguage?: boolean;
  options?: {
    fontSize?: number;
    language?: string;
    lineNumbers?: boolean;
    minimap?: boolean;
    tabSize?: number;
    theme?: 'dark' | 'hc-black' | 'light' | 'vs' | 'vs-dark';
    wordWrap?: 'bounded' | 'off' | 'on' | 'wordWrapColumn';
  };
  onChange?: (value: string) => void;
  onReady?: () => void;
  onError?: (error: Error) => void;
}

type EditorLanguage = string;

const languageMap: Record<string, EditorLanguage> = {
  javascript: 'javascript',
  js: 'javascript',
  typescript: 'typescript',
  ts: 'typescript',
  json: 'json',
  html: 'html',
  css: 'css',
  python: 'python',
  java: 'java',
  sql: 'sql',
  markdown: 'markdown',
  shell: 'shell',
  php: 'php',
  go: 'go',
  golang: 'go',
  ruby: 'ruby',
  c: 'c',
  'c++': 'cpp',
  cplusplus: 'cpp',
  cpp: 'cpp',
  'c#': 'csharp',
  csharp: 'csharp',
  lua: 'lua',
};

// Initialize monaco workers once
if (typeof window !== 'undefined') {
  initMonacoWorkers();
}

const CodeEditor: React.FC<CodeEditorProps> = ({
  value,
  height = '100%',
  disabled = false,
  autoDetectLanguage = true,
  options = {},
  onChange,
  onReady,
}) => {
  const containerRef = useRef<HTMLDivElement>(null);
  const editorRef = useRef<monaco.editor.IStandaloneCodeEditor | null>(null);
  const modelRef = useRef<monaco.editor.ITextModel | null>(null);
  const isUpdatingRef = useRef(false);
  const valueRef = useRef(value);

  const {
    language = 'javascript',
    theme,
    lineNumbers = true,
    tabSize = 2,
    minimap = false,
    fontSize = 14,
    wordWrap = 'on',
  } = options;

  const editorHeight =
    typeof height === 'number' ? `${Math.max(height, 200)}px` : height?.toString() || '100%';

  const getThemeName = useCallback((): 'vs' | 'vs-dark' | 'hc-black' => {
    if (theme && theme !== 'light' && theme !== 'dark') return theme;
    return isDarkMode() ? 'vs-dark' : 'vs';
  }, [theme]);

  const detectLanguage = useCallback(
    (content: string): EditorLanguage => {
      try {
        if (!content || content.trim() === '') {
          return language || 'plaintext';
        }
        if (content.trim().startsWith('{') || content.trim().startsWith('[')) {
          try {
            JSON.parse(content);
            return 'json';
          } catch {
            /* ignore */
          }
        }
        const detectedLanguage = hljs.highlightAuto(content);
        const detected = detectedLanguage.language;
        const detectedKey = typeof detected === 'string' ? detected.toLowerCase() : '';
        return languageMap[detectedKey] || 'plaintext';
      } catch (error) {
        return language || 'plaintext';
      }
    },
    [language],
  );

  // Initialize editor
  useEffect(() => {
    if (!containerRef.current) return;

    const initialLanguage = autoDetectLanguage ? detectLanguage(value || '') : language;

    modelRef.current = monaco.editor.createModel(value || '', initialLanguage);

    editorRef.current = monaco.editor.create(containerRef.current, {
      model: modelRef.current,
      theme: getThemeName(),
      automaticLayout: true,
      minimap: { enabled: minimap },
      lineNumbers: !lineNumbers ? 'off' : 'on',
      tabSize,
      insertSpaces: true,
      readOnly: disabled,
      scrollBeyondLastLine: false,
      wordWrap,
      fontSize,
      fontFamily: "'Monaco', 'Menlo', 'Ubuntu Mono', 'Consolas', monospace",
      lineHeight: 1.6,
      quickSuggestions: !disabled,
      codeLens: !disabled,
      folding: true,
      colorDecorators: true,
      renderLineHighlight: 'gutter',
      scrollbar: { vertical: 'visible', horizontal: 'auto' },
    });

    let changeTimeout: ReturnType<typeof setTimeout> | null = null;
    editorRef.current.onDidChangeModelContent(() => {
      if (isUpdatingRef.current || !editorRef.current) return;
      if (changeTimeout) clearTimeout(changeTimeout);
      changeTimeout = setTimeout(() => {
        const newValue = editorRef.current!.getValue() || '';
        valueRef.current = newValue;
        onChange?.(newValue);
      }, 100);
    });

    onReady?.();

    return () => {
      if (changeTimeout) clearTimeout(changeTimeout);
      editorRef.current?.dispose();
      editorRef.current = null;
      modelRef.current?.dispose();
      modelRef.current = null;
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  // Sync external value changes
  useEffect(() => {
    if (!editorRef.current || isUpdatingRef.current) return;
    const currentValue = editorRef.current.getValue();
    if (currentValue === value) return;

    isUpdatingRef.current = true;
    const currentPosition = editorRef.current.getPosition();
    editorRef.current.setValue(value || '');

    if (autoDetectLanguage && modelRef.current && value?.trim()) {
      const detected = detectLanguage(value);
      monaco.editor.setModelLanguage(modelRef.current, detected);
    }

    if (currentPosition && (value || '').length >= currentPosition.column) {
      editorRef.current.setPosition(currentPosition);
    }

    valueRef.current = value;
    isUpdatingRef.current = false;
  }, [value, autoDetectLanguage, detectLanguage]);

  // Theme changes
  useEffect(() => {
    const themeName = getThemeName();
    monaco.editor.setTheme(themeName);
    editorRef.current?.layout();
  }, [getThemeName]);

  // Disabled changes
  useEffect(() => {
    editorRef.current?.updateOptions({ readOnly: disabled });
  }, [disabled]);

  return (
    <div className="code-editor-wrapper">
      <div ref={containerRef} className="code-editor-container" style={{ height: editorHeight }} />
    </div>
  );
};

export default CodeEditor;
