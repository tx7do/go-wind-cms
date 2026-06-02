'use client';

import {useEffect} from 'react';
import {useEditor, EditorContent} from '@tiptap/react';
import StarterKit from '@tiptap/starter-kit';
import Placeholder from '@tiptap/extension-placeholder';
import {Bold, Italic, Code, List, Strikethrough} from 'lucide-react';
import {Toggle} from '@/components/ui/toggle';
import {Button} from '@/components/ui/button';
import XIcon from '@/plugins/xicon';

interface RichTextEditorProps {
    value: string;
    onChange: (content: string) => void;
    onSubmit: () => void;
    submitting?: boolean;
    placeholder?: string;
    maxLength?: number;
    submitLabel?: string;
}

export default function RichTextEditor({
                                           value,
                                           onChange,
                                           onSubmit,
                                           submitting = false,
                                           placeholder = '写下你的评论...',
                                           maxLength = 1000,
                                           submitLabel = '提交评论',
                                       }: RichTextEditorProps) {
    const editor = useEditor({
        extensions: [
            StarterKit,
            Placeholder.configure({
                placeholder,
            }),
        ],
        content: value,
        onUpdate: ({editor}) => {
            onChange(editor.getHTML());
        },
        editorProps: {
            attributes: {
                class: 'prose dark:prose-invert focus:outline-none min-h-[120px] max-h-[300px] overflow-y-auto p-4 text-sm text-foreground',
            },
        },
    });

    // 外部 value 变化时同步到编辑器（如提交后清空）
    useEffect(() => {
        if (!editor) return;
        const currentHTML = editor.getHTML();
        if (value !== currentHTML) {
            editor.commands.setContent(value || '');
        }
    }, [editor, value]);

    if (!editor) return null;

    // 使用纯文本计算字数
    const textLength = editor.getText().length;
    const isEmpty = textLength === 0;

    // Ctrl/Cmd + Enter 快捷键提交
    const handleKeyDown = (e: React.KeyboardEvent) => {
        if ((e.ctrlKey || e.metaKey) && e.key === 'Enter') {
            e.preventDefault();
            if (!isEmpty && !submitting) {
                onSubmit();
            }
        }
    };

    return (
        <div
            className="w-full overflow-hidden rounded-xl border border-border bg-card transition-all focus-within:border-primary focus-within:ring-1 focus-within:ring-primary/40"
            onKeyDown={handleKeyDown}
        >
            {/* 工具列 (Toolbar) */}
            <div className="flex items-center gap-1 border-b border-border bg-muted/50 p-1.5">
                <Toggle
                    size="sm"
                    pressed={editor.isActive('bold')}
                    onPressedChange={() => editor.commands.toggleBold()}
                    className="h-8 w-8 p-0"
                >
                    <Bold className="h-4 w-4"/>
                </Toggle>

                <Toggle
                    size="sm"
                    pressed={editor.isActive('italic')}
                    onPressedChange={() => editor.commands.toggleItalic()}
                    className="h-8 w-8 p-0"
                >
                    <Italic className="h-4 w-4"/>
                </Toggle>

                <Toggle
                    size="sm"
                    pressed={editor.isActive('strike')}
                    onPressedChange={() => editor.commands.toggleStrike()}
                    className="h-8 w-8 p-0"
                >
                    <Strikethrough className="h-4 w-4"/>
                </Toggle>

                <Toggle
                    size="sm"
                    pressed={editor.isActive('codeBlock')}
                    onPressedChange={() => editor.commands.toggleCodeBlock()}
                    className="h-8 w-8 p-0"
                >
                    <Code className="h-4 w-4"/>
                </Toggle>

                <Toggle
                    size="sm"
                    pressed={editor.isActive('bulletList')}
                    onPressedChange={() => editor.commands.toggleBulletList()}
                    className="h-8 w-8 p-0"
                >
                    <List className="h-4 w-4"/>
                </Toggle>

                <div className="ml-auto flex items-center gap-2 text-xs text-muted-foreground">
                    <XIcon name="carbon:keyboard" size={14}/>
                    <span className="max-sm:hidden">Ctrl + Enter</span>
                </div>
            </div>

            {/* 富文本编辑输入区域 */}
            <EditorContent editor={editor}/>

            {/* 底部操作栏 */}
            <div className="flex items-center justify-between border-t border-border/50 bg-muted/20 p-2">
                <span className={`text-xs ${textLength > maxLength ? 'text-destructive' : 'text-muted-foreground'}`}>
                    {textLength} / {maxLength}
                </span>
                <Button
                    size="sm"
                    onClick={onSubmit}
                    disabled={submitting || isEmpty || textLength > maxLength}
                    className="gap-2"
                >
                    {submitting ? (
                        <XIcon name="carbon:circle-dash" className="h-4 w-4 animate-spin"/>
                    ) : (
                        <XIcon name="carbon:send-alt" className="h-4 w-4"/>
                    )}
                    {submitLabel}
                </Button>
            </div>
        </div>
    );
}
