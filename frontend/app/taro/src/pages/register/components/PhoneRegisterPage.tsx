import {useState, useEffect} from 'react';
import {useTranslation} from 'react-i18next';

import '../register.scss';

export default function PhoneRegisterPage() {
  const {t} = useTranslation();

  const [phone, setPhone] = useState('');
  const [verificationCode, setVerificationCode] = useState('');
  const [codeSent, setCodeSent] = useState(false);
  const [codeCountdown, setCodeCountdown] = useState(0);

  useEffect(() => {
    let codeTimer: number | null = null;

    if (codeSent && codeCountdown > 0) {
      codeTimer = window.setInterval(() => {
        setCodeCountdown((prev) => {
          if (prev <= 1) {
            setCodeSent(false);
            return 0;
          }
          return prev - 1;
        });
      }, 1000);
    }

    return () => {
      if (codeTimer !== null) {
        clearInterval(codeTimer);
      }
    };
  }, [codeSent, codeCountdown]);

  /**
   * 发送验证码
   */
  const handleButtonSendVerifyCode = () => {
    if (!phone) {
      return;
    }

    setCodeSent(true);
    setCodeCountdown(60);
  };

  /**
   * 登录或者注册
   */
  const handleButtonRegisterOrLogin = () => {
    if (!phone || !verificationCode) {
      return;
    }
    console.log('手机号注册/登录:', {phone, verificationCode});
  };

  return (
    <div className="register-form">
      {/* Phone Number Group */}
      <div className="form-group">
        <label htmlFor="register-phone-number">
          {t('authentication.register.phone')}
        </label>
        <input
          id="register-phone-number"
          type="text"
          value={phone}
          onChange={(e) => setPhone(e.target.value)}
          placeholder={t('authentication.register.input_phone')}
          autoComplete="tel"
          className="input-field"
        />
      </div>

      {/* Verification Code Group */}
      <div className="form-group">
        <label htmlFor="register-phone-code">
          {t('authentication.register.code')}
        </label>
        <div className="code-input-row">
          <input
            id="register-phone-code"
            type="text"
            value={verificationCode}
            onChange={(e) => setVerificationCode(e.target.value)}
            placeholder={t('authentication.register.input_code')}
            maxLength={6}
            autoComplete="one-time-code"
            className="input-field"
          />
          <button
            type="button"
            disabled={codeSent}
            className={`send-code-btn ${codeSent ? 'disabled' : ''}`}
            onClick={handleButtonSendVerifyCode}
          >
            {codeSent ? `${codeCountdown}s` : t('authentication.register.send_code')}
          </button>
        </div>
      </div>

      {/* Register Button */}
      <button
        type="button"
        className="register-button"
        onClick={handleButtonRegisterOrLogin}
      >
        {t('authentication.register.register_or_login')}
      </button>
    </div>
  );
}
