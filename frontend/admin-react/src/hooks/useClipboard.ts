import {useState, useEffect} from 'react';
import i18next from 'i18next';

type CopyHandler = (text: string) => Promise<boolean>;

export function useClipboard(): [boolean, string, CopyHandler] {
    const [isCopied, setIsCopied] = useState(false);
    const [error, setError] = useState('');

    useEffect(() => {
        let timer: ReturnType<typeof setTimeout>;
        if (isCopied) {
            timer = setTimeout(() => setIsCopied(false), 1000);
        }
        return () => clearTimeout(timer);
    }, [isCopied]);

    const copyText: CopyHandler = async (text) => {
        try {
            // 现代Clipboard API
            if (navigator.clipboard?.writeText) {
                await navigator.clipboard.writeText(text);
                setIsCopied(true);
                return true;
            }

            // 兼容旧浏览器的execCommand方法
            const textArea = document.createElement('textarea');
            textArea.value = text;
            textArea.style.position = 'fixed'; // 避免滚动
            document.body.appendChild(textArea);
            textArea.select();

            const success = document.execCommand('copy');
            document.body.removeChild(textArea);

            if (success) {
                setIsCopied(true);
                return true;
            }

            throw new Error(i18next.t('common:clipboard.copyFailed'));
        } catch (err) {
            const message = err instanceof Error ? err.message : i18next.t('common:clipboard.copyDenied');
            setError(message);
            setIsCopied(false);
            return false;
        }
    };

    return [isCopied, error, copyText];
}
