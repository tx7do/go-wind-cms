export type ContentType = 'markdown' | 'html' | 'text';

export interface ContentViewerProps {
    content?: string;
    type?: ContentType;
    className?: string;
}
