import {useState, useMemo} from 'react';
import {useTranslation} from 'react-i18next';

import '../register.scss';

export default function AccountRegisterPage() {
  const {t} = useTranslation();

  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');

  // 用户名验证（3-20个字符，只能包含字母、数字、下划线）
  const isValidUsername = useMemo(() => {
    if (!username) return false;
    const usernameRegex = /^[a-zA-Z0-9_]{3,20}$/;
    return usernameRegex.test(username);
  }, [username]);

  // 密码强度验证（至少 6 个字符）
  const isValidPassword = useMemo(() => {
    return password.length >= 6;
  }, [password]);

  // 确认密码验证
  const isPasswordMatch = useMemo(() => {
    return password === confirmPassword && confirmPassword.length > 0;
  }, [password, confirmPassword]);

  // 表单是否可提交
  const isFormValid = useMemo(() => {
    return isValidUsername && isValidPassword && isPasswordMatch;
  }, [isValidUsername, isValidPassword, isPasswordMatch]);

  const handleButtonRegister = () => {
    if (!isFormValid) {
      return;
    }

    console.log('注册信息：', {
      username: username,
      password: password,
    });
  };

  return (
    <div className="register-form">
      {/* Username */}
      <div className="form-group">
        <label htmlFor="register-account-username">
          {t('authentication.register.username')}
        </label>
        <input
          id="register-account-username"
          type="text"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          placeholder={t('authentication.register.input_username')}
          autoComplete="username"
          className={`input-field ${username && !isValidUsername ? 'error' : ''}`}
        />
        {username && !isValidUsername && (
          <span className="error-hint">{t('authentication.register.invalid_username')}</span>
        )}
      </div>

      {/* Password */}
      <div className="form-group">
        <label htmlFor="register-account-password">
          {t('authentication.register.password')}
        </label>
        <input
          id="register-account-password"
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          placeholder={t('authentication.register.input_password')}
          autoComplete="new-password"
          className={`input-field ${password && !isValidPassword ? 'error' : ''}`}
        />
        {password && !isValidPassword && (
          <span className="error-hint">{t('authentication.register.invalid_password')}</span>
        )}
      </div>

      {/* Confirm Password */}
      <div className="form-group">
        <label htmlFor="register-account-confirm-password">
          {t('authentication.register.confirm_password')}
        </label>
        <input
          id="register-account-confirm-password"
          type="password"
          value={confirmPassword}
          onChange={(e) => setConfirmPassword(e.target.value)}
          placeholder={t('authentication.register.input_confirm_password')}
          autoComplete="new-password"
          className={`input-field ${confirmPassword && !isPasswordMatch ? 'error' : ''}`}
        />
        {confirmPassword && !isPasswordMatch && (
          <span className="error-hint">{t('authentication.register.password_not_match')}</span>
        )}
      </div>

      {/* Register Button */}
      <button
        type="button"
        className="register-button"
        disabled={!isFormValid}
        onClick={handleButtonRegister}
      >
        {t('authentication.register.register')}
      </button>
    </div>
  );
}
