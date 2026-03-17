import {useState, useEffect} from 'react';
import {useTranslations} from 'next-intl';
import styles from '../login.scss';

export default function PhoneLoginPage() {
  const t = useTranslations('authentication');
  const [phone, setPhone] = useState('');
  const [verificationCode, setVerificationCode] = useState('');
  const [codeSent, setCodeSent] = useState(false);
  const [countdown, setCountdown] = useState(0);

  useEffect(() => {
    let timer: number;
    if (countdown > 0) {
      timer = window.setTimeout(() => setCountdown(countdown - 1), 1000);
    } else if (codeSent && countdown === 0) {
      // 使用 setTimeout 延迟执行，避免在 effect 中同步调用 setState
      setTimeout(() => setCodeSent(false));
    }
    return () => clearTimeout(timer);
  }, [countdown, codeSent]);

  const handleSendCode = () => {
    if (!phone) return;
    setCodeSent(true);
    setCountdown(60);
  };

  const handleLogin = () => {
    if (!phone || !verificationCode) return;
    console.log('登录信息:', {phone, verificationCode});
  };

  return (
    <div className={styles['login-form']}>
      <div className={styles['form-group']}>
        <label htmlFor="login-phone">{t('register.phone')}</label>
        <input
          id="login-phone"
          type="tel"
          className={styles['input-field']}
          value={phone}
          onChange={(e) => setPhone(e.target.value)}
          placeholder={t('login.placeholder_phone')}
          autoComplete="tel"
        />
      </div>

      <div className={styles['form-group']}>
        <label htmlFor="login-code">{t('register.code')}</label>
        <div className={styles['code-input-row']}>
          <input
            id="login-code"
            type="text"
            className={styles['input-field']}
            value={verificationCode}
            onChange={(e) => setVerificationCode(e.target.value)}
            placeholder={t('login.placeholder_code')}
            maxLength={6}
            autoComplete="one-time-code"
            onKeyDown={(e) => e.key === 'Enter' && handleLogin()}
          />
          <button
            className={`${styles['send-code-btn']} ${codeSent ? styles.disabled : ''}`}
            onClick={handleSendCode}
            disabled={codeSent}
          >
            {codeSent ? `${countdown}s` : t('register.send_code')}
          </button>
        </div>
      </div>

      <button className={styles['login-button']} onClick={handleLogin}>
        {t('login.login')}
      </button>
    </div>
  );
}
