import Color from '@tiptap/extension-color';
import Highlight from '@tiptap/extension-highlight';
import HorizontalRule from '@tiptap/extension-horizontal-rule';
import Image from '@tiptap/extension-image';
import Link from '@tiptap/extension-link';
import Placeholder from '@tiptap/extension-placeholder';
import Subscript from '@tiptap/extension-subscript';
import Superscript from '@tiptap/extension-superscript';
import { Table } from '@tiptap/extension-table';
import TableCell from '@tiptap/extension-table-cell';
import TableHeader from '@tiptap/extension-table-header';
import TableRow from '@tiptap/extension-table-row';
import TaskItem from '@tiptap/extension-task-item';
import TaskList from '@tiptap/extension-task-list';
import TextAlign from '@tiptap/extension-text-align';
import { TextStyle } from '@tiptap/extension-text-style';
import Underline from '@tiptap/extension-underline';
import StarterKit from '@tiptap/starter-kit';
import { EditorContent, useEditor } from '@tiptap/react';
import Collaboration from '@tiptap/extension-collaboration';
import CollaborationCaret from '@tiptap/extension-collaboration-caret';
import { all, createLowlight } from 'lowlight';
import { useCallback, useEffect, useMemo, useRef } from 'react';
import { WebsocketProvider } from 'y-websocket';
import { IndexeddbPersistence } from 'y-indexeddb';
import * as Y from 'yjs';

import {
  CustomCodeBlockLowlight,
  CustomVideo,
  CustomIframe,
  LinkModal,
  CodeBlockModal,
  VideoModal,
  IframeModal,
  StatusBar,
  Toolbar,
  useTiptapEditor,
} from './TiptapEditor/index.ts';

import './tiptap-editor.css';

// 房间级协同资源缓存（模块级）：React StrictMode 模拟重挂载会复用 useMemo 值，
// 若在 cleanup 里 destroy，第二个挂载拿到的是已销毁的 provider/ydoc——
// 因此 bundle 只创建不销毁（随页面生命周期结束），房间数有限内存可控。
const collabBundleCache = new Map<
  string,
  { ydoc: Y.Doc; provider: WebsocketProvider; indexeddb: IndexeddbPersistence }
>();

function getCollabBundle(opts: NonNullable<CollaborationOptions>) {
  let bundle = collabBundleCache.get(opts.room);
  if (!bundle) {
    const ydoc = new Y.Doc();
    const provider = new WebsocketProvider(opts.wsUrl, opts.room, ydoc, {
      params: { token: opts.token },
    });
    // 本地持久化：离线草稿 + 服务重启后房间重建的种子来源
    const indexeddb = new IndexeddbPersistence(`gwc-collab:${opts.room}`, ydoc);
    bundle = { ydoc, provider, indexeddb };
    collabBundleCache.set(opts.room, bundle);
  }
  return bundle;
}

export interface CollaborationOptions {
  /** 协同房间名（如 post:12:zh-CN），即 ws URL 路径末段 */
  room: string;
  /** ws 服务地址（如 ws://localhost:6602） */
  wsUrl: string;
  /** 准入 token（query 参数传给服务端） */
  token: string;
  /** 协作者显示名（光标标签） */
  userName: string;
  /** 协作者光标颜色 */
  userColor: string;
}

export interface TiptapEditorProps {
  value: string;
  height?: number | string;
  disabled?: boolean;
  placeholder?: string;
  config?: Record<string, any>;
  showToolbar?: boolean;
  showStatusBar?: boolean;
  uploadImage?: (file: File) => Promise<string>;
  fullHeight?: boolean;
  onChange?: (value: string) => void;
  onReady?: (editor: any) => void;
  /** 传入即启用实时协同：多编辑者无冲突合并 + 在线光标；保存仍走页面显式保存 */
  collaboration?: CollaborationOptions | null;
}

const lowlight = createLowlight(all);

