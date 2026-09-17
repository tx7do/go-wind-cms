import { useState, useCallback } from 'react';

export const useLayoutState = () => {
  const [collapsed, setCollapsed] = useState(false);
  const [isMobile, setIsMobile] = useState(false);
  const [selectedKeys, setSelectedKeys] = useState<string[]>([]);
  const [openKeys, setOpenKeys] = useState<string[]>([]);

  // 设置选中菜单项（用于高亮）
  const setSelectedKeysSafe = useCallback((keys: string[]) => {
    // 可选：同步到 URL 或 localStorage
    setSelectedKeys(keys);
  }, []);

  return {
    collapsed,
    setCollapsed,
    isMobile,
    setIsMobile,
    selectedKeys,
    setSelectedKeys: setSelectedKeysSafe,
    openKeys,
    setOpenKeys,
  };
};
