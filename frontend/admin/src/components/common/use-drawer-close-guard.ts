import { useCallback, useEffect, useRef } from 'react';
import { App } from 'antd';
import { useTranslation } from 'react-i18next';

/**
 * 表单抽屉关闭守卫（规范见 docs/form_drawer_close_guard.md）
 *
 * 表单类 Drawer 的四条关闭路径（遮罩 / ESC / 标题栏 X / 底部取消）统一行为：
 * 无脏数据直接关闭；有脏数据弹确认框，确定后才关闭。
 * 脏 = 用户改过任何字段（onValuesChange 置脏，程序化回填不触发）；抽屉每次打开重置。
 * 提交成功后的关闭不经本守卫（页面直接调 onClose）。
 *
 * @example
 * ```tsx
 * const { guardClose, handleValuesChange } = useDrawerCloseGuard(open, onClose);
 * <DrawerForm
 *   open={open}
 *   onValuesChange={handleValuesChange}
 *   drawerProps={{ destroyOnHidden: true, size: 600, onClose: guardClose }}
 *   submitter={{ searchConfig: {...}, resetButtonProps: { onClick: guardClose } }}
 * />
 * ```
 */
export function useDrawerCloseGuard(open: boolean, onClose: () => void) {
  const { modal } = App.useApp();
  const { t } = useTranslation('common');
  const dirtyRef = useRef(false);

  useEffect(() => {
    if (open) dirtyRef.current = false;
  }, [open]);

  /** 关闭尝试（遮罩/ESC/X/取消按钮统一入口）：有脏数据先弹确认 */
  const guardClose = useCallback(() => {
    if (!dirtyRef.current) {
      onClose();
      return;
    }
    modal.confirm({
      title: t('drawer.unsavedMessage'),
      okText: t('drawer.discardButton'),
      cancelText: t('button.cancel'),
      onOk: () => {
        dirtyRef.current = false;
        onClose();
      },
    });
  }, [modal, t, onClose]);

  /** 挂到 DrawerForm / Form 的 onValuesChange：用户修改置脏 */
  const handleValuesChange = useCallback(() => {
    dirtyRef.current = true;
  }, []);

  /** 保存成功后手动清脏（后续关闭不再弹确认） */
  const markClean = useCallback(() => {
    dirtyRef.current = false;
  }, []);

  return { guardClose, handleValuesChange, markClean };
}
