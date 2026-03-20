import {useState} from 'react';
import {useTranslation} from 'react-i18next';

import {useAuthenticationStore} from "@/store/slices/authentication/hooks";

import '../index.scss';

export default function AccountLoginPage() {
  const {t} = useTranslation();

  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [rememberMe, setRememberMe] = useState(false);

  const authenticationStore = useAuthenticationStore();

  const handleLogin = async () => {
    if (!username || !password) {
      return;
    }
    console.log('登录信息:', {username, password, rememberMe});

    await authenticationStore.login({
      username: username,
      password: password,
    });
  };

  const handleForgotPassword = () => {
    console.log('忘记密码');
  };

  return (
    <div className="login-form">
      <div className="form-group">
        <label htmlFor="login-username">{t('authentication.login.username')}</label>
        <input
          id="login-username"
          type="text"
          className="input-field"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          placeholder={t('authentication.register.input_username')}
          autoComplete="username"
        />
      </div>

      <div className="form-group">
        <label htmlFor="login-password">{t('authentication.login.password')}</label>
        <input
          id="login-password"
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