const TiptapEditor: React.FC<TiptapEditorProps> = ({
  value,
  height = 500,
  disabled = false,
  placeholder,
  showToolbar = true,
  showStatusBar = true,
  uploadImage,
  onChange,
  onReady,
  collaboration,
}) => {
  const {
    isDark,
    fileInputRef,
    markdownInputRef,
    linkModalVisible,
    setLinkModalVisible,
    linkUrl,
    setLinkUrl,
    codeBlockModalVisible,
    setCodeBlockModalVisible,
    codeBlockLanguage,
    setCodeBlockLanguage,
    codeBlockContent,
    setCodeBlockContent,
    videoModalVisible,
    setVideoModalVisible,
    videoUrl,
    setVideoUrl,
    videoWidth,
    setVideoWidth,
    iframeModalVisible,
    setIframeModalVisible,
    iframeUrl,
    setIframeUrl,
    iframeWidth,
    setIframeWidth,
    iframeHeight,
    setIframeHeight,
    iframeTitle,
    setIframeTitle,
    iframeAllowFullscreen,
    setIframeAllowFullscreen,
    isFullscreen,
    editorConfig,
    syncValue,
    createToolbarActions,
    handleImageUpload,
    handleMarkdownImport,
    handleLinkOk,
    handleCodeBlockOk,
    handleVideoOk,
    handleIframeOk,
    getStatusInfo,
    setupDarkModeTracking,
  } = useTiptapEditor({ value, disabled, placeholder, uploadImage, onChange, onReady });

  // 协同资源：ydoc/provider 按 room 生命周期创建；room 变化（换语言）即整体重建
  const collabRoom = collaboration?.room;
  const yjsBundle = useMemo(() => {
    if (!collaboration) return null;
    return getCollabBundle(collaboration);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [collabRoom]);

  const editor = useEditor(
    {
      ...editorConfig,
      extensions: [
      StarterKit.configure({
        link: false,
        underline: false,
        horizontalRule: false,
        codeBlock: false,
        // Collaboration 自带历史（yjs undo），与内置 undoRedo 互斥
        undoRedo: yjsBundle ? false : undefined,
      }),
      Underline,
      Subscript,
      Superscript,
      TaskList,
      TaskItem.configure({ nested: true }),
      HorizontalRule,
      Highlight.configure({ multicolor: true }),
      Color,
      TextStyle,
      CustomCodeBlockLowlight.configure({ lowlight, defaultLanguage: 'javascript' }),
      Table.configure({ resizable: true }),
      TableRow,
      TableHeader,
      TableCell,
      Placeholder.configure({ placeholder: placeholder || '' }),
      Link.configure({ openOnClick: false, autolink: true }),
      Image.configure({ inline: true }),
      CustomVideo,
      CustomIframe,
      TextAlign.configure({ types: ['heading', 'paragraph'] }),
        ...(yjsBundle
          ? [
              Collaboration.configure({ document: yjsBundle.ydoc }),
              CollaborationCaret.configure({
                provider: yjsBundle.provider,
                user: { name: collaboration?.userName ?? '匿名', color: collaboration?.userColor ?? '#30b8db' },
              }),
            ]
          : []),
      ],
    },
    [yjsBundle],
  );

  // Sync external value（协同模式下 value 只作种子，外部回写会与 ydoc 打架）
  useEffect(() => {
    if (editor && !yjsBundle) {
      syncValue(editor, value);
    }
  }, [value, editor, syncValue, yjsBundle]);

  // 首次同步后种子化：房间为空（新文档/服务重启）时以当前 HTML 灌入
  const seededRef = useRef(false);
  useEffect(() => {
    if (!yjsBundle || !editor) return undefined;
    seededRef.current = false;
    const seed = () => {
      if (seededRef.current) return;
      try {
        const frag = yjsBundle.ydoc.getXmlFragment('default');
        // 只有真的灌入才消耗一次性标志：value 异步晚到时（先 synced 后内容）
        // 允许在后续 sync/value 变更时重试
        if (frag.length === 0 && value) {
          seededRef.current = true;
          editor.commands.setContent(value, { emitUpdate: false });
        }
      } catch (error) {
        console.error('Seed collab document failed:', error);
      }
    };
    if (yjsBundle.provider.synced) seed();
    yjsBundle.provider.on('sync', seed);
    return () => {
      yjsBundle.provider.off('sync', seed);
    };
  }, [yjsBundle, editor, value]);

  // Disabled changes
  useEffect(() => {
    editor?.setEditable(!disabled);
  }, [disabled, editor]);

  // Dark mode tracking
  useEffect(() => {
    return setupDarkModeTracking();
  }, [setupDarkModeTracking]);

  // Prevent body scroll when fullscreen
  useEffect(() => {
    if (isFullscreen) {
      document.body.classList.add('tiptap-editor-fullscreen-active');
    } else {
      document.body.classList.remove('tiptap-editor-fullscreen-active');
    }
    return () => {
      document.body.classList.remove('tiptap-editor-fullscreen-active');
    };
  }, [isFullscreen]);

  // Toolbar actions
  const isActive = useCallback(
    (name: string, options?: any) => {
      return editor?.isActive(name, options) || false;
    },
    [editor],
  );

  const toolbarActions = useMemo(
    () => createToolbarActions(editor),
    [editor, createToolbarActions],
  );

  // Status info
  const statusInfo = useMemo(() => getStatusInfo(editor), [editor, getStatusInfo]);

  const containerHeight = typeof height === 'number' ? `${height}px` : height;

  return (
    <div
      className={`tiptap-editor-wrapper${isDark ? ' tiptap-editor-dark' : ''}${isFullscreen ? ' tiptap-editor-fullscreen' : ''}`}
      data-theme={isDark ? 'dark' : 'light'}
      style={{ height: containerHeight }}
    >
      {/* Toolbar */}
      {showToolbar && <Toolbar editor={editor} isActive={isActive} actions={toolbarActions} isFullscreen={isFullscreen} />}

      {/* Editor Content */}
      <div className="tiptap-editor-content-wrapper">
        <EditorContent
          editor={editor}
          className={`tiptap-editor-content${isDark ? ' dark' : ''}`}
        />
      </div>

      {/* Status Bar */}
      {showStatusBar && (
        <StatusBar
          words={statusInfo.words}
          chars={statusInfo.chars}
          cursor={statusInfo.cursor}
          isDark={isDark}
        />
      )}

      {/* Hidden file inputs */}
      <input
        ref={fileInputRef}
        type="file"
        accept="image/*"
        style={{ display: 'none' }}
        onChange={(e) => handleImageUpload(e, editor)}
      />
      <input
        ref={markdownInputRef}
        type="file"
        accept=".md,text/markdown"
        style={{ display: 'none' }}
        onChange={(e) => handleMarkdownImport(e, editor)}
      />

      {/* Modals */}
      <LinkModal
        open={linkModalVisible}
        url={linkUrl}
        onOk={() => handleLinkOk(editor)}
        onCancel={() => {
          setLinkModalVisible(false);
          setLinkUrl('');
          editor?.chain().focus().run();
        }}
        onUrlChange={setLinkUrl}
      />

      <CodeBlockModal
        open={codeBlockModalVisible}
        language={codeBlockLanguage}
        content={codeBlockContent}
        onOk={() => handleCodeBlockOk(editor)}
        onCancel={() => {
          setCodeBlockModalVisible(false);
          setCodeBlockContent('');
        }}
        onLanguageChange={setCodeBlockLanguage}
        onContentChange={setCodeBlockContent}
      />

      <VideoModal
        open={videoModalVisible}
        url={videoUrl}
        width={videoWidth}
        onOk={() => handleVideoOk(editor)}
        onCancel={() => setVideoModalVisible(false)}
        onUrlChange={setVideoUrl}
        onWidthChange={setVideoWidth}
      />

      <IframeModal
        open={iframeModalVisible}
        url={iframeUrl}
        width={iframeWidth}
        height={iframeHeight}
        title={iframeTitle}
        allowFullscreen={iframeAllowFullscreen}
        onOk={() => handleIframeOk(editor)}
        onCancel={() => setIframeModalVisible(false)}
        onUrlChange={setIframeUrl}
        onWidthChange={setIframeWidth}
        onHeightChange={setIframeHeight}
        onTitleChange={setIframeTitle}
        onAllowFullscreenChange={setIframeAllowFullscreen}
      />
    </div>
  );
};

export default TiptapEditor;
