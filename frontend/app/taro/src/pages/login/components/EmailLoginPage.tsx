import {useState} from 'react';
import {useTranslation} from 'react-i18next';

import {useAuthenticationStore} from "@/store/slices/authentication/hooks";

import '../index.scss';

export default function EmailLoginPage() {
  const {t} = useTranslation();

  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [rememberMe, setRememberMe] = useState(false);

  const authenticationStore = useAuthenticationStore();

  const handleLogin = async () => {
    if (!email || !password) {
      return;
    }
    console.log('登录信息:', {email, password, rememberMe});

    await authenticationStore.login({
      email: email,
      password: password,
    });
  };

  const handleForgotPassword = () => {
    console.log('忘记密码');
  };

  return (
    <div className="login-form">
      <div className="form-group">
        <label htmlFor="login-email">{t('authentication.register.email')}</label>
        <input
          id="login-email"
          type="email"
          className="input-field"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          placeholder={t('authentication.login.placeholder_email')}
          autoComplete="email"
        />
      </div>

      <div className="form-group">
        <label htmlFor="login-password-email">{t('authentication.login.password')}</label>
        <input
          id="login-password-email"
          type="password"
          className="input-field"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          placeholder={t('authentication.login.placeholder_password')}
          autoComplete="current-password"
          onKeyDown={(e) => e.key === 'Enter' && handleLogin()}
        />
      </div>

      <div className="form-options">
        <label className="checkbox-label">
          <input
            type="checkbox"
            checked={rememberMe}
            onChange={(e) => setRememberMe(e.target.checked)}
          />
          <span>{t('authentication.login.remember_me')}</span>
        </label>
        <button className="text-btn" onClick={handleForgotPassword}>
          {t('authentication.login.forgot_password')}
        </button>
      </div>

      <button className="login-button" onClick={handleLogin}>
        {t('authentication.login.login')}
      </button>
    </div>
  );
}
