import { useEffect, useState } from 'react';
import type { CSSProperties } from 'react';
import { Avatar, Button, Input } from 'antd';
import { LockOutlined } from '@ant-design/icons';
import dayjs from 'dayjs';
import 'dayjs/locale/zh-cn';
import type { Dayjs } from 'dayjs';

import { useI18n } from '@/core/i18n';
import { usePreferencesStore } from '@/core/preferences/store';
import { useLockStore } from '@/stores';

interface LockScreenProps {
  avatar?: string;
  onToLogin: () => void;
}

/**
 * 时钟独立成子组件：只在锁屏可见时挂载（父组件 !visible 返回 null），
 * state 天然新鲜，无需在 effect 里同步 setState（react-hooks/set-state-in-effect）。
 */
const useNow = () => {
  const [now, setNow] = useState<Dayjs>(() => dayjs());

  useEffect(() => {
    const timer = window.setInterval(() => setNow(dayjs()), 1000);
    return () => window.clearInterval(timer);
  }, []);

  return now;
};

const LockClock = () => {
  const now = useNow();

  return (
    <div style={{ position: 'relative', fontSize: 160, fontWeight: 700, lineHeight: 1, color: 'var(--ant-color-text)' }}>
      {now.format('HH')}:{now.format('mm')}
      <span style={{ position: 'absolute', top: 12, marginLeft: 16, fontSize: 24, fontWeight: 500, color: 'var(--ant-color-text-secondary)' }}>
        {now.format('A')}
      </span>
    </div>
  );
};

const LockBottomBar = ({ locale, showTime }: { locale: string; showTime: boolean }) => {
  const now = useNow();
  const date = locale.startsWith('zh') ? now.locale('zh-cn').format('YYYY-MM-DD dddd') : now.format('YYYY-MM-DD dddd');

  return (
    <div style={{ position: 'absolute', bottom: 24, width: '100%', textAlign: 'center', color: 'var(--ant-color-text-secondary)' }}>
      {showTime && (
        <div style={{ fontSize: 28, marginBottom: 4, color: 'var(--ant-color-text)' }}>
          {now.format('HH:mm')} <span style={{ fontSize: 16 }}>{now.format('A')}</span>
        </div>
      )}
      <div style={{ fontSize: 24 }}>{date}</div>
    </div>
  );
};

/**
 * 全屏锁屏遮罩（行为基准：vue-vben lock-screen.vue）
 *
 * 时钟态与解锁表单态互斥：点「点击解锁」进表单，密码比对 lockScreenPassword，
 * 错了就地提示而非静默——静默失败会让人以为是输入框没生效。
 */
export const LockScreen = ({ avatar, onToLogin }: LockScreenProps) => {
  const { t } = useI18n('common');
  const locale = usePreferencesStore((state) => state.preferences.app.locale);
  const lockScreenEnabled = usePreferencesStore((state) => state.preferences.widget?.lockScreen ?? true);

  const isLockScreen = useLockStore((state) => state.isLockScreen);
  const lockScreenPassword = useLockStore((state) => state.lockScreenPassword);
  const unlockScreen = useLockStore((state) => state.unlockScreen);

  const [showUnlockForm, setShowUnlockForm] = useState(false);
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');

  const visible = lockScreenEnabled && isLockScreen;

  useEffect(() => {
    if (!visible) return undefined;
    const previous = document.body.style.overflow;
    document.body.style.overflow = 'hidden';
    return () => {
      document.body.style.overflow = previous;
    };
  }, [visible]);

  if (!visible) return null;

  const submit = () => {
    if (!password) return;
    if (password === lockScreenPassword) {
      unlockScreen();
      setPassword('');
      setError('');
    } else {
      setError(t('lockScreen.errorTip'));
    }
  };

  const toggleForm = () => {
    setShowUnlockForm((prev) => !prev);
    setPassword('');
    setError('');
  };

  const labelStyle: CSSProperties = {
    fontSize: 14,
    color: 'var(--ant-color-text-secondary)',
  };

  return (
    <div
      style={{
        position: 'fixed',
        inset: 0,
        zIndex: 2000,
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        justifyContent: 'center',
        gap: 24,
        background: 'var(--ant-color-bg-layout)',
      }}
    >
      {!showUnlockForm ? (
        <>
          <div
            style={{
              display: 'flex',
              flexDirection: 'column',
              alignItems: 'center',
              gap: 4,
              cursor: 'pointer',
              fontSize: 16,
              fontWeight: 600,
              color: 'var(--ant-color-text)',
            }}
            onClick={toggleForm}
          >
            <LockOutlined style={{ fontSize: 20 }} />
            <span style={labelStyle}>{t('lockScreen.unlock')}</span>
          </div>
          <LockClock />
        </>
      ) : (
        <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', width: 300 }} onKeyDown={(e) => e.key === 'Enter' && submit()}>
          <Avatar src={avatar} size={80} />
          <Input.Password
            autoFocus
            value={password}
            status={error ? 'error' : undefined}
            style={{ width: '100%', marginTop: 24 }}
            placeholder={t('lockScreen.placeholder')}
            onChange={(e) => {
              setPassword(e.target.value);
              setError('');
            }}
          />
          {error && <div style={{ width: '100%', marginTop: 8, fontSize: 12, color: 'var(--ant-color-error)' }}>{error}</div>}
          <Button type="primary" block style={{ marginTop: 16 }} onClick={submit}>
            {t('lockScreen.entry')}
          </Button>
          <Button block style={{ marginTop: 8 }} onClick={onToLogin}>
            {t('lockScreen.backToLogin')}
          </Button>
          <Button type="text" block style={{ marginTop: 8 }} onClick={toggleForm}>
            {t('button.back')}
          </Button>
        </div>
      )}

      <LockBottomBar locale={locale} showTime={showUnlockForm} />
    </div>
  );
};

export default LockScreen;
