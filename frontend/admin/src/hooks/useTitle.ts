import { useEffect } from 'react';

export function useTitle(title: string) {
  useEffect(() => {
    const prev = document.title;
    document.title = title;
    return () => {
      document.title = prev; // 离开时恢复
    };
  }, [title]);
}
