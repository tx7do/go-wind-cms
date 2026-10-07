import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import i18next from 'i18next';

// 结构化 JSON 编辑器（tree/form 视图）在 vben 侧由 json-editor-vue 提供，该库
// 依赖 Vue 运行时，React 下不可用。当前实现为带解析校验的 textarea 兜底；
// Phase 1 若需 parity，可改用框架无关的 vanilla-jsoneditor 以命令式挂载。

export interface JsonEditorProps {
  value: string;
  height?: number | string;
  disabled?: boolean;
  placeholder?: string;
  options?: {
    mode?: 'code' | 'form' | 'text' | 'tree' | 'view';
    modes?: ('code' | 'form' | 'text' | 'tree' | 'view')[];
    search?: boolean;
  };
  onChange?: (value: string) => void;
  onReady?: () => void;
  onError?: (error: Error) => void;
}

const JsonEditor: React.FC<JsonEditorProps> = ({
  value,
  height = 500,
  disabled = false,
  placeholder = '{}',
  options: propOptions,
  onChange,
  onReady,
  onError,
}) => {
  const { mode = 'text', search = true } = propOptions || {};

  const [parseError, setParseError] = useState('');
  const [jsonData, setJsonData] = useState<any>(null);
  const localValueRef = useRef(value);
  const containerRef = useRef<HTMLDivElement>(null);

  // Calculate editor height
  const editorHeight = useMemo(() => {
    let baseHeight = 500;
    if (typeof height === 'number') {
      baseHeight = height;
    } else if (typeof height === 'string') {
      const numericHeight = Number(height);
      if (!Number.isNaN(numericHeight)) {
        baseHeight = numericHeight;
      } else if (height.endsWith('px')) {
        const pxValue = Number(height.replace('px', ''));
        if (!Number.isNaN(pxValue)) baseHeight = pxValue;
      } else {
        return height;
      }
    }
    return `${Math.max(baseHeight - 40, 200)}px`;
  }, [height]);

  // Validate and format JSON
  const validateAndFormat = useCallback(
    (val: string) => {
      try {
        if (!val?.trim()) {
          setParseError('');
          return { parsed: null, formatted: '' };
        }
        const parsed = JSON.parse(String(val));
        const formatted = JSON.stringify(parsed, null, 2);
        setParseError('');
        return { parsed, formatted };
      } catch (error) {
        const err = error as Error;
        setParseError(i18next.t('editor:jsonParseError', { error: err.message || i18next.t('editor:unknownError') }));
        onError?.(err);
        return { parsed: null, formatted: val };
      }
    },
    [onError],
  );

  // Initialize data
  useEffect(() => {
    const { parsed } = validateAndFormat(value);
    localValueRef.current = value;
    if (parsed !== null && typeof parsed === 'object') {
      setJsonData(parsed);
    } else if (parsed === null) {
      setJsonData({});
    } else {
      setJsonData({ value: parsed });
    }
  }, []);

  // Watch for external value changes
  useEffect(() => {
    if (value !== localValueRef.current) {
      const { parsed } = validateAndFormat(value);
      localValueRef.current = value;
      try {
        setJsonData(parsed || JSON.parse(placeholder));
      } catch {
        setJsonData({});
      }
    }
  }, [value, placeholder, validateAndFormat]);

  // Handle editor content change
  const handleEditorChange = useCallback(
    (newValue: any) => {
      if (typeof newValue === 'string') {
        localValueRef.current = newValue;
        onChange?.(newValue);
        const { parsed } = validateAndFormat(newValue);
        if (parsed !== null && typeof parsed === 'object') {
          setJsonData(parsed);
        }
        return;
      }
      if (Array.isArray(newValue) || (newValue !== null && typeof newValue === 'object')) {
        try {
          const serialized = JSON.stringify(newValue, null, 2);
          if (serialized !== localValueRef.current) {
            localValueRef.current = serialized;
            onChange?.(serialized);
          }
          setJsonData(newValue);
          setParseError('');
        } catch (error) {
          const err = error as Error;
          setParseError(
            i18next.t('editor:jsonSerializeError', {
              error: err.message || i18next.t('editor:unknownError'),
            }),
          );
          onError?.(err);
        }
      }
    },
    [onChange, onError, validateAndFormat],
  );

  // Ready event
  useEffect(() => {
    onReady?.();
  }, []);

  return (
    <div className="json-editor-container">
      {parseError && <div className="error-message">{parseError}</div>}
      <div
        ref={containerRef}
        style={{ height: editorHeight, width: '100%' }}
        className="json-editor-core"
      >
        {/* json-editor-vue will be rendered here via imperative approach below */}
        <JsonJsonEditorFallback
          jsonData={jsonData}
          mode={mode}
          disabled={disabled}
          search={search}
          onChange={handleEditorChange}
        />
      </div>
    </div>
  );
};

/**
 * Fallback: use a simple textarea for JSON editing since json-editor-vue
 * requires Vue runtime. For a full-featured JSON editor in React,
 * consider using @monaco-editor/react or a dedicated React JSON editor.
 */
const JsonJsonEditorFallback: React.FC<{
  jsonData: any;
  mode: string;
  disabled: boolean;
  search: boolean;
  onChange: (value: any) => void;
}> = ({ jsonData, disabled, onChange }) => {
  const [text, setText] = useState(() => {
    try {
      return jsonData ? JSON.stringify(jsonData, null, 2) : '';
    } catch {
      return '';
    }
  });

  useEffect(() => {
    try {
      const newText = jsonData ? JSON.stringify(jsonData, null, 2) : '';
      if (newText !== text) {
        setText(newText);
      }
    } catch {
      /* ignore */
    }
  }, [jsonData]);

  const handleChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    const val = e.target.value;
    setText(val);
    onChange(val);
  };

  return (
    <textarea
      value={text}
      disabled={disabled}
      onChange={handleChange}
      style={{
        width: '100%',
        height: '100%',
        padding: '8px',
        fontFamily: "Monaco, Consolas, 'Courier New', monospace",
        fontSize: '14px',
        lineHeight: 1.6,
        resize: 'none',
        border: 'none',
        outline: 'none',
        backgroundColor: 'var(--ant-color-bg-container)',
        color: 'var(--ant-color-text)',
        boxSizing: 'border-box',
      }}
    />
  );
};

export default JsonEditor;
